package oauth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kitanoyoru/kitaDriveBot/internal/adapter/tokenstore"
	"github.com/kitanoyoru/kitaDriveBot/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const redirectURL = "http://127.0.0.1:8080/oauth/callback"

type Service struct {
	cfg   config.Config
	store *tokenstore.FileStore
}

func NewService(cfg config.Config) *Service {
	return &Service{
		cfg:   cfg,
		store: tokenstore.NewFileStore(cfg.GoogleTokenPath),
	}
}

func (s *Service) Store() *tokenstore.FileStore {
	return s.store
}

func (s *Service) OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{drive.DriveFileScope},
		Endpoint:     google.Endpoint,
	}
}

func (s *Service) RunInteractiveAuth(ctx context.Context) error {
	oauthCfg := s.OAuthConfig()
	state := fmt.Sprintf("kitadrivebot-%d", time.Now().Unix())

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("invalid oauth state")
			_, _ = fmt.Fprint(w, "Authentication failed: invalid state. You can close this window.")
			return
		}

		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			errCh <- fmt.Errorf("oauth error: %s", errMsg)
			_, _ = fmt.Fprintf(w, "Authentication failed: %s. You can close this window.", errMsg)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("missing oauth code")
			_, _ = fmt.Fprint(w, "Authentication failed: missing code. You can close this window.")
			return
		}

		_, _ = fmt.Fprint(w, "Authentication successful. You can close this window and return to the terminal.")
		codeCh <- code
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("oauth callback server: %w", err)
		}
	}()

	authURL := oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Println("Open this URL in your browser to authorize Google Drive access:")
	fmt.Println(authURL)

	var code string
	select {
	case <-ctx.Done():
		_ = server.Shutdown(context.Background())
		return ctx.Err()
	case err := <-errCh:
		_ = server.Shutdown(context.Background())
		return err
	case code = <-codeCh:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)

	token, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange oauth code: %w", err)
	}

	if token.RefreshToken == "" {
		return fmt.Errorf("no refresh token received; revoke app access in Google Account settings and run auth again")
	}

	if err := s.store.Save(token); err != nil {
		return err
	}

	fmt.Printf("Saved Google token to %s\n", s.cfg.GoogleTokenPath)
	return nil
}

func (s *Service) TokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	token, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	base := s.OAuthConfig().TokenSource(ctx, token)
	return oauth2.ReuseTokenSource(token, &persistingTokenSource{
		source: base,
		store:  s.store,
	}), nil
}

type persistingTokenSource struct {
	source oauth2.TokenSource
	store  *tokenstore.FileStore
}

func (p *persistingTokenSource) Token() (*oauth2.Token, error) {
	token, err := p.source.Token()
	if err != nil {
		return nil, err
	}

	if err := p.store.Save(token); err != nil {
		return nil, fmt.Errorf("persist refreshed token: %w", err)
	}

	return token, nil
}
