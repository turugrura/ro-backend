package service

import (
	"fmt"
	"ro-backend/appError"
	"ro-backend/repository"
	"time"
)

func NewRoPresetService(repo repository.RoPresetRepository, tagRepo repository.PresetTagRepository) RoPresetService {
	return roPresetService{presetRepo: repo, tagRepo: tagRepo}
}

type roPresetService struct {
	presetRepo repository.RoPresetRepository
	tagRepo    repository.PresetTagRepository
}

func (s roPresetService) ValidatePresetOwner(r CheckPresetOwnerRequest) (*repository.RoPreset, error) {
	res, err := s.presetRepo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: false,
	})
	if err != nil {
		return nil, err
	}

	if res.UserId != r.UserId {
		return nil, fmt.Errorf(appError.ErrNotMyPreset)
	}

	return res, nil
}

func (s roPresetService) UpdatePreset(id string, i repository.UpdatePresetInput) (*repository.RoPreset, error) {
	p, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: id, UserId: i.UserId})
	if err != nil {
		return nil, err
	}

	if p.IsPublished {
		return nil, fmt.Errorf(appError.ErrCannotUpdatePublishedPreset)
	}

	err = s.presetRepo.UpdatePreset(id, repository.UpdatePresetInput{
		Label: i.Label,
		Model: i.Model,
	})
	if err != nil {
		return nil, err
	}

	return s.presetRepo.FindPresetById(repository.FindPresetByIdInput{
		Id:           id,
		InCludeModel: false,
	})
}

func (s roPresetService) PublishPreset(id string, i repository.UpdatePresetInput) (*repository.RoPreset, error) {
	p, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: id, UserId: i.UserId})
	if err != nil {
		return nil, err
	}

	if p.IsPublished {
		return nil, fmt.Errorf(appError.ErrCannotUpdatePublishedPreset)
	}

	err = s.presetRepo.UpdatePreset(id, repository.UpdatePresetInput{
		PublishName: i.PublishName,
		IsPublished: true,
		PublishedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return s.presetRepo.FindPresetById(repository.FindPresetByIdInput{
		Id:           id,
		InCludeModel: false,
	})
}

func (s roPresetService) UnPublishPreset(id string, i repository.UpdatePresetInput) (*repository.RoPreset, error) {
	p, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: id, UserId: i.UserId})
	if err != nil {
		return nil, err
	}

	if !p.IsPublished {
		return p, nil
	}

	err = s.tagRepo.DeleteTagsByPresetId(p.Id)
	if err != nil {
		return nil, err
	}

	err = s.presetRepo.UnpublishedPreset(id)
	if err != nil {
		return nil, err
	}

	return s.presetRepo.FindPresetById(repository.FindPresetByIdInput{
		Id:           id,
		InCludeModel: false,
	})
}

func (s roPresetService) DeletePresetById(r CheckPresetOwnerRequest) (*int, error) {
	_, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	return s.presetRepo.DeletePresetById(r.Id)
}

func (s roPresetService) BulkCreatePresets(r repository.BulkCreatePresetInput) ([]repository.RoPreset, error) {
	return s.presetRepo.CreatePresets(r)
}

func (s roPresetService) FindPresetsByUserId(userId string, includeModel bool) ([]repository.RoPreset, error) {
	res, err := s.presetRepo.PartialSearchPresets(repository.PartialSearchRoPresetInput{
		UserId:       &userId,
		InCludeModel: includeModel,
	})
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (s roPresetService) CreatePreset(r repository.CreatePresetInput) (*repository.RoPreset, error) {
	res, err := s.presetRepo.CreatePreset(r)

	return (*repository.RoPreset)(res), err
}

func (s roPresetService) FindPresetById(r CheckPresetOwnerRequest) (*repository.RoPreset, error) {
	_, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	res, err := s.presetRepo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: true,
	})

	return (*repository.RoPreset)(res), err
}
