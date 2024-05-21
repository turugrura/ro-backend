package service

import (
	"ro-backend/repository"
)

func NewStoreService(repo repository.StoreRepository) StoreService {
	return storeService{storeRepo: repo}
}

type storeService struct {
	storeRepo repository.StoreRepository
}

func (s storeService) CreateStore(input repository.CreateStoreInput) (*repository.Store, error) {
	return s.storeRepo.CreateStore(input)
}

func (s storeService) FindStoreById(storeId string) (*repository.Store, error) {
	return s.storeRepo.FindStoreById(storeId)
}

func (s storeService) UpdateRatingStore(storeId string, input repository.UpdateRatingInput) (*repository.Store, error) {
	return s.storeRepo.UpdateRatingStore(storeId, input)
}

func (s storeService) UpdateStore(userId string, input repository.PatchStoreInput) (*repository.Store, error) {
	store, err := s.storeRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	return s.storeRepo.UpdateStore(store.Id.Hex(), input)
}
