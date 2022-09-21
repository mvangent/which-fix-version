package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

func (m model) findFixVersion() tea.Msg {

	rootCandidates, releases := getRemoteBranches()

	// fetch commit list from ma(in/ster)
	root := selectRoot(rootCandidates)
	// check latest release

	sortedReleases := getSortedReleases(releases)

	fmt.Println(sortedReleases)

	c := getRootCommit(m.commitHash, root)
	// if not present
	// return not yet in any release
	// if in release
	// go back one earlier
	// if not present return release
	// if still present go back even further

	// fmt.Printf("Finding commit hash %s ", ch)
	var message string

	if c == nil {
		message = "No such hash in the root of this repo"
		return fixVersionMsg(message)
	} else {
		message = "No fixed version found"

		for _, version := range sortedReleases {
			if isCommitPresentOnBranch(c, releases[version]) {
				message = version
				break
			}
		}

		return fixVersionMsg(message)
	}
}

func isCommitPresentOnBranch(rootCommit *object.Commit, branch string) bool {
	result := false

	Info("git clone https://github.com/vpofe/just-in-time")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:           "https://github.com/vpofe/just-in-time",
		ReferenceName: plumbing.ReferenceName(branch),
	})

    fmt.Println(branch)

	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	Info("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Message == rootCommit.Message {
			result = true
			return nil
		}

		return nil
	})

	CheckIfError(err)

	return result
}

func getSortedReleases(releases map[string]string) []string {
	versions := make([]string, 0)
	for k := range releases {
		versions = append(versions, k)
	}

	// FIXME: do real semver number sort instead of string alphabetical sort
	sort.Strings(versions)

	for i := len(versions)/2 - 1; i >= 0; i-- {
		opp := len(versions) - 1 - i
		versions[i], versions[opp] = versions[opp], versions[i]
	}

	return versions
}

func selectRoot(rootCandidates []string) string {
	// TODO: this should come as default from a flag, lets have, main, master, development fallback
	return rootCandidates[0]
}

func getRootCommit(hash string, rootBranch string) *object.Commit {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	// FIXME: repo should be stored centrally
	Info("git clone https://github.com/vpofe/just-in-time")
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/vpofe/just-in-time",
	})

	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	Info("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	var commit *object.Commit
	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.Hash.String())
		if c.Hash.String() == hash {
			commit = c
			return nil
		}
		return nil
	})

	CheckIfError(err)

	return commit
}

func getRemoteBranches() ([]string, map[string]string) {

	Info("Get all remote branches")

	time.Sleep(4 * time.Second)

	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/vpofe/just-in-time"},
	})

	refs, err := remote.List(&git.ListOptions{})

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	releases := make(map[string]string)
	rootCandidates := make([]string, 0)

	for _, ref := range refs {
		s := ref.String()
		if strings.Contains(s, "refs/heads/") {
			branchName := strings.SplitAfter(s, "refs/heads/")[1]

			var branchVersion string

			if strings.Contains(branchName, "release/") {
				branchVersion = strings.SplitAfter(branchName, "release/")[1]
				releases[branchVersion] = branchName
			} else if strings.Contains(branchName, "releases/") {
				branchVersion = strings.SplitAfter(branchName, "releases/")[1]
				releases[branchVersion] = branchName
			} else if strings.Contains(branchName, "release-") {
				branchVersion = strings.SplitAfter(branchName, "release-")[1]
				releases[branchVersion] = branchName
			} else if branchName == "main" || branchName == "master" || branchName == "development" {
				rootCandidates = append(rootCandidates, branchName)
			}
		}
	}

	return rootCandidates, releases
}
