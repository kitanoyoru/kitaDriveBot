package fake

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
)

type FileFetcher struct {
	Content string
}

func (f *FileFetcher) Fetch(_ context.Context, _ string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(f.Content)), nil
}

type DriveRepository struct {
	mu            sync.Mutex
	EnsuredPaths  []domain.FolderPath
	UploadedFiles []UploadedFile
}

type UploadedFile struct {
	FolderID string
	Name     string
	Content  string
}

func (d *DriveRepository) EnsureFolderPath(_ context.Context, path domain.FolderPath) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.EnsuredPaths = append(d.EnsuredPaths, path)
	return "folder-" + path.String(), nil
}

func (d *DriveRepository) UploadPDF(_ context.Context, folderID, name string, content io.Reader) (string, error) {
	data, err := io.ReadAll(content)
	if err != nil {
		return "", err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.UploadedFiles = append(d.UploadedFiles, UploadedFile{
		FolderID: folderID,
		Name:     name,
		Content:  string(data),
	})

	return "https://drive.example/" + name, nil
}
