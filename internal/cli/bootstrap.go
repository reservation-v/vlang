package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/bootstrap"
)

func bootstrap1(args []string) (string, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	if err := fs.Parse(args); err != nil {
		return "", err
	}

	return *dirPtr, nil
}

func RunBootstrap(args []string) error {
	dir, err := bootstrap1(args)
	if err != nil {
		return err
	}

	err = bootstrap.Vendor(dir)
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
