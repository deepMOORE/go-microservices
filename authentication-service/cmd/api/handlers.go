package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Authentication started")

	var requestPayload struct {
		Email     string `json:"email"`
		Passoword string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		log.Println("Cannot read json with " + err.Error())
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("Error " + err.Error())
		app.errorJSON(w, errors.New("invalid crenedtials"), http.StatusBadRequest)
	}

	valid, err := user.PasswordMatches(requestPayload.Passoword)
	if err != nil || !valid {
		log.Println("Password missmatch " + requestPayload.Email)
		app.errorJSON(w, errors.New("invalid crenedtials"), http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
