package handler

import (
	"net/http"
	"ro-backend/configuration"
)

func NewHelpCheckHandler() HelpCheckHandler {
	return helpCheckHandler{}
}

type HelpCheckHandler interface {
	Ping(http.ResponseWriter, *http.Request)
}

type helpCheckHandler struct {
}

type PingPongResponse struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

func (h helpCheckHandler) Ping(w http.ResponseWriter, r *http.Request) {
	response := PingPongResponse{
		Message: "Pong",
		From:    configuration.Config.Environment,
	}

	WriteOK(w, response)
}
