package service

import "ro-backend/repository"

func NewMovieTranslatorService(repo repository.MovieTranslatorRepository) MovieTranslatorService {
	return movieTranslatorService{r: repo}
}

type movieTranslatorService struct {
	r repository.MovieTranslatorRepository
}

func (m movieTranslatorService) GetAllEpisodes() ([]repository.MovieInfo, error) {
	return m.r.GetAllEpisodes()
}

func (m movieTranslatorService) GetEpisode(ss float32, ep float32) (*repository.MovieTranslator, error) {
	return m.r.GetEpisode(ss, ep)
}

func (m movieTranslatorService) PatchSentence(ss float32, ep float32, sentence repository.PatchSentenceInput) (*repository.MovieTranslator, error) {
	return m.r.PatchSentence(ss, ep, sentence)
}
