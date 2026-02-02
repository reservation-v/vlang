package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/inspect"
)

type inspectFlags struct {
	Dir string
	Out OutputFlags
}

func RunInspect(args []string) error {
	inspectFlgs, inspectFlagsErr := parseInspectFlags(args)
	if inspectFlagsErr != nil {
		return fmt.Errorf("inspect parse flags: %w", inspectFlagsErr)
	}

	absDir, absPathErr := absPath(inspectFlgs.Dir)
	if absPathErr != nil {
		return fmt.Errorf("get absolute path: %w", absPathErr)
	}
	inspectFlgs.Dir = absDir

	info, inspectErr := inspect.Inspect(absDir)
	if inspectErr != nil {
		return fmt.Errorf("inspect: %w", inspectErr)
	}

	writer, closeFn, existed, openWriterErr := openOutputWriter(inspectFlgs.Out.Output)
	if openWriterErr != nil {
		return fmt.Errorf("open output writer: %w", openWriterErr)
	}

	defer func() {
		closeErr := closeFn()
		if closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
		}
	}()

	writeOutputErr := WriteOutputInspect(writer, inspectFlgs.Out.Format, info)
	if writeOutputErr != nil {
		return fmt.Errorf("write output: %w", writeOutputErr)
	}

	if inspectFlgs.Out.Output != "" {
		if existed {
			fmt.Fprintln(os.Stderr, "file was overwritten")
		} else {
			fmt.Fprintln(os.Stderr, "file was created")
		}
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
