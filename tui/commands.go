package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpofe/which-fix-version/git"
)

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
