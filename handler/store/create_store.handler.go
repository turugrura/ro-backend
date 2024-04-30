package store_handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
)

type CreateStoreRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Fb            string `json:"fb"`
	CharacterName string `json:"characterName"`
}

func (r CreateStoreRequest) toCreateInput(userId string) repository.CreateStoreInput {
	return repository.CreateStoreInput{
		OwnerId:       userId,
		Name:          r.Name,
		Description:   r.Description,
		Fb:            r.Fb,
		CharacterName: r.CharacterName,
	}
}

func (s storeHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	var d CreateStoreRequest
	json.NewDecoder(r.Body).Decode(&d)

	store, err := s.service.CreateStore(d.toCreateInput(userId))
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
