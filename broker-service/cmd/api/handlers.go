package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}

	_ = app.WriteJson(w, http.StatusOK, payload)
}

func (app *Config) HandlerSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default: 
		app.ErrorJson(w, errors.New("unknown action"))

	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsondata, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authenticate-service/authenticate", bytes.NewBuffer(jsondata))
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJson(w, errors.New("invalid credentials"))
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.ErrorJson(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	if jsonFromService.Error {
		app.ErrorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.WriteJson(w, http.StatusAccepted, payload)

}
