package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Vendor(dir string) error {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = dir

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("go mod vendor: %w", err)
	}

	vendorPath := filepath.Join(dir, "vendor")
	info, err := os.Stat(vendorPath)

	if err == nil && info.IsDir() {
		return nil
	}

	return fmt.Errorf("vendor in dir=%q not found", dir)
}
