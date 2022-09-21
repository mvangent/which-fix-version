package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

func (m model) findFixVersion() tea.Msg {

	rootCandidates, releases:= getRemoteBranches()

	// fetch commit list from ma(in/ster)
	root := selectRoot(rootCandidates)
	// check latest release

    sortedReleases := getSortedReleases(releases)

    fmt.Println(sortedReleases)
	// if not present
	// return not yet in any release
	// if in release
	// go back one earlier
	// if not present return release
	// if still present go back even further

	// fmt.Printf("Finding commit hash %s ", ch)
	return fixVersionMsg(root)
}

func getSortedReleases(releases map[string]string) []string {
    versions := make([]string, 0)
    for k := range releases {
        versions = append(versions, k)
    }

    sort.Strings(versions)

   return versions 
}

func selectRoot(rootCandidates []string) string {
	// TODO: this should come as default from a flag, lets have, main, master, development fallback
    return rootCandidates[0]
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

	// check if main or master
	for k, v:= range releases {
        fmt.Printf("releases: %s, branchName: %s \n", k, v)
	}

	return rootCandidates, releases
}
