package bootstrap

import (
	"fmt"

	"github.com/reservation-v/vlang/internal/inspect"
)

type ProjectInfo struct {
	Dir        string `json:"dir"`
	ModulePath string `json:"module_path"`
	ImportPath string `json:"import_path"`
	Name       string `json:"name"`
	HasVendor  bool   `json:"has_vendor"`
}

func Inspect(dir string) (ProjectInfo, error) {
	facts, err := inspect.Inspect(dir)
	if err != nil {
		return ProjectInfo{},
			fmt.Errorf("failed to inspect %s: %w", dir, err)
	}

	return ProjectInfo{
		Dir:        facts.Dir,
		ModulePath: facts.ModulePath,
		ImportPath: facts.ImportPath,
		Name:       facts.Name,
		HasVendor:  facts.HasVendor,
	}, nil

}
