package broker_service

import (
	"broker-service/internal/dto"
	"net/http"
)

type JsonService interface {
	WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
}

func (s *service) Broker(w http.ResponseWriter, r *http.Request) {
	payload := dto.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = s.jsonService.WriteJSON(w, http.StatusOK, payload)
}

func New(jsonService JsonService) *service {
	return &service{
		jsonService: jsonService,
	}
}

type service struct {
	jsonService JsonService
}
