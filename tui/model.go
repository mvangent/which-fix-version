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
		t.CharLimit = 32
		switch i {
		case 0:
			t.CursorStyle = cursorStyle
			t.Placeholder = "Commit Hash"
			t.CharLimit = 40
			t.SetValue(gc.CommitHash)
		case 1:
			t.Placeholder = "Repository URL"
			t.CharLimit = 100
			t.SetValue(gc.URL)
			t.Blur()
		case 2:
			t.Placeholder = "Remote Name"
			t.CharLimit = 100
			t.SetValue(gc.RemoteName)
			t.Blur()

		case 3:
			t.Placeholder = "Development Branch Name"
			t.CharLimit = 20
			t.SetValue(gc.DevelopBranchName)
			t.Blur()

		case 4:
			t.Placeholder = "Release Identifiers"
			t.CharLimit = 120
			t.SetValue(strings.Join(gc.ReleaseBranchPrependIdentifiers, " "))
			t.Blur()

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
