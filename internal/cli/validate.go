package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/validate"
)

type validateFlags struct {
	Stage string
	Dir   string
	Out   OutputFlags
}

func RunValidate(args []string) error {
	validateFlgs, parseErr := parseValidateFlags(args)
	if parseErr != nil {
		return fmt.Errorf("validate parse flags: %w", parseErr)
	}

	absDir, absErr := absPath(validateFlgs.Dir)
	if absErr != nil {
		return fmt.Errorf("get absolute path: %w", absErr)
	}
	validateFlgs.Dir = absDir

	var report validate.Report
	switch validateFlgs.Stage {
	case "pre":
		report = validate.Pre(absDir)
	default:
		return fmt.Errorf("stage %s not supported", validateFlgs.Stage)
	}

	writer, closeFn, existed, openWriterErr := openOutputWriter(validateFlgs.Out.Output)

	if openWriterErr != nil {
		return fmt.Errorf("open output writer: %w", openWriterErr)
	}

	defer func() {
		closeErr := closeFn()
		if closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
		}
	}()

	writeOutputErr := WriteOutputValidate(writer, validateFlgs.Out.Format, report)
	if writeOutputErr != nil {
		return fmt.Errorf("write output: %w", writeOutputErr)
	}

	if validateFlgs.Out.Output != "" {
		if existed {
			fmt.Fprintln(os.Stderr, "file was overwritten")
		} else {
			fmt.Fprintln(os.Stderr, "file was created")
		}
	}

	return nil
}

func parseValidateFlags(args []string) (validateFlags, error) {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := addDirFlag(fs)
	format, output := addOutputFlags(fs)
	stage := fs.String("stage", "", "stage name")
	if err := fs.Parse(args); err != nil {
		return validateFlags{}, err
	}

	validateFs := validateFlags{
		Dir:   *dirPtr,
		Stage: *stage,
		Out:   OutputFlags{Format: *format, Output: *output},
	}

	return validateFs, nil
}
