package tokenstore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

type FileStore struct {
	path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Load() (*oauth2.Token, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("read token file %s: %w", s.path, err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parse token file %s: %w", s.path, err)
	}

	return &token, nil
}

func (s *FileStore) Save(token *oauth2.Token) error {
	if token == nil {
		return fmt.Errorf("token is nil")
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return fmt.Errorf("create token directory: %w", err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return fmt.Errorf("write token file %s: %w", s.path, err)
	}

	return nil
}

func (s *FileStore) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}
