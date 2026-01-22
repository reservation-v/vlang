package cli

import (
	"flag"
	"os"
)

func Bootstrap(args []string) (string, error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dirPtr := fs.String("dir", ".", "upstream project directory")
	if err := fs.Parse(args); err != nil {
		return "", err
	}

	return *dirPtr, nil
}
