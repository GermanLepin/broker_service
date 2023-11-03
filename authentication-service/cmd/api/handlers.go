package main

import (
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := c.readJSON(w, r, &requestPayload)
	if err != nil {
		c.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := c.Models.User.GetUserByEmail(requestPayload.Email)
	if err != nil {
		c.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	fmt.Println(requestPayload.Password)
	fmt.Println(user.Password)

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		c.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	fmt.Println(payload)

	c.writeJSON(w, http.StatusAccepted, payload)
}
