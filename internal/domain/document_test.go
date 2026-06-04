package domain_test

import (
	"errors"
	"testing"

	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
)

func TestIncomingFile_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    domain.IncomingFile
		wantErr error
	}{
		{
			name: "valid pdf by mime",
			file: domain.IncomingFile{
				ID:       "1",
				Name:     "invoice.pdf",
				Size:     1024,
				MIMEType: domain.MIMEPDF,
			},
		},
		{
			name: "valid pdf by extension",
			file: domain.IncomingFile{
				ID:   "1",
				Name: "invoice.PDF",
				Size: 1024,
			},
		},
		{
			name: "not pdf",
			file: domain.IncomingFile{
				ID:       "1",
				Name:     "image.png",
				Size:     1024,
				MIMEType: "image/png",
			},
			wantErr: domain.ErrNotPDF,
		},
		{
			name: "too large",
			file: domain.IncomingFile{
				ID:       "1",
				Name:     "big.pdf",
				Size:     domain.MaxTelegramFileSz + 1,
				MIMEType: domain.MIMEPDF,
			},
			wantErr: domain.ErrFileTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.file.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestIncomingFile_FileName(t *testing.T) {
	t.Parallel()

	if got := (domain.IncomingFile{}).FileName(); got != domain.DefaultPDFName {
		t.Fatalf("FileName() = %q, want %q", got, domain.DefaultPDFName)
	}

	if got := (domain.IncomingFile{Name: "  report.pdf  "}).FileName(); got != "report.pdf" {
		t.Fatalf("FileName() = %q, want report.pdf", got)
	}
}

func TestParseFolderPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		raw      string
		fallback string
		want     string
		segments []string
	}{
		{
			name:     "uses raw path",
			raw:      " Work/Invoices/2026 ",
			fallback: "Telegram",
			want:     "Work/Invoices/2026",
			segments: []string{"Work", "Invoices", "2026"},
		},
		{
			name:     "uses fallback",
			raw:      "",
			fallback: "Telegram",
			want:     "Telegram",
			segments: []string{"Telegram"},
		},
		{
			name:     "trims slashes",
			raw:      "/Receipts/June/",
			fallback: "Telegram",
			want:     "Receipts/June",
			segments: []string{"Receipts", "June"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := domain.ParseFolderPath(tt.raw, tt.fallback)
			if path.String() != tt.want {
				t.Fatalf("String() = %q, want %q", path.String(), tt.want)
			}

			segments := path.Segments()
			if len(segments) != len(tt.segments) {
				t.Fatalf("Segments() = %v, want %v", segments, tt.segments)
			}
			for i := range tt.segments {
				if segments[i] != tt.segments[i] {
					t.Fatalf("Segments()[%d] = %q, want %q", i, segments[i], tt.segments[i])
				}
			}
		})
	}
}
