package service

import (
	"authentication-service/internal/dto"
	"errors"
	"fmt"
	"net/http"
)

type JsonService interface {
	ErrorJSON(w http.ResponseWriter, err error, status ...int) error
	WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error
	ReadJSON(w http.ResponseWriter, r *http.Request, data any) error
}

type UserRepository interface {
	GetUserByEmail(email string) (*dto.User, error)
	PasswordMatches(plainText string, u dto.User) (bool, error)
}

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

	fmt.Println(requestPayload.Password)
	fmt.Println(user.Password)

	valid, err := s.userRepository.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		s.jsonService.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := dto.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	fmt.Println(payload)

	s.jsonService.WriteJSON(w, http.StatusAccepted, payload)
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
