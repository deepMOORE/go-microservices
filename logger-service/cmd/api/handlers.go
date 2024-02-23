package main

import (
	"log"
	"log-service/data"
	"net/http"

	"github.com/deepMOORE/tools"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	log.Print("Starting write log...")
	var requestPayload JSONPayload

	log.Print("Reading json...")
	err := tools.ReadJSON(w, r, &requestPayload)

	if err != nil {
		log.Print("Error on reading json ", err)
		tools.ErrorJSON(w, err)
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Print("Mongo insert error ", err)
		tools.ErrorJSON(w, err)
		return
	}

	resp := tools.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	tools.WriteJSON(w, http.StatusAccepted, resp)
}
