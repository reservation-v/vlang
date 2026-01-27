package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/reservation-v/vlang/internal/modfile"
)

type ProjectInfo struct {
	Dir        string `json:"dir"`
	ModulePath string `json:"module_path"`
	ImportPath string `json:"import_path"`
	Name       string `json:"name"`
	HasVendor  bool   `json:"has_vendor"`
}

func Inspect(dir string) (ProjectInfo, error) {
	goModPath := filepath.Join(dir, "go.mod")

	file, err := os.ReadFile(goModPath)
	if err != nil {
		return ProjectInfo{}, fmt.Errorf("read go.mod: %w", err)
	}

	modulePath, err := modfile.ParseModulePath(file)
	if err != nil {
		return ProjectInfo{}, fmt.Errorf("parse module path: %w", err)
	}

	importPath := modulePath
	importSlice := strings.Split(importPath, "/")
	last := importSlice[len(importSlice)-1]

	var name string
	if len(last) >= 2 {
		if last[0] == 'v' && unicode.IsDigit(rune(importSlice[len(importSlice)-1][1])) {
			name = importSlice[len(importSlice)-2]
		} else {
			name = importSlice[len(importSlice)-1]
		}
	} else if len(importSlice) == 1 {
		name = importSlice[0]
	}

	var hasVendor bool

	vendorPath := filepath.Join(dir, "vendor")
	info, err := os.Stat(vendorPath)

	if err == nil {
		hasVendor = info.IsDir()
	} else if os.IsNotExist(err) {
		hasVendor = false
	} else {
		return ProjectInfo{}, fmt.Errorf("stat vendor: %w", err)
	}

	return ProjectInfo{
		Dir:        dir,
		ModulePath: modulePath,
		ImportPath: importPath,
		Name:       name,
		HasVendor:  hasVendor,
	}, nil

}
