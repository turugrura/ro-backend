package service

import "ro-backend/repository"

type CheckPresetOwnerInput struct {
	Id     string
	UserId string
}

type RoPresetService interface {
	FindPresetById(CheckPresetOwnerInput) (*repository.RoPreset, error)
	FindPresetsByUserId(string) (*[]repository.FindPreset, error)
	CreatePreset(repository.CreatePresetInput) (*repository.RoPreset, error)
	UpdatePreset(repository.UpdatePresetInput) (*repository.RoPreset, error)
	BulkCreatePresets(repository.BulkCreatePresetInput) (*[]repository.RoPreset, error)
	DeletePresetById(CheckPresetOwnerInput) (*int, error)
}
