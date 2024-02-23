package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/deepMOORE/tools"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := tools.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload)

	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		log.Println("Auth Detected " + requestPayload.Auth.Email)
		app.authenticate(w, requestPayload.Auth)
	case "log":
		log.Println("Log Detected " + requestPayload.Log.Name)
		app.logItem(w, requestPayload.Log)
	default:
		tools.ErrorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, a LogPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("ERROR ON BUILDING REQUEST: ", err)
		tools.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		log.Println("ERROR ON CLIENT DO: ", err)
		tools.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Println("ERROR STATUS IS NOT ACCEPTED: ", err)
		tools.ErrorJSON(w, err)
		return
	}

	var payload tools.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("Auth action: Unable to create request with " + err.Error())
		tools.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("Auth action: Error response" + err.Error())
		tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		tools.ErrorJSON(w, errors.New("invalid credentials"))
	} else if response.StatusCode != http.StatusAccepted {
		tools.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService tools.JsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		tools.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload tools.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	tools.WriteJSON(w, http.StatusAccepted, payload)
}
