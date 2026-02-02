package inspect

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/reservation-v/vlang/internal/modfile"
)

type Facts struct {
	Dir          string `json:"dir"`
	ModulePath   string `json:"module_path"`
	ImportPath   string `json:"import_path"`
	Name         string `json:"name"`
	GoVersion    string `json:"go_version"`
	HasVendor    bool   `json:"has_vendor"`
	HasGearDir   bool   `json:"has_gear_dir"`
	HasGearRules bool   `json:"has_gear_rules"`
	HasGearSpec  bool   `json:"has_gear_spec"`
}

// func RunInspect(args []string) error {}

func Inspect(dir string) (Facts, error) {
	goModPath := filepath.Join(dir, "go.mod")

	file, readErr := os.ReadFile(goModPath)
	if readErr != nil {
		return Facts{}, fmt.Errorf("read go.mod: %w", readErr)
	}

	modulePath, parseModErr := modfile.ParseModulePath(file)
	if parseModErr != nil {
		return Facts{}, fmt.Errorf("parse module path: %w", parseModErr)
	}

	importPath := modulePath
	importSlice := strings.Split(importPath, "/")
	last := importSlice[len(importSlice)-1]

	var name string
	if len(last) >= 2 {
		if isMajorVersionSegment(last) {
			name = importSlice[len(importSlice)-2]
		} else {
			name = importSlice[len(importSlice)-1]
		}
	} else if len(importSlice) == 1 {
		name = importSlice[0]
	} else {
		return Facts{}, fmt.Errorf("invalid import path: %q", importPath)
	}

	name = strings.ToLower(name)

	goVersion, goParseErr := modfile.ParseGoVersion(file)
	if goParseErr != nil {
		return Facts{}, fmt.Errorf("parse go version: %w", goParseErr)
	}

	hasVendor, hasVendorErr := hasDir(dir, "vendor")
	if hasVendorErr != nil {
		return Facts{}, hasVendorErr
	}

	hasGearDir, hasGearDirErr := hasDir(dir, ".gear")
	if hasGearDirErr != nil {
		return Facts{}, hasGearDirErr
	}

	hasGearRules, hasGearRulesErr := hasFile(dir, filepath.Join(".gear", "rules"))
	if hasGearRulesErr != nil {
		return Facts{}, hasGearRulesErr
	}

	specName := filepath.Join(".gear", name+".spec")
	hasGearSpec, hasSpecErr := hasFile(dir, specName)
	if hasSpecErr != nil {
		return Facts{}, hasSpecErr
	}

	return Facts{
		Dir:          dir,
		ModulePath:   modulePath,
		ImportPath:   importPath,
		Name:         name,
		GoVersion:    goVersion,
		HasVendor:    hasVendor,
		HasGearDir:   hasGearDir,
		HasGearRules: hasGearRules,
		HasGearSpec:  hasGearSpec,
	}, nil
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
