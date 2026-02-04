package validate

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/reservation-v/vlang/internal/inspect"
	"github.com/reservation-v/vlang/internal/modfile"
)

const (
	accessW = 0x2
	accessX = 0x1
)

type Severity string

const (
	SeverityOK   Severity = "OK"
	SeverityWarn Severity = "WARN"
	SeverityErr  Severity = "ERROR"
)

type Issue struct {
	Severity Severity `json:"severity"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	Path     string   `json:"path"`
}

type Report struct {
	Stage      string   `json:"stage"`
	Verdict    Severity `json:"verdict"`
	Issues     []Issue  `json:"issues"`
	ModulePath string   `json:"module_path"`
	Name       string   `json:"name"`
}

func checkGoMod(dir string) (data []byte, issue *Issue) {
	goModPath := filepath.Join(dir, "go.mod")
	_, statErr := os.Stat(goModPath)
	file, readErr := os.ReadFile(goModPath)

	if os.IsNotExist(statErr) {
		return file, &Issue{
			Severity: SeverityErr,
			Code:     "GO_MOD_MISSING",
			Message:  "go.mod file is missing",
			Path:     goModPath,
		}
	}

	if readErr != nil {
		return nil, &Issue{
			Severity: SeverityErr,
			Code:     "GO_MOD_READ_FAILED",
			Message:  "go.mod file cannot be read",
			Path:     goModPath,
		}
	}

	return file, nil
}

func checkModule(data []byte) (modulePath string, issue *Issue) {
	modulePath, err := modfile.ParseModulePath(data)
	if err != nil {
		return "", &Issue{
			Severity: SeverityErr,
			Code:     "MODULE_PATH_INVALID",
			Message:  "invalid module path",
			Path:     modulePath,
		}
	}

	return modulePath, nil
}

func checkName(modulePath string) (name string, issue *Issue) {
	name, err := inspect.NameFromModulePath(modulePath)
	if err != nil {
		return "", &Issue{
			Severity: SeverityErr,
			Code:     "MODULE_PATH_INVALID",
			Message:  "invalid module path",
			Path:     modulePath,
		}
	}

	return name, nil
}

func checkWritable(dir string) *Issue {
	gearDir := filepath.Join(dir, ".gear")
	info, err := os.Stat(gearDir)

	if err == nil {
		if !info.IsDir() {
			return &Issue{
				Severity: SeverityErr,
				Code:     "GEAR_IS_A_FILE",
				Message:  ".gear is a file",
				Path:     gearDir,
			}
		}
		if accessErr := syscall.Access(gearDir, accessW|accessX); accessErr != nil {
			return &Issue{
				Severity: SeverityErr,
				Code:     "GEAR_IS_NOT_ACCESSIBLE",
				Message:  ".gear is not accessible",
				Path:     gearDir,
			}
		}
		return nil
	}

	if os.IsNotExist(err) {
		accessErr := syscall.Access(dir, accessW|accessX)
		if accessErr != nil {
			return &Issue{
				Severity: SeverityErr,
				Code:     "DIR_IS_NOT_ACCESSIBLE",
				Message:  "src dir is not accessible",
				Path:     dir,
			}
		}
		return nil
	}

	return &Issue{
		Severity: SeverityErr,
		Code:     "OS_STAT_FAILED",
		Message:  "os stat method failed",
		Path:     gearDir,
	}
}

func maxSeverity(issues []Issue) Severity {
	maxSeverity := SeverityOK
	for _, issue := range issues {
		warn := issue.Severity == SeverityWarn
		err := issue.Severity == SeverityErr
		if warn {
			maxSeverity = SeverityWarn
			continue
		} else if err {
			maxSeverity = SeverityErr
			break
		}
	}

	return maxSeverity
}
