package tui

import (
	"fmt"
	"strings"

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


func (m Model) View() string {
	if m.isDone {
		return fmt.Sprintf("\n\n Fix version = %s", m.fixVersion)
	}

	if m.isPending {
		str := fmt.Sprintf("\n\n  %s Scanning release branches for %s...press q to quit\n\n", m.spinner.View(), m.commitHash)
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

func (m Model) mapTuiInputsToGitConfig() git.GitConfig {
	return git.GitConfig{
		CommitHash:                      m.inputs[0].Value(),
		URL:                             m.inputs[1].Value(),
		RemoteName:                      m.inputs[2].Value(),
		DevelopBranchName:               m.inputs[3].Value(),
		ReleaseBranchPrependIdentifiers: strings.Split(m.inputs[4].Value(), " "),
	}
}

func (m Model) findFixVersion() tea.Msg {
	gitConfig := m.mapTuiInputsToGitConfig()

	releases := git.FormatRemoteBranches(&gitConfig)

	sortedReleases := git.GetSortedReleases(releases)

	rootCommit := git.GetRootCommit(&gitConfig)

	var message string

	if rootCommit == nil {
		message = "No such hash in the root of this repo"
		return fixVersionMsg(message)
	} else {
		message = "No fixed version found"

		fixedVersions := make([]string, 0)

		for _, version := range sortedReleases {
			if git.IsCommitPresentOnBranch(&gitConfig, rootCommit, releases[version]) {
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
