package git

import (
	"os/exec"
	"strconv"
	"strings"
)

func FormatLocalBranches(gitConfig *GitConfig) float64 {
	fetchCommand := exec.Command("git", "fetch")
	fetchCommand.Dir = gitConfig.Path
	err := fetchCommand.Run()

	CheckIfError(err)

	prg := "git"

	arg1 := "branch"
	arg2 := "--remotes"
	arg3 := "--contains"
	arg4 := gitConfig.CommitHash

	searchCmd := exec.Command(prg, arg1, arg2, arg3, arg4)
	searchCmd.Dir = gitConfig.Path
	stdout, err := searchCmd.Output()

	CheckIfError(err)

	results := string(stdout[:])

	branchList := strings.Split(results, " ")

	releases := make(map[string]string)

	for _, branchName := range branchList {
		var branchVersion string

		for _, releaseIdentifier := range gitConfig.ReleaseBranchPrependIdentifiers {
			if strings.Contains(branchName, releaseIdentifier) {
				branchVersion = strings.SplitAfter(branchName, releaseIdentifier)[1]
				releases[branchVersion] = branchName
			}
		}
	}

	releaseVersions := make([]float64, len(releases))

	for k := range releases {
		// FIXME: handle release versions with non number characters
		stripped := strings.ReplaceAll(k, "\n", "")

		i, err := strconv.ParseFloat(stripped, 64)

		CheckIfError(err)

		releaseVersions = append(releaseVersions, i)
	}

	firstVersion := 0.0

	for _, v := range releaseVersions {
		if firstVersion == 0.0 {
			firstVersion = v
		}

		if v > 0 && v < firstVersion {
			firstVersion = v
		}
	}

	// FIXME: TRIM back to original length
	return firstVersion
}
