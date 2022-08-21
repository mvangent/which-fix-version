package main

import (
	"fmt"
	// "strings"
	//"net/http"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	// c "github.com/vpofe/just-in-time/httpclient"
	git "github.com/go-git/go-git/v5"
)

type model struct {
	status string
	err    error
}

var url = "https://charm.sh"

func checkServer() tea.Msg {

	/* response, err := c.HTTP.Get(url)

	   // if err != nil {
			return errMsg(err)
		}
	*/
	// check for release versions
	// Clones the given repository in memory, creating the remote, the local
	// branches and fetching the objects, exactly as:
	Info("git clone https://github.com/go-git/go-billy")

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/go-git/go-billy",
	})

	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	Info("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	CheckIfError(err)

	if err != nil {
		return errMsg(err)
	}

	// releases := strings.Split("--", string(stdout))

	return statusMsg(string(stdout))

	// return statusMsg(response.StatusCode)
}

type statusMsg string

type errMsg error

func (m model) Init() tea.Cmd {
	return checkServer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case statusMsg:
		// The server returned a status message. Save it to our model. Also
		// tell the Bubble Tea runtime we want to exit because we have nothing
		// else to do. We'll still be able to render a final view with our
		// status message.
		m.status = string(msg)
		return m, tea.Quit

	case errMsg:
		// There was an error. Note it in the model. And tell the runtime
		// we're done and want to quit.
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		// Ctrl+c exits. Even with short running programs it's good to have
		// a quit key, just incase your logic is off. Users will be very
		// annoyed if they can't exit.
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// If we happen to get any other messages, don't do anything.
	return m, nil
}

func (m model) View() string {
	// If there's an error, print it out and don't do anything else.
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Tell the user we're doing something.
	s := fmt.Sprintf("Checking %s ... ", url)

	// When the server responds with a status, add it to the current line.
	if len(m.status) > 0 {
		s += fmt.Sprintf("%d %s!", m.status)
	}

	// Send off whatever we came up with above for rendering.
	return "\n" + s + "\n\n"
}

func main() {
	p := tea.NewProgram(model{})

	if err := p.Start(); err != nil {
		fmt.Printf("Something went wrong, ups. Error: %v", err)
		os.Exit(1)
	}
}
