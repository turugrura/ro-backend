package service

import (
	"fmt"
	"ro-backend/repository"

	"golang.org/x/exp/slices"
)

func NewRoPresetService(repo repository.RoPresetRepository) RoPresetService {
	return roPresetService{repo: repo}
}

type roPresetService struct {
	repo repository.RoPresetRepository
}

func (s roPresetService) FindPresetsByTags(r FindPresetsByTagsRequest) (*repository.PartialSearchRoPresetResult, error) {
	res, err := s.repo.PartialSearchPresets(repository.PartialSearchRoPresetInput{
		ClassId:      &r.ClassId,
		Tag:          &r.Tag,
		Skip:         &r.Skip,
		Take:         &r.Take,
		InCludeModel: true,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s roPresetService) AddTags(r AddTagsRequest) (*repository.RoPreset, error) {
	preset, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	// tags := []string{}
	// for _, v := range preset.Tags {
	// 	_, wantAdd := slices.BinarySearch(r.Tags, v)
	// 	if wantAdd {
	// 		continue
	// 	}
	// 	tags = append(tags, v)
	// }

	err = s.repo.UpdatePresetTags(repository.UpdateTagsInput{
		Id:   preset.Id,
		Tags: r.Tags,
	})
	if err != nil {
		return nil, err
	}

	return s.repo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: false,
	})
}

func (s roPresetService) RemoveTags(r AddTagsRequest) (*repository.RoPreset, error) {
	preset, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	tags := []string{}
	for _, v := range preset.Tags {
		wantDelete := slices.Contains(r.Tags, v)
		if wantDelete {
			continue
		}
		tags = append(tags, v)
	}

	err = s.repo.UpdatePreset(repository.UpdatePresetInput{
		Id:   preset.Id,
		Tags: tags,
	})
	if err != nil {
		return nil, err
	}

	return s.repo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: false,
	})
}

func (s roPresetService) ValidatePresetOwner(r CheckPresetOwnerRequest) (*repository.RoPreset, error) {
	res, err := s.repo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: false,
	})
	if err != nil {
		return nil, err
	}

	if res.UserId != r.UserId {
		return nil, fmt.Errorf("not my preset")
	}

	return res, nil
}

func (s roPresetService) UpdatePreset(r repository.UpdatePresetInput) (*repository.RoPreset, error) {
	_, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdatePreset(r)
	if err != nil {
		return nil, err
	}

	return s.repo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: false,
	})
}

func (s roPresetService) DeletePresetById(r CheckPresetOwnerRequest) (*int, error) {
	_, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	return s.repo.DeletePresetById(r.Id)
}

func (s roPresetService) BulkCreatePresets(r repository.BulkCreatePresetInput) (*[]repository.RoPreset, error) {
	return s.repo.CreatePresets(r)
}

func (s roPresetService) FindPresetsByUserId(userId string) (*[]repository.RoPreset, error) {
	res, err := s.repo.PartialSearchPresets(repository.PartialSearchRoPresetInput{
		UserId:       &userId,
		InCludeModel: false,
	})
	if err != nil {
		return nil, err
	}

	return &res.Items, nil
}

func (s roPresetService) CreatePreset(r repository.CreatePresetInput) (*repository.RoPreset, error) {
	res, err := s.repo.CreatePreset(r)

	return (*repository.RoPreset)(res), err
}

func (s roPresetService) FindPresetById(r CheckPresetOwnerRequest) (*repository.RoPreset, error) {
	_, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{Id: r.Id, UserId: r.UserId})
	if err != nil {
		return nil, err
	}

	res, err := s.repo.FindPresetById(repository.FindPresetByIdInput{
		Id:           r.Id,
		InCludeModel: true,
	})

	return (*repository.RoPreset)(res), err
}
