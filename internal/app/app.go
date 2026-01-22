package app

import (
	"fmt"

	"github.com/reservation-v/vlang/internal/bootstrap"
	"github.com/reservation-v/vlang/internal/cli"
)

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command given")
	}
	subCommand := args[0]
	cmdArgs := args[1:]

	switch subCommand {
	case "bootstrap":
		err := runBootstrap(cmdArgs)
		if err != nil {
			return err
		}
	default:
		fmt.Printf("unknown subcommand %q\n", subCommand)
	}

	return nil
}

func runBootstrap(args []string) error {
	dir, err := cli.Bootstrap(args)
	if err != nil {
		return err
	}

	projectInfo, err := bootstrap.Inspect(dir)
	if err != nil {
		return err
	}
	fmt.Printf("Project info: %+v\n", projectInfo)

	return nil
}
