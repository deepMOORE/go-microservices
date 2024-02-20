package main

import (
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

	payload := tools.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	tools.WriteJSON(w, http.StatusAccepted, payload)
}
