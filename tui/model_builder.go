package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vpofe/which-fix-version/git"
)

type ModelBuilder interface {
	Build() Model

	AddInputs(searchMode SearchMode, gc *git.GitConfig) ModelBuilder
	AddSpinner() ModelBuilder
	InitUI() ModelBuilder
}

type modelBuilder struct{ m Model }

func NewBuilder() ModelBuilder {
	builder := &modelBuilder{m: Model{}}

	return builder
}

func (mb modelBuilder) Build() Model {
	return mb.m
}

func (mb modelBuilder) InitUI() ModelBuilder {
	mb.m.isPending = false
	mb.m.isDone = false
	mb.m.isInit = true

	return mb
}

func (mb modelBuilder) AddSpinner() ModelBuilder {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	mb.m.spinner = s

	return mb

}

func (mb modelBuilder) AddInputs(searchMode SearchMode, gc *git.GitConfig) ModelBuilder {
	m := mb.m
	switch searchMode {
	case Remote:
		m.searchMode = Remote

		m.inputs = make([]textinput.Model, 5)

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
				t.Placeholder = "Development Branch Name"
				t.CharLimit = 20
				t.SetValue(gc.DevelopBranchName)
				t.Blur()
			case 2:
				t.Placeholder = "Release Identifiers"
				t.CharLimit = 120
				t.SetValue(strings.Join(gc.ReleaseBranchPrependIdentifiers, " "))
				t.Blur()
			case 3:
				t.Placeholder = "Repository URL"
				t.CharLimit = 100
				t.SetValue(gc.URL)
				t.Blur()
			case 4:
				t.Placeholder = "Remote Name"
				t.CharLimit = 100
				t.SetValue(gc.RemoteName)
				t.Blur()
			}

			m.inputs[i] = t
		}
	case Local:
		m.searchMode = Local

		m.inputs = make([]textinput.Model, 4)

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
				t.Placeholder = "Development Branch Name"
				t.CharLimit = 20
				t.SetValue(gc.DevelopBranchName)
				t.Blur()
			case 2:
				t.Placeholder = "Release Identifiers"
				t.CharLimit = 120
				t.SetValue(strings.Join(gc.ReleaseBranchPrependIdentifiers, " "))
				t.Blur()
			case 3:
				t.Placeholder = "Path"
				t.CharLimit = 120
				t.SetValue(gc.Path)
				t.Blur()
			}

			m.inputs[i] = t
		}
	default:
		panic("Invalid searchMode")

	}

	mb.m = m

	return mb
}
