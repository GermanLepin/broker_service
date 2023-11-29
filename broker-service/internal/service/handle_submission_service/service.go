package handle_submission_service

import (
	"broker-service/internal/dto"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type (
	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
		WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
		ReadJSON(w http.ResponseWriter, r *http.Request, data any) error
	}
)

func (s *service) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.RequestPayload

	err := s.jsonService.ReadJSON(w, r, &requestPayload)
	if err != nil {
		s.jsonService.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "authentication":
		s.authenticate(w, requestPayload.AuthenticationPayload)
	case "logging":
		s.loggingItem(w, requestPayload.LoggingPayload)
	default:
		s.jsonService.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (s *service) loggingItem(w http.ResponseWriter, entry dto.LoggingPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	loggerServiceURL := "http://logger-service/log"

	request, err := http.NewRequest(http.MethodPost, loggerServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	var payload dto.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	s.jsonService.WriteJSON(w, http.StatusAccepted, payload)
}

func (s *service) authenticate(w http.ResponseWriter, entry dto.AuthenticationPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	authenticationServiceURL := "http://authentication-service/authenticate"

	request, err := http.NewRequest(http.MethodPost, authenticationServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}
	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		s.jsonService.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		s.jsonService.ErrorJSON(w, errors.New("error calling authentication service"))
		return
	}

	var jsonFromService dto.JsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		s.jsonService.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload dto.JsonResponse
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonFromService.Data

	s.jsonService.WriteJSON(w, http.StatusAccepted, payload)
}

func New(jsonService JsonService) *service {
	return &service{
		jsonService: jsonService,
	}
}

type service struct {
	jsonService JsonService
}
