package store_handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"

	"github.com/gorilla/mux"
)

type UpdateRatingRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

func (r UpdateRatingRequest) toReviewInput(userId string) repository.UpdateRatingInput {
	return repository.UpdateRatingInput{
		ReviewerId: userId,
		Rating:     r.Rating,
		Comment:    r.Comment,
	}
}

func (s storeHandler) ReviewStore(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	storeId := mux.Vars(r)["storeId"]

	var d UpdateRatingRequest
	json.NewDecoder(r.Body).Decode(&d)

	if d.Rating < 0 || d.Rating > 5 {
		core.WriteErr(w, "rating must be a value between 0 - 5")
		return
	}

	store, err := s.service.UpdateRatingStore(storeId, d.toReviewInput(userId))
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
