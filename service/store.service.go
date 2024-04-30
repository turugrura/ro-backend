package service

import (
	"fmt"
	"ro-backend/appError"
	"ro-backend/repository"
)

func NewStoreService(repo repository.StoreRepository) StoreService {
	return storeService{storeRepo: repo}
}

type storeService struct {
	storeRepo repository.StoreRepository
}

func (s storeService) validateStoreOwner(userId, storeId string) (*repository.Store, error) {
	store, err := s.storeRepo.FindStoreById(storeId)
	if err != nil {
		return nil, err
	}

	if store.OwnerId != userId {
		return nil, fmt.Errorf(appError.ErrNotMyPreset)
	}

	return store, nil
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

func (s storeService) UpdateStore(userId, storeId string, input repository.PatchStoreInput) (*repository.Store, error) {
	_, err := s.validateStoreOwner(userId, storeId)
	if err != nil {
		return nil, err
	}

	return s.storeRepo.UpdateStore(storeId, input)
}
