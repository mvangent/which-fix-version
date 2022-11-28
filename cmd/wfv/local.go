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

func NewFindLocalCommand() *FindLocalCommand {
	fc := &FindLocalCommand{
		fs: flag.NewFlagSet("local", flag.ContinueOnError),
	}

	fc.fs.StringVar(&fc.commitHash, "commitHash", "", "the main/master/development/custom branch commit hash to find the minimal fix version")
	fc.fs.StringVar(&fc.developmentBranchName, "developmentBranchName", "main", "name of the central development branch")
	fc.fs.StringVar(&fc.releaseBranchFormats, "releaseBranchFormats", "", "all string characters in the branchname before the release version. For example: /starproject/ios/releases/ . Take multiple format separate by a space character")
	fc.fs.StringVar(&fc.path, "path", "", "the absolute path to the local git repository")
	fc.fs.BoolVar(&fc.skipFetch, "skipFetch", false, "Set this to true if network calls should be skipped to get the latest references on the local repository")

	return fc
}

type FindLocalCommand struct {
	fs *flag.FlagSet

	developmentBranchName string
	releaseBranchFormats  string
	commitHash            string
	path                  string
	skipFetch             bool
}

func (g *FindLocalCommand) Name() string {
	return g.fs.Name()
}

func (g *FindLocalCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *FindLocalCommand) Run() error {
	app := app.NewApp(&git.GitConfig{
		CommitHash:            g.commitHash,
		Path:                  g.path,
		DevelopmentBranchName: g.developmentBranchName,
		ReleaseBranchFormats:  strings.Split(g.releaseBranchFormats, " "),
		SkipFetch:             g.skipFetch,
	}, tui.Local)

	if err := tea.NewProgram(app.Model).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	return nil
}
