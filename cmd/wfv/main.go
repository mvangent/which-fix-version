package main

import (
	"fmt"
	"os"
)

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
}

var (
	availableCommands = "valid options: \n- local         find fix version in local repository \n- remote        scan a remote repo for the fix version"
)

func root(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please choose your subcommand, %s", availableCommands)
	}

	cmds := []Runner{
		NewFindRemoteCommand(),
		NewFindLocalCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s, %s", subcommand, availableCommands)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
