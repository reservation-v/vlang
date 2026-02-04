package inspect

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/reservation-v/vlang/internal/modfile"
)

type Info struct {
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

func Inspect(dir string) (Info, error) {
	goModPath := filepath.Join(dir, "go.mod")

	file, readErr := os.ReadFile(goModPath)
	if readErr != nil {
		return Info{}, fmt.Errorf("read go.mod: %w", readErr)
	}

	modulePath, parseModErr := modfile.ParseModulePath(file)
	if parseModErr != nil {
		return Info{}, fmt.Errorf("parse module path: %w", parseModErr)
	}

	importPath := modulePath
	name, parseErr := NameFromModulePath(importPath)
	if parseErr != nil {
		return Info{}, fmt.Errorf("parse import path: %w", parseErr)
	}

	goVersion, goParseErr := modfile.ParseGoVersion(file)
	if goParseErr != nil {
		return Info{}, fmt.Errorf("parse go version: %w", goParseErr)
	}

	hasVendor, hasVendorErr := hasDir(dir, "vendor")
	if hasVendorErr != nil {
		return Info{}, hasVendorErr
	}

	hasGearDir, hasGearDirErr := hasDir(dir, ".gear")
	if hasGearDirErr != nil {
		return Info{}, hasGearDirErr
	}

	hasGearRules, hasGearRulesErr := hasFile(dir, filepath.Join(".gear", "rules"))
	if hasGearRulesErr != nil {
		return Info{}, hasGearRulesErr
	}

	specName := filepath.Join(".gear", name+".spec")
	hasGearSpec, hasSpecErr := hasFile(dir, specName)
	if hasSpecErr != nil {
		return Info{}, hasSpecErr
	}

	return Info{
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
