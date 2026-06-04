package tokenstore_test

import (
	"testing"

	"github.com/kitanoyoru/kitaDriveBot/internal/adapter/tokenstore"
	"golang.org/x/oauth2"
)

func TestFileStore_SaveLoadExists(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := dir + "/token.json"
	store := tokenstore.NewFileStore(path)

	if store.Exists() {
		t.Fatal("Exists() = true before save")
	}

	token := &oauth2.Token{
		AccessToken:  "access",
		RefreshToken: "refresh",
		TokenType:    "Bearer",
	}

	if err := store.Save(token); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !store.Exists() {
		t.Fatal("Exists() = false after save")
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Fatalf("AccessToken = %q, want %q", loaded.AccessToken, token.AccessToken)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Fatalf("RefreshToken = %q, want %q", loaded.RefreshToken, token.RefreshToken)
	}
}
