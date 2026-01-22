package bootstrap

import (
	"os"
	"path/filepath"
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

func TestInspect(t *testing.T) {
	tests := []struct {
		name       string
		modulePath string
		makeVendor bool
		makeFile   bool // create vendor as a file
		wantName   string
		wantVendor bool
		wantErr    bool
	}{
		{
			name:       "ok_no_vendor",
			modulePath: "github.com/reservation-v/vlang",
			makeVendor: false,
			wantName:   "vlang",
			wantVendor: false,
			wantErr:    false,
		},
		{
			name:       "ok_with_vendor_dir",
			modulePath: "github.com/reservation-v/vlang",
			makeVendor: true,
			wantName:   "vlang",
			wantVendor: true,
			wantErr:    false,
		},
		{
			name:       "ok_semver_import_suffix_v2",
			modulePath: "github.com/reservation-v/vlang/v2",
			makeVendor: false,
			wantName:   "vlang",
			wantVendor: false,
			wantErr:    false,
		},
		{
			name:       "vendor_is_file",
			modulePath: "github.com/reservation-v/vlang",
			makeFile:   true,
			wantName:   "vlang",
			wantVendor: false,
			wantErr:    false,
		},
		{
			name:    "err_no_go_mod",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			if tc.modulePath != "" {
				writeGoMod(t, dir, tc.modulePath)
			}

			if tc.makeVendor {
				if err := os.Mkdir(filepath.Join(dir, "vendor"), 0o755); err != nil {
					t.Fatalf("mkdir vendor: %v", err)
				}
			}

			if tc.makeFile {
				vendorPath := filepath.Join(dir, "vendor")
				if err := os.WriteFile(vendorPath, []byte("not a dir"), 0o644); err != nil {
					t.Fatalf("write vendor file: %v", err)
				}
			}

			got, err := Inspect(dir)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Inspect() error: %v", err)
			}

			if got.Dir != dir {
				t.Fatalf("Dir: got %q, want %q", got.Dir, dir)
			}
			if got.ModulePath != tc.modulePath {
				t.Fatalf("ModulePath: got %q, want %q", got.ModulePath, tc.modulePath)
			}
			if got.ImportPath != tc.modulePath {
				t.Fatalf("ImportPath: got %q, want %q", got.ImportPath, tc.modulePath)
			}
			if got.Name != tc.wantName {
				t.Fatalf("Name: got %q, want %q", got.Name, tc.wantName)
			}
			if got.HasVendor != tc.wantVendor {
				t.Fatalf("HasVendor: got %v, want %v", got.HasVendor, tc.wantVendor)
			}
		})
	}
}
