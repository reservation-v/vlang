package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/bootstrap"
)

func RunBootstrap(args []string) error {
	bootstrapFlags, err := parseBootstrapFlags(args)
	if err != nil {
		return fmt.Errorf("bootstrap parse flags: %w", err)
	}

	absDir, err := absPath(bootstrapFlags.Dir)
	if err != nil {
		return err
	}
	bootstrapFlags.Dir = absDir

	vendorInfo, getVendorErr := getVendorInfo(bootstrapFlags.Vendor, bootstrapFlags.Dir)
	if getVendorErr != nil {
		return fmt.Errorf("get_vendor info: %w", getVendorErr)
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

type bootstrapFlags struct {
	Dir    string
	Vendor bool
	Out    OutputFlags
}

func parseBootstrapFlags(args []string) (bootstrapFlags, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := addDirFlag(fs)
	needVendor := fs.Bool("vendor", true, "enable/disable vendoring (true/false)")
	format, output := addOutputFlags(fs)
	if err := fs.Parse(args); err != nil {
		return bootstrapFlags{}, err
	}

	bsFlags := bootstrapFlags{
		Dir:    *dirPtr,
		Vendor: *needVendor,
		Out:    OutputFlags{Format: *format, Output: *output},
	}

	return bsFlags, nil
}

func getVendorInfo(needVendor bool, dir string) (VendorInfo, error) {
	vendorInfo := VendorInfo{}
	if !needVendor {
		vendorInfo.Status = "skipped"
		vendorInfo.Enabled = false
	} else {
		hadVendorBefore, err := bootstrap.Vendor(dir)
		if err != nil {
			return VendorInfo{}, fmt.Errorf("vendor: %w", err)
		}
		if hadVendorBefore {
			vendorInfo.Status = "updated"
			vendorInfo.Enabled = true
		} else {
			vendorInfo.Status = "created"
			vendorInfo.Enabled = true
		}
	}

	return vendorInfo, nil
}
