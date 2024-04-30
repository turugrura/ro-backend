package store_handler

import (
	"net/http"
	"ro-backend/service"
)

func NewStoreHandler(s service.StoreService) StoreHandler {
	return storeHandler{service: s}
}

type StoreHandler interface {
	FindStoreById(w http.ResponseWriter, r *http.Request)
	CreateStore(w http.ResponseWriter, r *http.Request)
	UpdateStore(w http.ResponseWriter, r *http.Request)
	ReviewStore(w http.ResponseWriter, r *http.Request)
}

type storeHandler struct {
	service service.StoreService
}
