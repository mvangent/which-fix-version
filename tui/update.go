package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) executeSearchCommand() (tea.Model, tea.Cmd) {
	m.commitHash = m.inputs[0].Value()
	m.isPending = true
	switch m.searchMode {
	case Remote:
		return m, tea.Batch(m.spinner.Tick, m.findFixVersionRemote)
	case Local:
		return m, tea.Batch(m.spinner.Tick, m.findFixVersionLocal)
	default:
		panic("Invalid SearchMode, abort..")
	}

}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	skipInteraction := true

	for i := 0; i < len(m.inputs); i++ {
		if len(m.inputs[i].Value()) == 0 {
			skipInteraction = false
			break
		}
	}

	if skipInteraction && m.isInit {
		m.isInit = false
		return m.executeSearchCommand()
	}

	m.isInit = false

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
				return m.executeSearchCommand()
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
					m.inputs[i].CursorStyle = focusedStyle
				} else {
					// Remove focused state
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
					m.inputs[i].CursorStyle = noStyle
				}
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
