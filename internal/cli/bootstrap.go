package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/bootstrap"
)

type bootstrapFlags struct {
	Dir    string
	Vendor bool
	Out    OutputFlags
}

func RunBootstrap(args []string) error {
	bootstrapFlgs, err := parseBootstrapFlags(args)
	if err != nil {
		return fmt.Errorf("bootstrap parse flags: %w", err)
	}

	absDir, err := absPath(bootstrapFlgs.Dir)
	if err != nil {
		return err
	}
	bootstrapFlgs.Dir = absDir

	vendorInfo, err := getVendorInfo(bootstrapFlgs.Vendor, bootstrapFlgs.Dir)
	if err != nil {
		return fmt.Errorf("get_vendor info: %w", err)
	}

	projectInfo, err := bootstrap.Inspect(bootstrapFlgs.Dir)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	writer, closeFn, existed, err := openOutputWriter(bootstrapFlgs.Out.Output)

	if err != nil {
		return fmt.Errorf("open output writer: %w", err)
	}

	defer func() {
		closeErr := closeFn()
		if closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
		}
	}()

	err = WriteOutput(writer, bootstrapFlgs.Out.Format, projectInfo, vendorInfo)
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	if bootstrapFlgs.Out.Output != "" {
		if existed {
			fmt.Fprintln(os.Stderr, "file was overwritten")
		} else {
			fmt.Fprintln(os.Stderr, "file was created")
		}
	}

	return nil
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
