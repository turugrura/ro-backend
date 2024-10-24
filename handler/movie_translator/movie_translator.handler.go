package movie_translator_handler

import (
	"net/http"
	"ro-backend/service"
)

func NewMovieTranslatorHandler(s service.MovieTranslatorService) MovieTranslatorHandler {
	return movieTranslatorHandler{service: s}
}

type MovieTranslatorHandler interface {
	GetAllEpisodes(w http.ResponseWriter, r *http.Request)
	GetEpisode(w http.ResponseWriter, r *http.Request)
	PatchSentence(w http.ResponseWriter, r *http.Request)
}

type movieTranslatorHandler struct {
	service service.MovieTranslatorService
}
