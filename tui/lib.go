package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vpofe/which-fix-version/git"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	isPending  bool
	isDone     bool
	commitHash string
	spinner    spinner.Model
	fixVersion string
}

var url = "git@github.com:vpofe/which-fix-version.git"

func InitialModel() Model {
	m := Model{
		inputs:    make([]textinput.Model, 5),
		isPending: false,
		isDone:    false,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Commit Hash"
			t.CursorStyle = cursorStyle
			t.Focus()
			t.CharLimit = 40
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.SetValue("4d5671af9791dc8b03f4596505d98fa9f3757e7d")
		case 1:
			t.Placeholder = "Repository URL"
			t.CharLimit = 100
			t.SetValue(url)
		case 2:
			t.Placeholder = "Remote Name"
			t.CharLimit = 100
			t.SetValue("origin")
		case 3:
			t.Placeholder = "Development Branch Name"
			t.CharLimit = 20
			t.SetValue("main")
		case 4:
			t.Placeholder = "Release Identifiers"
			t.CharLimit = 120
			t.SetValue("release- release/ releases/")
		}

		m.inputs[i] = t
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.spinner = s

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type fixVersionMsg string
type errMsg error

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case fixVersionMsg:
		m.isPending = false
		m.isDone = true
		m.fixVersion = string(msg)

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

			// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			if m.isDone {
				return m, tea.Quit
			}

			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.commitHash = m.inputs[0].Value()
				m.isPending = true
				return m, tea.Batch(m.spinner.Tick, m.findFixVersion)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)

		}

		// Handle character input and blinking
		cmd := m.updateInputs(msg)

		return m, cmd

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.isDone {
		return fmt.Sprintf("\n\n Fix version = %s", m.fixVersion)
	}

	if m.isPending {
		str := fmt.Sprintf("\n\n   %s Scanning release branches for %s...press q to quit\n\n", m.spinner.View(), m.commitHash)
		return str
	}

	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func (m Model) findFixVersion() tea.Msg {
	// FIXME make readable map for inputs
	repoUrl := m.inputs[1].Value()

	releaseIdentifiers := make([]string, 0)

	releaseIdentifiers = append(releaseIdentifiers, strings.Split(m.inputs[4].Value(), " ")...)

	rootCandidates, releases := git.FormatRemoteBranches(repoUrl, m.inputs[3].Value(), releaseIdentifiers, m.inputs[2].Value())

	// fetch commit list from ma(in/ster)
	root := git.SelectRoot(rootCandidates)
	// check latest release

	sortedReleases := git.GetSortedReleases(releases)

	c := git.GetRootCommit(repoUrl, m.commitHash, root)

	var message string

	if c == nil {
		message = "No such hash in the root of this repo"
		return fixVersionMsg(message)
	} else {
		message = "No fixed version found"

		fixedVersions := make([]string, 0)

		for _, version := range sortedReleases {
			if git.IsCommitPresentOnBranch(repoUrl, c, releases[version], m.inputs[2].Value()) {
				fixedVersions = append(fixedVersions, version)
			} else {
				// Cancel looking further if previous doesn't have a fixed version any longer
				if len(fixedVersions) > 0 {
					break
				}
			}
		}

		if len(fixedVersions) > 0 {
			return fixVersionMsg(fixedVersions[len(fixedVersions)-1])
		} else {
			return fixVersionMsg("No fixed version found")
		}
	}
}
