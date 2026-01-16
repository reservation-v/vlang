package modfile

import "testing"

func TestParseModulePath(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{
			name:    "ok_simple",
			in:      "module github.com/a/b\n",
			want:    "github.com/a/b",
			wantErr: false,
		},
		{
			name:    "ok_leading_spaces",
			in:      "   module github.com/a/b\n",
			want:    "github.com/a/b",
			wantErr: false,
		},
		{
			name:    "ok_many_spaces",
			in:      "module     github.com/a/b   \n",
			want:    "github.com/a/b",
			wantErr: false,
		},
		{
			name:    "err_no_module",
			in:      "go 1.22\nrequire example.com/x v1.0.0\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "err_empty_path",
			in:      "module   \n",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseModulePath([]byte(tt.in))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseModulePath() want error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseModulePath() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("ParseModulePath() got = %v, want %v", got, tt.want)
			}

		})
	}

}
