package store_handler

import (
	"net/http"
	"ro-backend/core"
)

func (s storeHandler) FindMyStore(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	store, err := s.service.FindMyStore(userId)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
