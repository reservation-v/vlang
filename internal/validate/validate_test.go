package validate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
)

func writeGoMod(t *testing.T, dir, modulePath string) {
	t.Helper()

	content := "module " + modulePath + "\n\ngo 1.25\n"
	path := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
}

func findIssue(issues []Issue, code string) *Issue {
	for i := range issues {
		if issues[i].Code == code {
			return &issues[i]
		}
	}
	return nil
}

func TestCheckGoModMissing(t *testing.T) {
	dir := t.TempDir()

	_, issue := checkGoMod(dir)
	if issue == nil {
		t.Fatalf("expected issue, got nil")
	}
	if issue.Code != "GO_MOD_MISSING" {
		t.Fatalf("issue code: got %q, want %q", issue.Code, "GO_MOD_MISSING")
	}
	if !strings.HasSuffix(issue.Path, filepath.Join(dir, "go.mod")) {
		t.Fatalf("issue path: got %q, want suffix %q", issue.Path, filepath.Join(dir, "go.mod"))
	}
}

func TestCheckGoModOK(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/example/project")

	data, issue := checkGoMod(dir)
	if issue != nil {
		t.Fatalf("unexpected issue: %+v", issue)
	}
	if len(data) == 0 {
		t.Fatalf("expected go.mod data, got empty")
	}
}

func TestCheckModuleInvalid(t *testing.T) {
	data := []byte("module github.com/a/b extra\n")
	_, issue := checkModule(data)
	if issue == nil {
		t.Fatalf("expected issue, got nil")
	}
	if issue.Code != "MODULE_PATH_INVALID" {
		t.Fatalf("issue code: got %q, want %q", issue.Code, "MODULE_PATH_INVALID")
	}
}

func TestCheckNameSemver(t *testing.T) {
	name, issue := checkName("github.com/example/project/v2")
	if issue != nil {
		t.Fatalf("unexpected issue: %+v", issue)
	}
	if name != "project" {
		t.Fatalf("name: got %q, want %q", name, "project")
	}
}

func TestCheckWritableGearDirOK(t *testing.T) {
	dir := t.TempDir()
	gearDir := filepath.Join(dir, ".gear")

	if err := os.Mkdir(gearDir, 0o755); err != nil {
		t.Fatalf("mkdir .gear: %v", err)
	}

	issue := checkWritable(dir)
	if issue != nil {
		t.Fatalf("unexpected issue: %+v", issue)
	}
}

func TestCheckWritableGearIsFile(t *testing.T) {
	dir := t.TempDir()
	gearPath := filepath.Join(dir, ".gear")

	if err := os.WriteFile(gearPath, []byte("not a dir"), 0o644); err != nil {
		t.Fatalf("write .gear: %v", err)
	}

	issue := checkWritable(dir)
	if issue == nil {
		t.Fatalf("expected issue, got nil")
	}
	if issue.Code != "GEAR_IS_A_FILE" {
		t.Fatalf("issue code: got %q, want %q", issue.Code, "GEAR_IS_A_FILE")
	}
}

func TestCheckWritableNoGearDirNotWritable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix-only permission behavior")
	}

	dir := t.TempDir()
	origPerm, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat dir: %v", err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatalf("chmod dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(dir, origPerm.Mode().Perm())
	})

	if accessErr := syscall.Access(dir, accessW|accessX); accessErr == nil {
		t.Skip("dir still writable; skipping permission test")
	}

	issue := checkWritable(dir)
	if issue == nil {
		t.Fatalf("expected issue, got nil")
	}
	if issue.Code != "DIR_IS_NOT_ACCESSIBLE" {
		t.Fatalf("issue code: got %q, want %q", issue.Code, "DIR_IS_NOT_ACCESSABLE")
	}
}

func TestMaxSeverity(t *testing.T) {
	cases := []struct {
		name   string
		issues []Issue
		want   Severity
	}{
		{name: "empty", issues: nil, want: SeverityOK},
		{name: "warn", issues: []Issue{{Severity: SeverityWarn}}, want: SeverityWarn},
		{name: "error", issues: []Issue{{Severity: SeverityWarn}, {Severity: SeverityErr}}, want: SeverityErr},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := maxSeverity(tc.issues); got != tc.want {
				t.Fatalf("maxSeverity: got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPreOK(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "github.com/example/project")

	report := Pre(dir)
	if report.Verdict != SeverityOK {
		t.Fatalf("verdict: got %q, want %q", report.Verdict, SeverityOK)
	}
	if report.ModulePath != "github.com/example/project" {
		t.Fatalf("module path: got %q, want %q", report.ModulePath, "github.com/example/project")
	}
	if report.Name != "project" {
		t.Fatalf("name: got %q, want %q", report.Name, "project")
	}
	if len(report.Issues) != 0 {
		t.Fatalf("issues: expected none, got %d", len(report.Issues))
	}
}

func TestPreMissingGoMod(t *testing.T) {
	dir := t.TempDir()

	report := Pre(dir)
	if report.Verdict != SeverityErr {
		t.Fatalf("verdict: got %q, want %q", report.Verdict, SeverityErr)
	}
	if findIssue(report.Issues, "GO_MOD_MISSING") == nil {
		t.Fatalf("expected GO_MOD_MISSING issue")
	}
}
