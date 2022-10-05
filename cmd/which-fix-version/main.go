package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpofe/which-fix-version/app"
)

func main() {
	app := app.NewApp()
	if err := tea.NewProgram(app.Model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
