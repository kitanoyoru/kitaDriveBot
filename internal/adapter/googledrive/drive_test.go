package googledrive

import "testing"

func TestBuildFolderQuery(t *testing.T) {
	t.Parallel()

	query := buildFolderQuery("root", "Work/Invoices")
	want := "name = 'Work/Invoices' and 'root' in parents and mimeType = 'application/vnd.google-apps.folder' and trashed = false"
	if query != want {
		t.Fatalf("buildFolderQuery() = %q, want %q", query, want)
	}
}

func TestEscapeQueryValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain", input: "Work", want: "Work"},
		{name: "quote", input: "O'Brien", want: `O\'Brien`},
		{name: "backslash", input: `path\name`, want: `path\\name`},
		{name: "both", input: `O'Brien\docs`, want: `O\'Brien\\docs`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := escapeQueryValue(tt.input); got != tt.want {
				t.Fatalf("escapeQueryValue() = %q, want %q", got, tt.want)
			}
		})
	}
}
