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

	writer, closeFn, existed, err := openOutputWriter(bootstrapFlags.Out.Output)

	if err != nil {
		return fmt.Errorf("open output writer: %w", err)
	}

	defer func() {
		closeErr := closeFn()
		if closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
		}
	}()

	err = WriteOutput(writer, bootstrapFlags.Out.Format, projectInfo, vendorInfo)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if bootstrapFlags.Out.Output != "" {
		if existed {
			fmt.Fprintln(os.Stderr, "file was overwritten")
		} else {
			fmt.Fprintln(os.Stderr, "file was created")
		}
	}

	return nil
}

type BootstrapFlags struct {
	Dir    string
	Vendor bool
	Out    OutputFlags
}

func parseBootstrapFlags(args []string) (BootstrapFlags, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	needVendor := fs.Bool("vendor", true, "enable/disable vendoring (true/false)")
	format, output := addOutputFlags(fs)
	if err := fs.Parse(args); err != nil {
		return BootstrapFlags{}, err
	}

	bsFlags := BootstrapFlags{
		Dir:    *dirPtr,
		Vendor: *needVendor,
		Out:    OutputFlags{Format: *format, Output: *output},
	}

	return bsFlags, nil
}
