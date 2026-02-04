package inspect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func NameFromModulePath(modulePath string) (string, error) {
	moduleSlice := strings.Split(modulePath, "/")
	last := moduleSlice[len(moduleSlice)-1]

	var name string
	if len(last) >= 2 {
		if isMajorVersionSegment(last) {
			name = moduleSlice[len(moduleSlice)-2]
		} else {
			name = moduleSlice[len(moduleSlice)-1]
		}
	} else if len(moduleSlice) == 1 {
		name = moduleSlice[0]
	} else {
		return "", fmt.Errorf("invalid module path: %q", modulePath)
	}

	name = strings.ToLower(name)
	return name, nil
}

func isMajorVersionSegment(s string) bool {
	if len(s) < 2 || s[0] != 'v' {
		return false
	}
	for _, r := range s[1:] {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func hasDir(path, dir string) (bool, error) {
	dirPath := filepath.Join(path, dir)
	info, err := os.Stat(dirPath)
	var hasDir bool

	if err == nil {
		hasDir = info.IsDir()
	} else if os.IsNotExist(err) {
		hasDir = false
	} else {
		return false, fmt.Errorf("stat: %w", err)
	}
	return hasDir, nil
}

func hasFile(path, dir string) (bool, error) {
	dirPath := filepath.Join(path, dir)
	info, err := os.Stat(dirPath)
	var hasFile bool

	if err == nil {
		hasFile = !info.IsDir()
	} else if os.IsNotExist(err) {
		hasFile = false
	} else {
		return false, fmt.Errorf("stat: %w", err)
	}
	return hasFile, nil
}
