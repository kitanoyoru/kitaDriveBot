package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/config"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Service struct {
	srv *drive.Service

	logger *logger.Logger
}

func NewService(config *config.GoogleDriveConfig, logger *logger.Logger) (*Service, error) {
	s := Service{}

	s.logger = logger

	srv, err := s.initGoogleDriveSrv(config)
	if err != nil {
		return nil, err
	}
	s.srv = srv

	return &s, nil
}

func (s *Service) initGoogleDriveSrv(config *config.GoogleDriveConfig) (*drive.Service, error) {
	b, err := os.ReadFile(config.CredentialsPath)
	if err != nil {
		return nil, err
	}

	cfg, err := google.ConfigFromJSON(b, drive.DriveMetadataScope)
	if err != nil {
		return nil, err
	}

	client, err := s.getClient(cfg, config.TokenPath)
	if err != nil {
		return nil, err
	}

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *Service) getClient(cfg *oauth2.Config, tokenPath string) (*http.Client, error) {
	token, err := s.loadTokenFromFile(tokenPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if token == nil {
		token, err = s.getTokenFromWeb(cfg)
		if err != nil {
			return nil, err
		}

		err = s.saveTokenToFile(token, tokenPath)
		if err != nil {
			return nil, err
		}
	}

	return cfg.Client(context.Background(), token), nil
}

func (s *Service) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) loadTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) saveTokenToFile(token *oauth2.Token, tokenPath string) error {
	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return err
	}

	return nil
}
