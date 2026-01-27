package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/reservation-v/vlang/internal/bootstrap"
)

type VendorInfo struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type output struct {
	ProjectInfo bootstrap.ProjectInfo `json:"project_info"`
	Vendor      VendorInfo            `json:"vendor"`
}

func WriteOutput(format string, projectInfo bootstrap.ProjectInfo, vendorInfo VendorInfo) error {
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output{projectInfo, vendorInfo})
	case "text":
		printer(projectInfo, vendorInfo)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}

	return nil
}

func printer(projectInfo bootstrap.ProjectInfo, vendorInfo VendorInfo) {
	fmt.Println(
		"Project Info:",
		"\nName:", projectInfo.Name,
		"\nDir:", projectInfo.Dir,
		"\nModulePath:", projectInfo.ModulePath,
		"\nImportPath:", projectInfo.ImportPath,
		"\nVendorStatus", vendorInfo.Status,
	)
}
