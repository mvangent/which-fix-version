package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vpofe/which-fix-version/git"
)

type Model struct {
	isInit     bool
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	isPending  bool
	isDone     bool
	commitHash string
	spinner    spinner.Model
	fixVersion string
}

func InitialModel(gc *git.GitConfig) Model {
	m := Model{
		inputs:    make([]textinput.Model, 5),
		isPending: false,
		isDone:    false,
		isInit:    true,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Commit Hash"
			t.Focus()
			t.CharLimit = 40
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.SetValue(gc.CommitHash)
		case 1:
			t.Placeholder = "Repository URL"
			t.CharLimit = 100
			t.SetValue(gc.URL)
		case 2:
			t.Placeholder = "Remote Name"
			t.CharLimit = 100
			t.SetValue(gc.RemoteName)
		case 3:
			t.Placeholder = "Development Branch Name"
			t.CharLimit = 20
			t.SetValue(gc.DevelopBranchName)
		case 4:
			t.Placeholder = "Release Identifiers"
			t.CharLimit = 120
			t.SetValue(strings.Join(gc.ReleaseBranchPrependIdentifiers, " "))
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

