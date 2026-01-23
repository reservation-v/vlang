package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/reservation-v/vlang/internal/bootstrap"
)

func parseBootstrapFlags(args []string) (string, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	if err := fs.Parse(args); err != nil {
		return "", err
	}

	return *dirPtr, nil
}

func RunBootstrap(args []string) error {
	dir, err := parseBootstrapFlags(args)
	if err != nil {
		return fmt.Errorf("bootstrap parse flags: %w", err)
	}

	hadVendorBefore, err := bootstrap.Vendor(dir)
	if err != nil {
		return fmt.Errorf("vendor: %w", err)
	}
	if hadVendorBefore {
		fmt.Println("vendor updated")
	} else {
		fmt.Println("vendor created")
	}

	projectInfo, err := bootstrap.Inspect(dir)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	projectInfo.Dir, err = filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("bootstrap absolute path: %w", err)
	}

	fmt.Printf("Project info: %+v\n", projectInfo)

	return nil
}
