package main

import (
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/modfile"
)

func run() error {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("read go.mod: %w", err)
	}

	modulePath, err := modfile.ParseModulePath(data)
	if err != nil {
		return fmt.Errorf("parse module path: %w", err)
	}

	fmt.Println(modulePath)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
