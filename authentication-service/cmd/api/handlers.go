package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/deepMOORE/tools"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Authentication started")

	var requestPayload struct {
		Email     string `json:"email"`
		Passoword string `json:"password"`
	}

	err := tools.ReadJSON(w, r, &requestPayload)

	if err != nil {
		log.Println("Cannot read json with " + err.Error())
		tools.ErrorJSON(w, err, http.StatusBadRequest)
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("Error " + err.Error())
		tools.ErrorJSON(w, errors.New("invalid crenedtials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Passoword)
	if err != nil || !valid {
		log.Println("Password missmatch " + requestPayload.Email)
		tools.ErrorJSON(w, errors.New("invalid crenedtials"), http.StatusBadRequest)
	}

	_ = app.logRequest("authenticate", fmt.Sprintf("%s logged in", user.Email))

	payload := tools.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	tools.WriteJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name string, data string) error {
	var entity struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entity.Name = name
	entity.Data = data

	jsonData, _ := json.MarshalIndent(entity, "", "\t")
	logServiceUrl := "http://logger-service/log"

	reqeust, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(reqeust)

	if err != nil {
		return err
	}

	return nil
}
