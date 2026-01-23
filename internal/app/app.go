package app

import (
	"fmt"

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
		err := cli.RunBootstrap(cmdArgs)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown subcommand %q", subCommand)
	}

	return nil
}
