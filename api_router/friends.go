package api_router

import (
	_movieTranslatorHandler "ro-backend/handler/movie_translator"
	"ro-backend/repository"
	"ro-backend/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouterFriend(c *mongo.Collection, router AppRouter) {
	var friendTranslatorRepo = repository.NewMovieTranslatorRepository(c)
	var friendTranslatorService = service.NewMovieTranslatorService(friendTranslatorRepo)
	var friendTranslatorHandler = _movieTranslatorHandler.NewMovieTranslatorHandler(friendTranslatorService)

	router.Get("/friends", friendTranslatorHandler.GetAllEpisodes)
	router.Get("/friends/{ss}/{ep}", friendTranslatorHandler.GetEpisode)
	router.Post("/friends/{ss}/{ep}", friendTranslatorHandler.PatchSentence)
}
