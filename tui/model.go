package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpofe/which-fix-version/git"
)

type SearchMode int64

const (
	Local SearchMode = iota
	Remote
)

func (searchMode SearchMode) String() string {
	switch searchMode {
	case Local:
		return "local"
	case Remote:
		return "remote"
	}
	return "unknown"
}

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
	searchMode SearchMode
}

func InitialModel(gc *git.GitConfig, searchMode SearchMode) Model {

	mb := NewBuilder()

	m := mb.InitUI().AddInputs(searchMode, gc).AddSpinner().Build()

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
