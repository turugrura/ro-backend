package service

import "ro-backend/repository"

type CheckPresetOwnerRequest struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type FindPresetsByTagsRequest struct {
	ClassId int
	Skip    int
	Take    int
}

type RoPresetService interface {
	FindPresetById(CheckPresetOwnerRequest) (*repository.RoPreset, error)
	FindPresetsByUserId(string) ([]repository.RoPreset, error)
	CreatePreset(repository.CreatePresetInput) (*repository.RoPreset, error)
	BulkCreatePresets(repository.BulkCreatePresetInput) ([]repository.RoPreset, error)
	UpdatePreset(repository.UpdatePresetInput) (*repository.RoPreset, error)
	DeletePresetById(CheckPresetOwnerRequest) (*int, error)
}
