package test_test

import (
	"context"
	"testing"

	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
	"github.com/kitanoyoru/kitaDriveBot/internal/usecase"
	"github.com/kitanoyoru/kitaDriveBot/test/fake"
)

// Integration-style test using in-memory fakes. Live Telegram/Google tests belong behind a build tag.
func TestUploadFlowWithFakes(t *testing.T) {
	t.Parallel()

	fetcher := &fake.FileFetcher{Content: "pdf-content"}
	drive := &fake.DriveRepository{}
	svc := usecase.NewUploadService(fetcher, drive)

	file := domain.IncomingFile{
		ID:       "telegram-file-id",
		Name:     "report.pdf",
		Size:     11,
		MIMEType: domain.MIMEPDF,
	}
	path := domain.ParseFolderPath("Work/Reports", "Telegram")

	result, err := svc.StorePDF(context.Background(), file, path)
	if err != nil {
		t.Fatalf("StorePDF() error = %v", err)
	}

	if result.WebViewLink != "https://drive.example/report.pdf" {
		t.Fatalf("WebViewLink = %q", result.WebViewLink)
	}
	if len(drive.EnsuredPaths) != 1 {
		t.Fatalf("EnsuredPaths len = %d, want 1", len(drive.EnsuredPaths))
	}
	if drive.EnsuredPaths[0].String() != "Work/Reports" {
		t.Fatalf("EnsuredPaths[0] = %q", drive.EnsuredPaths[0].String())
	}
	if len(drive.UploadedFiles) != 1 {
		t.Fatalf("UploadedFiles len = %d, want 1", len(drive.UploadedFiles))
	}
	if drive.UploadedFiles[0].Content != "pdf-content" {
		t.Fatalf("UploadedFiles[0].Content = %q", drive.UploadedFiles[0].Content)
	}
}
