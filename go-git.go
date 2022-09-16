package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

func checkServer() tea.Msg {

	// check for release versions
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	Info("git clone the test repo")

	// Create the remote with repository URL
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/vpofe/just-in-time"},
	})

	// Gets the HEAD history from HEAD, just like this command:
	Info("git list all remote branches")

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, ref := range refs {
		fmt.Println(ref.String())
	}

	// ... just iterates over the commits, printing it
	CheckIfError(err)

	if err != nil {
		return errMsg(err)
	}

	// releases := strings.Split("--", string(stdout))

	return statusMsg(string("great success"))

	// return statusMsg(response.StatusCode)
}
