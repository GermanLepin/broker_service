package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type (
	BrokerService interface {
		Broker(w http.ResponseWriter, r *http.Request)
	}

	HandleSubmissionService interface {
		HandleSubmission(w http.ResponseWriter, r *http.Request)
	}
)

func (s *service) NewRoutes() http.Handler {
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

	mux.Post("/", s.brokerService.Broker)
	mux.Post("/handle", s.handleSubmissionService.HandleSubmission)

	return mux
}

func New(
	brokerService BrokerService,
	handleSubmissionService HandleSubmissionService,
) *service {
	return &service{
		brokerService:           brokerService,
		handleSubmissionService: handleSubmissionService,
	}
}

type service struct {
	brokerService           BrokerService
	handleSubmissionService HandleSubmissionService
}
