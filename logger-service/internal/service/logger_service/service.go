package logger_service

import (
	"logger-service/internal/dto"
	"net/http"
)

type (
	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
		WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
		ReadJSON(w http.ResponseWriter, r *http.Request, data any) error
	}

	LogRepository interface {
		InsertLogEntry(entry dto.LogEntry) error
	}
)

func (s *service) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.JsonPayload
	_ = s.jsonService.ReadJSON(w, r, requestPayload)

	event := dto.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := s.logRepository.InsertLogEntry(event)
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	response := dto.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	s.jsonService.WriteJSON(w, http.StatusAccepted, response)
}

func New(
	jsonService JsonService,
	logRepository LogRepository,
) *service {
	return &service{
		jsonService:   jsonService,
		logRepository: logRepository,
	}
}

type service struct {
	jsonService   JsonService
	logRepository LogRepository
}
