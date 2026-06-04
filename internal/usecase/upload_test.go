package usecase_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
	"github.com/kitanoyoru/kitaDriveBot/internal/usecase"
	"github.com/kitanoyoru/kitaDriveBot/internal/usecase/mocks"
	"go.uber.org/mock/gomock"
)

func TestUploadService_StorePDF(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fetcher := mocks.NewMockFileFetcher(ctrl)
	drive := mocks.NewMockDriveRepository(ctrl)

	svc := usecase.NewUploadService(fetcher, drive)

	file := domain.IncomingFile{
		ID:       "file-id",
		Name:     "invoice.pdf",
		Size:     100,
		MIMEType: domain.MIMEPDF,
	}
	path := domain.ParseFolderPath("Work/Invoices", "Telegram")

	fetcher.EXPECT().
		Fetch(gomock.Any(), "file-id").
		Return(io.NopCloser(strings.NewReader("pdf-bytes")), nil)

	drive.EXPECT().
		EnsureFolderPath(gomock.Any(), path).
		Return("folder-id", nil)

	drive.EXPECT().
		UploadPDF(gomock.Any(), "folder-id", "invoice.pdf", gomock.Any()).
		Return("https://drive.google.com/file/d/abc/view", nil)

	result, err := svc.StorePDF(context.Background(), file, path)
	if err != nil {
		t.Fatalf("StorePDF() error = %v", err)
	}

	if result.FolderPath != "Work/Invoices" {
		t.Fatalf("FolderPath = %q, want Work/Invoices", result.FolderPath)
	}
	if result.FileName != "invoice.pdf" {
		t.Fatalf("FileName = %q, want invoice.pdf", result.FileName)
	}
	if result.WebViewLink == "" {
		t.Fatal("WebViewLink is empty")
	}
}

func TestUploadService_StorePDF_ValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := usecase.NewUploadService(mocks.NewMockFileFetcher(ctrl), mocks.NewMockDriveRepository(ctrl))

	_, err := svc.StorePDF(
		context.Background(),
		domain.IncomingFile{Name: "image.png", MIMEType: "image/png"},
		domain.ParseFolderPath("", "Telegram"),
	)
	if !errors.Is(err, domain.ErrNotPDF) {
		t.Fatalf("StorePDF() error = %v, want ErrNotPDF", err)
	}
}

func TestUploadService_StorePDF_FetchError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fetcher := mocks.NewMockFileFetcher(ctrl)
	svc := usecase.NewUploadService(fetcher, mocks.NewMockDriveRepository(ctrl))

	file := domain.IncomingFile{
		ID:       "file-id",
		Name:     "invoice.pdf",
		Size:     100,
		MIMEType: domain.MIMEPDF,
	}

	fetcher.EXPECT().
		Fetch(gomock.Any(), "file-id").
		Return(nil, errors.New("network down"))

	_, err := svc.StorePDF(context.Background(), file, domain.ParseFolderPath("", "Telegram"))
	if err == nil {
		t.Fatal("StorePDF() error = nil, want error")
	}
}
