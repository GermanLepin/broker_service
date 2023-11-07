package routes

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type AuthenticateService interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

func (s *service) Routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/authenticate", s.authenticateService.Authenticate)

	return mux
}

func New(authenticateService AuthenticateService) *service {
	return &service{
		authenticateService: authenticateService,
	}
}

type service struct {
	authenticateService AuthenticateService
}
