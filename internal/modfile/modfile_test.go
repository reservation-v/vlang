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
			name:    "ok_inline_comment",
			in:      "module github.com/a/b // comment\n",
			want:    "github.com/a/b",
			wantErr: false,
		},
		{
			name:    "ok_tabs",
			in:      "module\tgithub.com/a/b\n",
			want:    "github.com/a/b",
			wantErr: false,
		},
		{
			name:    "ok_crlf",
			in:      "module github.com/a/b\r\n",
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
			name:    "ok_not_first_line",
			in:      "go 1.22\nrequire example.com/x v1.0.0\nmodule github.com/a/b\n",
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
		{
			name:    "err_comment_only",
			in:      "module // comment\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "err_extra_tokens",
			in:      "module github.com/a/b extra\n",
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

func TestParseGoVersion(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{
			name:    "ok_simple",
			in:      "module github.com/a/b\ngo 1.22\n",
			want:    "1.22",
			wantErr: false,
		},
		{
			name:    "ok_inline_comment",
			in:      "go 1.25 // comment\n",
			want:    "1.25",
			wantErr: false,
		},
		{
			name:    "ok_tabs",
			in:      "go\t1.25\n",
			want:    "1.25",
			wantErr: false,
		},
		{
			name:    "ok_crlf",
			in:      "go 1.25\r\n",
			want:    "1.25",
			wantErr: false,
		},
		{
			name:    "ok_not_first_line",
			in:      "module github.com/a/b\nrequire example.com/x v1.0.0\ngo 1.23\n",
			want:    "1.23",
			wantErr: false,
		},
		{
			name:    "ok_ignores_toolchain",
			in:      "toolchain go1.25.1\ngo 1.24\n",
			want:    "1.24",
			wantErr: false,
		},
		{
			name:    "err_no_go",
			in:      "module github.com/a/b\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "err_empty",
			in:      "go   \n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "err_comment_only",
			in:      "go // comment\n",
			want:    "",
			wantErr: true,
		},
		{
			name:    "err_extra_tokens",
			in:      "go 1.25 extra\n",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGoVersion([]byte(tt.in))

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseGoVersion() want error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseGoVersion() unexpected error: %v", err)
			}

			if got != tt.want {
				t.Errorf("ParseGoVersion() got = %v, want %v", got, tt.want)
			}

		})
	}
}
