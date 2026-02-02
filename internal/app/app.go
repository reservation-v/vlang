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
		return cli.RunBootstrap(cmdArgs)
	case "inspect":
		return cli.RunInspect(cmdArgs)
	default:
		return fmt.Errorf("unknown subcommand %q", subCommand)
	}
}
