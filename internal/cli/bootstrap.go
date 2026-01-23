package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/reservation-v/vlang/internal/bootstrap"
)

func RunBootstrap(args []string) error {
	bootstrapFlags, err := parseBootstrapFlags(args)
	if err != nil {
		return fmt.Errorf("bootstrap parse flags: %w", err)
	}

	if bootstrapFlags.Vendor == true {
		hadVendorBefore, err := bootstrap.Vendor(bootstrapFlags.Dir)
		if err != nil {
			return fmt.Errorf("vendor: %w", err)
		}
		if hadVendorBefore {
			fmt.Println("vendor updated")
		} else {
			fmt.Println("vendor created")
		}
	} else if bootstrapFlags.Vendor != false {
		return fmt.Errorf("unknown vendor flag: %q", bootstrapFlags.Vendor)
	}

	projectInfo, err := bootstrap.Inspect(bootstrapFlags.Dir)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	projectInfo.Dir, err = filepath.Abs(bootstrapFlags.Dir)
	if err != nil {
		return fmt.Errorf("bootstrap absolute path: %w", err)
	}

	fmt.Printf("Project info: %+v\n", projectInfo)

	return nil
}

type BootstrapFlags struct {
	Dir    string
	Vendor bool
}

func parseBootstrapFlags(args []string) (BootstrapFlags, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	needVendor := fs.Bool("vendor", true, "enable/disable vendoring (yes/no)")
	if err := fs.Parse(args); err != nil {
		return BootstrapFlags{}, err
	}

	bsFlags := BootstrapFlags{
		Dir:    *dirPtr,
		Vendor: *needVendor,
	}

	return bsFlags, nil
}
