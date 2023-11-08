package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action      string                `json:"action"`
	AuthPayload AuthenticationPayload `json:"authentication,omitempty"`
}

type AuthenticationPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = c.writeJSON(w, http.StatusOK, payload)
}

func (c *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := c.readJSON(w, r, &requestPayload)
	if err != nil {
		c.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "authentication":
		c.authenticate(w, requestPayload.AuthPayload)
	default:
		c.errorJSON(w, errors.New("unknown action"))
	}
}

func (c *Config) authenticate(w http.ResponseWriter, a AuthenticationPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		c.errorJSON(w, err)
		return
	}
	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		c.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		c.errorJSON(w, errors.New("error calling authentication service"))
		return
	}

	// create a varabiel we'll read response.Body into
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		c.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		c.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	c.writeJSON(w, http.StatusAccepted, payload)
}
