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

	absDir, err := filepath.Abs(bootstrapFlags.Dir)
	if err != nil {
		return fmt.Errorf("bootstrap absolute path: %w", err)
	}
	bootstrapFlags.Dir = absDir

	vendorInfo := VendorInfo{}
	if !bootstrapFlags.Vendor {
		vendorInfo.Status = "skipped"
		vendorInfo.Enabled = false
	} else {
		hadVendorBefore, err := bootstrap.Vendor(bootstrapFlags.Dir)
		if err != nil {
			return fmt.Errorf("vendor: %w", err)
		}
		if hadVendorBefore {
			vendorInfo.Status = "updated"
			vendorInfo.Enabled = true
		} else {
			vendorInfo.Status = "created"
			vendorInfo.Enabled = true
		}
	}

	projectInfo, err := bootstrap.Inspect(bootstrapFlags.Dir)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	err = WriteOutput(bootstrapFlags.Format, projectInfo, vendorInfo)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

type BootstrapFlags struct {
	Dir    string
	Vendor bool
	Format string
}

func parseBootstrapFlags(args []string) (BootstrapFlags, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	needVendor := fs.Bool("vendor", true, "enable/disable vendoring (true/false)")
	format := fs.String("format", "json", "how to format project info")
	if err := fs.Parse(args); err != nil {
		return BootstrapFlags{}, err
	}

	bsFlags := BootstrapFlags{
		Dir:    *dirPtr,
		Vendor: *needVendor,
		Format: *format,
	}

	return bsFlags, nil
}
