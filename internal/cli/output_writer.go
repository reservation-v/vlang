package cli

import (
	"fmt"
	"io"
	"os"
)

func writeOutputWriter(output string, writeFn func(writer io.Writer) error) error {
	writer, closeFn, existed, openErr := openOutputWriter(output)
	if openErr != nil {
		return fmt.Errorf("open output writer: %w", openErr)
	}

	defer func() {
		closeErr := closeFn()
		if closeErr != nil {
			fmt.Fprintln(os.Stderr, closeErr)
		}
	}()

	writeErr := writeFn(writer)
	if writeErr != nil {
		return fmt.Errorf("write output: %w", writeErr)
	}

	if output != "" {
		if existed {
			fmt.Fprintln(os.Stderr, "file was overwritten")
		} else {
			fmt.Fprintln(os.Stderr, "file was created")
		}
	}

	return nil
}
