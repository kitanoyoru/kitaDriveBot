package domain

import (
	"path/filepath"
	"strings"
)

type IncomingFile struct {
	ID       string
	Name     string
	Size     int64
	MIMEType string
}

func (f IncomingFile) Validate() error {
	if !f.IsPDF() {
		return ErrNotPDF
	}
	if f.Size > MaxTelegramFileSz {
		return ErrFileTooLarge
	}
	return nil
}

func (f IncomingFile) IsPDF() bool {
	if strings.EqualFold(f.MIMEType, MIMEPDF) {
		return true
	}
	return strings.EqualFold(filepath.Ext(f.Name), ".pdf")
}

func (f IncomingFile) FileName() string {
	name := strings.TrimSpace(f.Name)
	if name == "" {
		return DefaultPDFName
	}
	return name
}

type FolderPath struct {
	value string
}

func ParseFolderPath(raw, fallback string) FolderPath {
	path := strings.Trim(strings.TrimSpace(raw), "/")
	if path == "" {
		path = strings.Trim(strings.TrimSpace(fallback), "/")
	}
	return FolderPath{value: path}
}

func (p FolderPath) String() string {
	return p.value
}

func (p FolderPath) Segments() []string {
	path := strings.Trim(p.value, "/")
	if path == "" {
		return nil
	}

	raw := strings.Split(path, "/")
	segments := make([]string, 0, len(raw))
	for _, segment := range raw {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		segments = append(segments, segment)
	}

	return segments
}
