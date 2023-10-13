package main

import (
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	log := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(log)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	res := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	_ = app.WriteJSON(w, http.StatusAccepted, res)
}
