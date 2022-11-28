package git

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// FIXME: split local config with remote one
type GitConfig struct {
	CommitHash            string
	URL                   string
	RemoteName            string
	DevelopmentBranchName string
	ReleaseBranchFormats  []string
	Path                  string
	SkipFetch             bool
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func IsCommitPresentOnBranch(config *GitConfig, rootCommit *object.Commit, branch string) bool {
	result := false

	var r *git.Repository
	var err error
	var ref *plumbing.Reference

	if config.Path != "" {
		r, err = git.PlainOpen(config.Path)

		CheckIfError(err)

		ref, err = r.Reference(plumbing.ReferenceName(branch), true)
	} else {
		r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:           config.URL,
			RemoteName:    config.RemoteName,
			ReferenceName: plumbing.ReferenceName(branch),
			SingleBranch:  true,
		})

		CheckIfError(err)

		// Gets the HEAD history from HEAD, just like this command:
		// ... retrieves the branch pointed by HEAD
		ref, err = r.Head()
	}

	CheckIfError(err)

	// ... retrieves the commit history
	// FIXME: needs to be configurable
	since := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		// FIXME: get to the bottom of isAncestor logic
		isAncestor, parseErr := rootCommit.IsAncestor(c)

		CheckIfError(parseErr)
		if isAncestor {
			result = true
			return nil
		}

		return nil
	})

	CheckIfError(err)

	return result
}

func GetSortedReleases(releases map[string]string) []string {
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

func GetRootCommit(gitConfig *GitConfig) *object.Commit {
	var r *git.Repository
	var err error

	if gitConfig.Path != "" {
		r, err = git.PlainOpen(gitConfig.Path)
	} else {

		// Clones the given repository, creating the remote, the local branches
		// and fetching the objects, everything in memory:
		r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:        gitConfig.URL,
			RemoteName: gitConfig.RemoteName,
			// FIXME: figure out why plumbing is not working
			ReferenceName: plumbing.ReferenceName(strings.Join([]string{"refs/heads", gitConfig.DevelopmentBranchName}, "/")),
			SingleBranch:  true,
		})
	}

	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// FIXME: yes this is hardcoded and will fixed later, but don't want to ddos enterprise repos
	since := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2099, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	var commit *object.Commit
	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		if c.Hash.String() == gitConfig.CommitHash {
			commit = c
			return nil
		}
		return nil
	})

	CheckIfError(err)

	return commit
}

// RemoteRemoteBranches fetches remote branches from the repo origin and filters out the root and release branches
func FormatRemoteBranches(gitConfig *GitConfig) map[string]string {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: gitConfig.RemoteName,
		URLs: []string{gitConfig.URL},
	})

	refs, err := remote.List(&git.ListOptions{})

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	releases := make(map[string]string)

	for _, ref := range refs {
		s := ref.String()
		// FIXME: find a better helper
		if strings.Contains(s, "refs/heads/") {
			branchName := strings.SplitAfter(s, " ")[1]

			var branchVersion string

			for _, releaseIdentifier := range gitConfig.ReleaseBranchFormats {
				if strings.Contains(branchName, releaseIdentifier) {
					branchVersion = strings.SplitAfter(branchName, releaseIdentifier)[1]
					releases[branchVersion] = branchName
				}
			}
		}
	}

	return releases
}
