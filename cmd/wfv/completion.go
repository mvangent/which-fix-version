package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)

var (
	bashCompletion string
)

func NewCompletionCommand() *CompletionCommand {
	fc := &CompletionCommand{
		fs: flag.NewFlagSet("completion", flag.ContinueOnError),
	}
	return fc
}

type CompletionCommand struct {
	fs *flag.FlagSet
}

func (g *CompletionCommand) Name() string {
	return g.fs.Name()
}

func (g *CompletionCommand) Init(args []string) error {
	return g.fs.Parse(args)
}

func (g *CompletionCommand) Run() error {
	data, err := base64.StdEncoding.DecodeString(bashCompletion)
	if err != nil {
		return fmt.Errorf("error: the shell completions script could not be decoded")
	}
	fmt.Fprintf(os.Stdout, "%s\n", string(data))
	return nil
}
