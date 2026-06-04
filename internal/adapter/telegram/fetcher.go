package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTimeout = 2 * time.Minute

type FileFetcher struct {
	token      string
	httpClient *http.Client
	timeout    time.Duration
}

func NewFileFetcher(token string) *FileFetcher {
	return &FileFetcher{
		token: token,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		timeout: defaultTimeout,
	}
}

func (f *FileFetcher) Fetch(ctx context.Context, fileID string) (io.ReadCloser, error) {
	ctx, cancel := f.withTimeout(ctx)
	defer cancel()

	filePath, err := f.resolveFilePath(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return f.downloadFile(ctx, filePath)
}

type getFileResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		FilePath string `json:"file_path"`
	} `json:"result"`
}

func (f *FileFetcher) resolveFilePath(ctx context.Context, fileID string) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", f.token, fileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("telegram getFile returned status %d", resp.StatusCode)
	}

	var payload getFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode getFile response: %w", err)
	}

	if !payload.OK || strings.TrimSpace(payload.Result.FilePath) == "" {
		return "", fmt.Errorf("telegram getFile returned empty file path")
	}

	return payload.Result.FilePath, nil
}

func (f *FileFetcher) downloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", f.token, filePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("telegram file download returned status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (f *FileFetcher) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, f.timeout)
}
