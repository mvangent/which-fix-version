package main

import (
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	c "github.com/vpofe/just-in-time/httpclient"
)

type model struct {
	status int
	err    error
}

var url = "https://charm.sh"

func checkServer() tea.Msg {
	response, err := c.HTTP.Get(url)

	if err != nil {
		return errMsg(err)
	}

	return statusMsg(response.StatusCode)
}

type statusMsg int

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
		m.status = int(msg)
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
	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
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
