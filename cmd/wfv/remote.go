package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vpofe/which-fix-version/app"
	"github.com/vpofe/which-fix-version/git"
	"github.com/vpofe/which-fix-version/tui"
)

func NewFindRemoteCommand() *FindRemoteCommand {
	fc := &FindRemoteCommand{
		fs: flag.NewFlagSet("remote", flag.ContinueOnError),
	}

	fc.fs.StringVar(&fc.commitHash, "commitHash", "", "the main/master/development/custom branch commit hash to find the minimal fix version")
	fc.fs.StringVar(&fc.url, "url", "", "git repository url")
	fc.fs.StringVar(&fc.remoteName, "remoteName", "origin", "remote name to fetch branches from")
	fc.fs.StringVar(&fc.developmentBranchName, "developmentBranchName", "main", "name of the central development branch")
	fc.fs.StringVar(&fc.releaseBranchFormats, "releaseBranchFormats", "", "all string characters in the branchname before the release version. For example: /starproject/ios/releases/. Takes multiple formats seperated by a space character")

	return fc
}

type FindRemoteCommand struct {
	fs *flag.FlagSet

	url                   string
	remoteName            string
	developmentBranchName string
	releaseBranchFormats  string
	commitHash            string
}

func (g *FindRemoteCommand) Name() string {
	return g.fs.Name()
}

func (g *FindRemoteCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *FindRemoteCommand) Run() error {
	app := app.NewApp(&git.GitConfig{
		CommitHash:            g.commitHash,
		URL:                   g.url,
		RemoteName:            g.remoteName,
		DevelopmentBranchName: g.developmentBranchName,
		ReleaseBranchFormats:  strings.Split(g.releaseBranchFormats, " "),
	}, tui.Remote)

	if err := tea.NewProgram(app.Model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	return nil
}
