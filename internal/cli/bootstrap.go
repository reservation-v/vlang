package cli

import (
	"flag"
	"fmt"
	"io"
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

	var writer io.Writer
	if bootstrapFlags.Output == "" {
		writer = os.Stdout
	} else {
		info, err := os.Stat(bootstrapFlags.Output)

		var exists bool
		if err == nil {
			if info.IsDir() {
				return fmt.Errorf("cannot write to directory")
			}
			exists = true
			fmt.Fprintln(os.Stderr, "file already exists")
		} else if os.IsNotExist(err) {
			exists = false
		} else {
			return fmt.Errorf("os stat: %w", err)
		}

		file, err := os.Create(bootstrapFlags.Output)
		if err != nil {
			return fmt.Errorf("create file: %w", err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintln(os.Stderr, "close output file:", err)
			}
		}()

		writer = file

		if exists {
			fmt.Fprintln(os.Stderr, "file was overwritten")
		} else {
			fmt.Fprintln(os.Stderr, "file was created")
		}

	}

	err = WriteOutput(writer, bootstrapFlags.Format, projectInfo, vendorInfo)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

type BootstrapFlags struct {
	Dir    string
	Vendor bool
	Format string
	Output string
}

func parseBootstrapFlags(args []string) (BootstrapFlags, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	needVendor := fs.Bool("vendor", true, "enable/disable vendoring (true/false)")
	format := fs.String("format", "json", "how to format project info")
	output := fs.String("output", "", "write output to file (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return BootstrapFlags{}, err
	}

	bsFlags := BootstrapFlags{
		Dir:    *dirPtr,
		Vendor: *needVendor,
		Format: *format,
		Output: *output,
	}

	return bsFlags, nil
}
