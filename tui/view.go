package tui

import (
	"fmt"
	"strings"
)

/* Static UI components */
var (
	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func (m Model) View() string {
	/* Start of State */

	// ui state
	isDone := m.isDone
	isPending := m.isPending
	// components
	inputs := m.inputs
	spinner := m.spinner
	// data
	commitHash := m.commitHash
	fixVersion := m.fixVersion
	currentVersion := m.currentVersion

	/* End of State */

	/* Start of UI */
	var b strings.Builder

	// Results
	if isDone {
		b.WriteString(fmt.Sprintf("\n Fix version = %s", fixVersion))
	}

	// Search in progress
	if isPending {
		b.WriteString(fmt.Sprintf("\n\n  %s Scanning release branch: %s, for %s...press q to quit\n\n", spinner.View(), currentVersion, commitHash))
	}

	// Interactive Inputs
	if !isPending && !isDone {

		// Input fields
		for i := range inputs {
			b.WriteString(inputs[i].View())
			if i < len(inputs)-1 {
				b.WriteRune('\n')
			}
		}

		// Submit button
		button := &blurredButton
		if m.focusIndex == len(inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		// Mode Info
		b.WriteString(helpStyle.Render("cursor mode is "))
		b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
		b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	}

	b.WriteString(helpStyle.Render("\n\npress Esc to quit"))

	return b.String()

	/* End of UI */
}
