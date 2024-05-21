package store_handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
)

type PatchStoreRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Fb            string `json:"fb"`
	CharacterName string `json:"characterName"`
}

func (r PatchStoreRequest) toPatchInput() repository.PatchStoreInput {
	return repository.PatchStoreInput{
		Name:          r.Name,
		Description:   r.Description,
		Fb:            r.Fb,
		CharacterName: r.CharacterName,
	}
}

func (s storeHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	d := PatchStoreRequest{}
	json.NewDecoder(r.Body).Decode(&d)

	store, err := s.service.UpdateStore(userId, d.toPatchInput())
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
