package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type OutputFlags struct {
	Format string
	Output string
}

func addOutputFlags(fs *flag.FlagSet) (formatPtr *string, outputPtr *string) {
	formatPtr = fs.String("format", "json", "how to format project info")
	outputPtr = fs.String("output", "", "write output to file (default: stdout)")
	return formatPtr, outputPtr
}

func openOutputWriter(path string) (w io.Writer, closeFn func() error, existed bool, err error) {
	if path == "" {
		return os.Stdout, func() error { return nil }, false, nil
	}
	info, statErr := os.Stat(path)

	if statErr == nil {
		if info.IsDir() {
			return nil, func() error { return nil }, false,
				fmt.Errorf("cannot write to directory: %q", path)
		}
		existed = true
	} else if os.IsNotExist(statErr) {
		existed = false
	} else {
		return nil, func() error { return nil }, false,
			fmt.Errorf("stat %q: %w", path, statErr)
	}

	file, createErr := os.Create(path)
	if createErr != nil {
		return nil, func() error { return nil }, false,
			fmt.Errorf("create %q: %w", path, createErr)
	}

	w = file

	closeFn = func() error {
		if err := file.Close(); err != nil {
			return fmt.Errorf("close output file: %w", err)
		}
		return nil
	}
	return w, closeFn, existed, nil
}
