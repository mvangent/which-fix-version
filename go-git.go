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

	//	go m.getRemoteBranches()

	// fetch commit list from ma(in/ster)

	// check latest release
	// if not present
	// return not yet in any release
	// if in release
	// go back one earlier
	// if not present return release
	// if still present go back even further

	// fmt.Printf("Finding commit hash %s ", ch)
	return m.spinner.Tick()
}

func (m model) getRemoteBranches() tea.Msg {

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

	branches := make([]string, 0)

	for _, ref := range refs {
		s := ref.String()
		if strings.Contains(s, "refs/heads/") {
			branches = append(branches, s)
		}
	}

	// check if main or master
	for _, branch := range branches {
		fmt.Println(branch)
	}

	return fixVersionMsg("0.0.1")
}
