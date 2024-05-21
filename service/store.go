package service

import "ro-backend/repository"

type StoreService interface {
	FindStoreById(storeId string) (*repository.Store, error)
	CreateStore(input repository.CreateStoreInput) (*repository.Store, error)
	UpdateStore(userId string, input repository.PatchStoreInput) (*repository.Store, error)
	UpdateRatingStore(storeId string, input repository.UpdateRatingInput) (*repository.Store, error)
}
