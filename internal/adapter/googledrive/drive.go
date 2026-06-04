package googledrive

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/kitanoyoru/kitaDriveBot/internal/domain"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const defaultTimeout = 2 * time.Minute

type Repository struct {
	client        *drive.Service
	rootFolderID  string
	timeout       time.Duration
	folderCacheMu sync.RWMutex
	folderCache   map[string]string
}

func NewRepository(ctx context.Context, tokenSource oauth2.TokenSource, rootFolderID string) (*Repository, error) {
	client, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("create drive client: %w", err)
	}

	root := strings.TrimSpace(rootFolderID)
	if root == "" {
		root = "root"
	}

	return &Repository{
		client:       client,
		rootFolderID: root,
		timeout:      defaultTimeout,
		folderCache:  make(map[string]string),
	}, nil
}

func (r *Repository) EnsureFolderPath(ctx context.Context, path domain.FolderPath) (string, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	segments := path.Segments()
	if len(segments) == 0 {
		return r.rootFolderID, nil
	}

	parentID := r.rootFolderID
	cacheKey := r.rootFolderID

	for _, segment := range segments {
		cacheKey = cacheKey + "/" + segment

		if folderID, ok := r.cachedFolderID(cacheKey); ok {
			parentID = folderID
			continue
		}

		folderID, err := r.findOrCreateFolder(ctx, parentID, segment)
		if err != nil {
			return "", err
		}

		r.setCachedFolderID(cacheKey, folderID)
		parentID = folderID
	}

	return parentID, nil
}

func (r *Repository) UploadPDF(ctx context.Context, folderID, name string, content io.Reader) (string, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	if strings.TrimSpace(name) == "" {
		name = domain.DefaultPDFName
	}

	file := &drive.File{
		Name:     name,
		MimeType: domain.MIMEPDF,
		Parents:  []string{folderID},
	}

	created, err := r.client.Files.Create(file).
		Media(content, googleapi.ContentType(domain.MIMEPDF)).
		Fields("id, webViewLink").
		Context(ctx).
		Do()
	if err != nil {
		return "", fmt.Errorf("upload pdf: %w", err)
	}

	return created.WebViewLink, nil
}

func (r *Repository) findOrCreateFolder(ctx context.Context, parentID, name string) (string, error) {
	query := buildFolderQuery(parentID, name)

	list, err := r.client.Files.List().
		Q(query).
		Fields("files(id)").
		PageSize(1).
		Spaces("drive").
		Context(ctx).
		Do()
	if err != nil {
		return "", fmt.Errorf("list folder %q: %w", name, err)
	}

	if len(list.Files) > 0 {
		return list.Files[0].Id, nil
	}

	created, err := r.client.Files.Create(&drive.File{
		Name:     name,
		MimeType: domain.FolderMimeType,
		Parents:  []string{parentID},
	}).Fields("id").Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("create folder %q: %w", name, err)
	}

	return created.Id, nil
}

func buildFolderQuery(parentID, name string) string {
	return fmt.Sprintf(
		"name = '%s' and '%s' in parents and mimeType = '%s' and trashed = false",
		escapeQueryValue(name),
		parentID,
		domain.FolderMimeType,
	)
}

func escapeQueryValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	return strings.ReplaceAll(value, `'`, `\'`)
}

func (r *Repository) cachedFolderID(key string) (string, bool) {
	r.folderCacheMu.RLock()
	defer r.folderCacheMu.RUnlock()

	id, ok := r.folderCache[key]
	return id, ok
}

func (r *Repository) setCachedFolderID(key, id string) {
	r.folderCacheMu.Lock()
	defer r.folderCacheMu.Unlock()
	r.folderCache[key] = id
}

func (r *Repository) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, r.timeout)
}
