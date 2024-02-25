package main

import (
	"net/http"

	"github.com/deepMOORE/tools"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSmtpMessage(msg)
	if err != nil {
		tools.ErrorJSON(w, err)
		return
	}

	payload := tools.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	tools.WriteJSON(w, http.StatusAccepted, payload)
}
