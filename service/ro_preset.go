package service

import "ro-backend/repository"

type CheckPresetOwnerRequest struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type AddTagsRequest struct {
	CheckPresetOwnerRequest
	Tags []string `json:"tags"`
}

type FindPresetsByTagsRequest struct {
	ClassId int
	Tag     string
	Skip    int
	Take    int
}

type RoPresetService interface {
	FindPresetById(CheckPresetOwnerRequest) (*repository.RoPreset, error)
	FindPresetsByUserId(string) (*[]repository.RoPreset, error)
	FindPresetsByTags(FindPresetsByTagsRequest) (*repository.PartialSearchRoPresetResult, error)
	CreatePreset(repository.CreatePresetInput) (*repository.RoPreset, error)
	BulkCreatePresets(repository.BulkCreatePresetInput) (*[]repository.RoPreset, error)
	AddTags(AddTagsRequest) (*repository.RoPreset, error)
	RemoveTags(AddTagsRequest) (*repository.RoPreset, error)
	UpdatePreset(repository.UpdatePresetInput) (*repository.RoPreset, error)
	DeletePresetById(CheckPresetOwnerRequest) (*int, error)
}
