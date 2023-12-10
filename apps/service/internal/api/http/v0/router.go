package v0

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/service"
)

const (
	apiPrefix = "/api/v0"

	apiBaseRoutesPrefix   = "/"
	apiPersonRoutesPrefix = "/person"
)

type HTTPApi struct {
	service *service.Service
}

func NewHTTPApi(service *service.Service) *HTTPApi {
	return &HTTPApi{
		service,
	}
}

func (api *HTTPApi) GetRouter() (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	return r, nil
}
