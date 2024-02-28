package handler

import "net/http"

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
}

func (h helpCheckHandler) Ping(w http.ResponseWriter, r *http.Request) {
	response := PingPongResponse{
		Message: "Pong",
	}

	WriteOK(w, response)
}
