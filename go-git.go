package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

func (m model) findFixVersion() tea.Msg {

	branches := getRemoteBranches()

	// fetch commit list from ma(in/ster)
	root := getRoot(branches)
	// check latest release
	// if not present
	// return not yet in any release
	// if in release
	// go back one earlier
	// if not present return release
	// if still present go back even further

	// fmt.Printf("Finding commit hash %s ", ch)
	return fixVersionMsg(root)
}

func getRoot(root map[string]string) string {
	// TODO: this should come as default from a flag, lets have, main, master, development fallback
	if _, exists := root["main"]; exists {
		return "main"
	} else {
		if _, exists := root["master"]; exists {
			return "master"
		} else {
			if _, exists := root["development"]; exists {
				return "development"
			} else {
				panic("no root found")
			}
		}
	}
}

func getRemoteBranches() map[string]string {

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

	branches := map[string]string{}

	for _, ref := range refs {
		s := ref.String()
		if strings.Contains(s, "refs/heads/") {
			branchName := strings.SplitAfter(s, "refs/heads/")[1]

			var branchVersion string

			if strings.Contains(branchName, "release/") {
				branchVersion = strings.SplitAfter(branchName, "release/")[1]
			} else if strings.Contains(branchName, "releases/") {
				branchVersion = strings.SplitAfter(branchName, "releases/")[1]
			} else if strings.Contains(branchName, "release-") {
				branchVersion = strings.SplitAfter(branchName, "release-")[1]
			} else if branchName == "main" || branchName == "master" || branchName == "development" {
				branchVersion = "rootCandidate"
			} else {
				branchVersion = "source"
			}

			branches[branchName] = branchVersion
		}
	}

	// check if main or master
	for _, branch := range branches {
		fmt.Println(branch)
	}

	return branches
}
