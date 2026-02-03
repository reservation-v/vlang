package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/reservation-v/vlang/internal/bootstrap"
	"github.com/reservation-v/vlang/internal/inspect"
)

type VendorInfo struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type output struct {
	ProjectInfo bootstrap.ProjectInfo `json:"project_info"`
	Vendor      VendorInfo            `json:"vendor"`
}

func WriteOutputInspect(w io.Writer, format string, info inspect.Info) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	case "text":
		return printInspect(w, info)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

func WriteOutput(w io.Writer, format string, projectInfo bootstrap.ProjectInfo, vendorInfo VendorInfo) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(output{projectInfo, vendorInfo})
	case "text":
		err := printer(w, projectInfo, vendorInfo)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}

	return nil
}

func printer(w io.Writer, projectInfo bootstrap.ProjectInfo, vendorInfo VendorInfo) error {
	_, err := fmt.Fprintln(w,
		"Project Info:",
		"\nName:", projectInfo.Name,
		"\nDir:", projectInfo.Dir,
		"\nModulePath:", projectInfo.ModulePath,
		"\nImportPath:", projectInfo.ImportPath,
		"\nVendorStatus:", vendorInfo.Status,
	)
	if err != nil {
		return fmt.Errorf("printer: %w", err)
	}

	return nil
}

func printInspect(w io.Writer, info inspect.Info) error {
	_, err := fmt.Fprintln(w,
		"Inspect Info",
		"\nDir:", info.Dir,
		"\nModulePath:", info.ModulePath,
		"\nImportPath:", info.ImportPath,
		"\nName:", info.Name,
		"\nGoVersion:", info.GoVersion,
		"\nHasVendor:", info.HasVendor,
		"\nHasGearDir:", info.HasGearDir,
		"\nHasGearRules:", info.HasGearRules,
		"\nHasGearSpec:", info.HasGearSpec,
	)
	if err != nil {
		return fmt.Errorf("inspect printer: %w", err)
	}

	return nil
}
