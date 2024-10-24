package service

import "ro-backend/repository"

type MovieTranslatorService interface {
	GetAllEpisodes() ([]repository.MovieInfo, error)
	GetEpisode(ss, ep float32) (*repository.MovieTranslator, error)
	PatchSentence(ss, ep float32, sentence repository.PatchSentenceInput) (*repository.MovieTranslator, error)
}
