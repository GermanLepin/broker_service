package authentication_service

import (
	"authentication-service/internal/dto"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type (
	JsonService interface {
		ErrorJSON(w http.ResponseWriter, err error, status ...int) error
		WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
		ReadJSON(w http.ResponseWriter, r *http.Request, data any) error
	}

	UserRepository interface {
		GetUserByEmail(email string) (*dto.User, error)
		PasswordMatches(plainText string, u dto.User) (bool, error)
	}
)

func (s *service) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload dto.RequestPayload

	err := s.jsonService.ReadJSON(w, r, &requestPayload)
	if err != nil {
		s.jsonService.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := s.userRepository.GetUserByEmail(requestPayload.Email)
	if err != nil {
		s.jsonService.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := s.userRepository.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		s.jsonService.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	err = s.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		s.jsonService.ErrorJSON(w, err)
		return
	}

	payload := dto.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	s.jsonService.WriteJSON(w, http.StatusAccepted, payload)
}

func (s *service) logRequest(name, data string) error {
	var entry dto.LogEntry

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	loggerServiceURL := "http://logger-service/log"

	request, err := http.NewRequest(http.MethodPost, loggerServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil

}

func New(
	jsonService JsonService,
	userRepository UserRepository,
) *service {
	return &service{
		jsonService:    jsonService,
		userRepository: userRepository,
	}
}

type service struct {
	jsonService    JsonService
	userRepository UserRepository
}
