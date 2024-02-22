package service

import (
	"fmt"
	"ro-backend/repository"
)

func NewRoPresetService(repo repository.RoPresetRepository) RoPresetService {
	return roPresetService{repo: repo}
}

type roPresetService struct {
	repo repository.RoPresetRepository
}

func (s roPresetService) ValidatePresetOwner(i CheckPresetOwnerInput) error {
	res, err := s.repo.FindPresetById(i.Id)
	if err != nil {
		return err
	}

	if res.UserId != i.UserId {
		return fmt.Errorf("not my preset")
	}

	return nil
}

// UpdatePreset implements RoPresetService.
func (s roPresetService) UpdatePreset(i repository.UpdatePresetInput) (*repository.RoPreset, error) {
	err := s.ValidatePresetOwner(CheckPresetOwnerInput{Id: i.Id, UserId: i.UserId})
	if err != nil {
		return nil, err
	}

	return s.repo.UpdatePreset(i)
}

// DeletePresetById implements RoPresetService.
func (s roPresetService) DeletePresetById(i CheckPresetOwnerInput) (*int, error) {
	err := s.ValidatePresetOwner(CheckPresetOwnerInput{Id: i.Id, UserId: i.UserId})
	if err != nil {
		return nil, err
	}

	return s.repo.DeletePresetById(i.Id)
}

// BulkCreatePresets implements RoPresetService.
func (s roPresetService) BulkCreatePresets(i repository.BulkCreatePresetInput) (*[]repository.RoPreset, error) {
	return s.repo.CreatePresets(i)
}

// FindPresetsByUserId implements RoPresetService.
func (s roPresetService) FindPresetsByUserId(userId string) (*[]repository.FindPreset, error) {
	return s.repo.FindPresetsByUserId(userId)
}

// CreatePreset implements RoPresetService.
func (s roPresetService) CreatePreset(r repository.CreatePresetInput) (*repository.RoPreset, error) {
	res, err := s.repo.CreatePreset(r)

	return (*repository.RoPreset)(res), err
}

// FindPresetById implements RoPresetService.
func (s roPresetService) FindPresetById(i CheckPresetOwnerInput) (*repository.RoPreset, error) {
	err := s.ValidatePresetOwner(CheckPresetOwnerInput{Id: i.Id, UserId: i.UserId})
	if err != nil {
		return nil, err
	}

	res, err := s.repo.FindPresetById(i.Id)

	return (*repository.RoPreset)(res), err
}
