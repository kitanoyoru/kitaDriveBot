package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
)

//go:generate mockgen -destination=mocks/file_fetcher.go -package=mocks github.com/kitanoyoru/kitaDriveBot/internal/usecase FileFetcher
//go:generate mockgen -destination=mocks/drive_repository.go -package=mocks github.com/kitanoyoru/kitaDriveBot/internal/usecase DriveRepository

type FileFetcher interface {
	Fetch(ctx context.Context, fileID string) (io.ReadCloser, error)
}

type DriveRepository interface {
	EnsureFolderPath(ctx context.Context, path domain.FolderPath) (folderID string, err error)
	UploadPDF(ctx context.Context, folderID, name string, content io.Reader) (webViewLink string, err error)
}

type UploadService struct {
	fetcher FileFetcher
	drive   DriveRepository
}

func NewUploadService(fetcher FileFetcher, drive DriveRepository) *UploadService {
	return &UploadService{
		fetcher: fetcher,
		drive:   drive,
	}
}

type StorePDFResult struct {
	FolderPath  string
	FileName    string
	WebViewLink string
}

func (s *UploadService) StorePDF(ctx context.Context, file domain.IncomingFile, path domain.FolderPath) (StorePDFResult, error) {
	if err := file.Validate(); err != nil {
		return StorePDFResult{}, err
	}

	content, err := s.fetcher.Fetch(ctx, file.ID)
	if err != nil {
		return StorePDFResult{}, fmt.Errorf("fetch file: %w", err)
	}
	defer func() {
		_ = content.Close()
	}()

	folderID, err := s.drive.EnsureFolderPath(ctx, path)
	if err != nil {
		return StorePDFResult{}, fmt.Errorf("ensure folder path: %w", err)
	}

	fileName := file.FileName()
	webViewLink, err := s.drive.UploadPDF(ctx, folderID, fileName, content)
	if err != nil {
		return StorePDFResult{}, fmt.Errorf("upload pdf: %w", err)
	}

	return StorePDFResult{
		FolderPath:  path.String(),
		FileName:    fileName,
		WebViewLink: webViewLink,
	}, nil
}
