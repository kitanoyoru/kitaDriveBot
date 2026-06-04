package telegram

import (
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
)

func TestIsOwner(t *testing.T) {
	t.Parallel()

	if IsOwner(nil, 42) {
		t.Fatal("IsOwner(nil) = true, want false")
	}

	if !IsOwner(&models.User{ID: 42}, 42) {
		t.Fatal("IsOwner(owner) = false, want true")
	}

	if IsOwner(&models.User{ID: 7}, 42) {
		t.Fatal("IsOwner(other) = true, want false")
	}
}

func TestMapDocument(t *testing.T) {
	t.Parallel()

	file := MapDocument(&models.Document{
		FileID:   "abc",
		FileName: "report.pdf",
		FileSize: 100,
		MimeType: domain.MIMEPDF,
	})

	if file.ID != "abc" {
		t.Fatalf("ID = %q, want abc", file.ID)
	}
	if file.Name != "report.pdf" {
		t.Fatalf("Name = %q, want report.pdf", file.Name)
	}
	if file.Size != 100 {
		t.Fatalf("Size = %d, want 100", file.Size)
	}
	if file.MIMEType != domain.MIMEPDF {
		t.Fatalf("MIMEType = %q, want %q", file.MIMEType, domain.MIMEPDF)
	}
}
