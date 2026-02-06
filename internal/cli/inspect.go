package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/reservation-v/vlang/internal/inspect"
)

type inspectFlags struct {
	Dir string
	Out OutputFlags
}

func RunInspect(args []string) error {
	inspectFlgs, parseErr := parseInspectFlags(args)
	if parseErr != nil {
		return fmt.Errorf("inspect parse flags: %w", parseErr)
	}

	absDir, absErr := absPath(inspectFlgs.Dir)
	if absErr != nil {
		return fmt.Errorf("get absolute path: %w", absErr)
	}
	inspectFlgs.Dir = absDir

	info, inspectErr := inspect.Inspect(absDir)
	if inspectErr != nil {
		return fmt.Errorf("inspect: %w", inspectErr)
	}

	writeErr := writeOutputWriter(inspectFlgs.Out.Output, func(w io.Writer) error {
		return WriteOutputInspect(w, inspectFlgs.Out.Format, info)
	})
	if writeErr != nil {
		return fmt.Errorf("write output: %w", writeErr)
	}

	return nil
}

func parseInspectFlags(args []string) (inspectFlags, error) {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := addDirFlag(fs)
	format, output := addOutputFlags(fs)
	if err := fs.Parse(args); err != nil {
		return inspectFlags{}, err
	}

	inspectFs := inspectFlags{
		Dir: *dirPtr,
		Out: OutputFlags{Format: *format, Output: *output},
	}

	return inspectFs, nil
}
