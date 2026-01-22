package main

import (
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/app"
)

func run() error {
	return app.Run(os.Args[1:])
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
