package store_handler

import (
	"net/http"
	"ro-backend/core"

	"github.com/gorilla/mux"
)

func (s storeHandler) FindStoreById(w http.ResponseWriter, r *http.Request) {
	storeId := mux.Vars(r)["storeId"]

	store, err := s.service.FindStoreById(storeId)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
