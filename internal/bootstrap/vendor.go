package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Vendor(dir string) (bool, error) {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Dir = dir

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stderr

	var hadVendorBefore bool
	vendorPath := filepath.Join(dir, "vendor")
	info, err := os.Stat(vendorPath)

	if err == nil {
		hadVendorBefore = info.IsDir()
	} else if os.IsNotExist(err) {
		hadVendorBefore = false
	} else {
		return false, fmt.Errorf("vendor check: %w", err)
	}

	err = cmd.Run()
	if err != nil {
		return false, fmt.Errorf("go mod vendor: %w", err)
	}

	info, err = os.Stat(vendorPath)

	if err == nil && info.IsDir() {
		if hadVendorBefore {
			return true, nil
		}
		return false, nil
	}

	return false, fmt.Errorf("vendorDir in %q not found", dir)
}
