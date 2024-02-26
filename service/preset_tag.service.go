package service

import (
	"fmt"
	"ro-backend/appError"
	"ro-backend/repository"
	"slices"
)

func NewPresetTagService(tRepo repository.PresetTagRepository, pRepo repository.RoPresetRepository) PresetTagService {
	return presetTagService{tRepo: tRepo, pRepo: pRepo}
}

type presetTagService struct {
	pRepo repository.RoPresetRepository
	tRepo repository.PresetTagRepository
}

func (s presetTagService) ValidatePresetOwner(r CheckPresetOwnerRequest) (*repository.RoPreset, error) {
	res, err := s.pRepo.FindPresetById(repository.FindPresetByIdInput{
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

func (s presetTagService) AttachTags(userId string, presets []repository.RoPreset) ([]PresetWithTags, error) {
	presetIds := []string{}
	for _, v := range presets {
		presetIds = append(presetIds, v.Id)
	}

	tgs, err := s.tRepo.FindByPresetIds(presetIds)
	if err != nil {
		return nil, err
	}

	presetTagsMap := map[string][]TagWithLiked{}
	for _, v := range tgs {
		if presetTagsMap[v.PresetId] == nil {
			presetTagsMap[v.PresetId] = []TagWithLiked{}
		}
		presetTagsMap[v.PresetId] = append(presetTagsMap[v.PresetId], TagWithLiked{
			PresetTag: v,
			Liked:     v.IsILike(userId),
		})
	}

	presetTags := []PresetWithTags{}
	for _, v := range presets {
		presetTags = append(presetTags, PresetWithTags{
			RoPreset: v,
			Tags:     presetTagsMap[v.Id],
		})
	}

	return presetTags, nil
}

func (s presetTagService) CreateTags(i repository.CreateTagInput) (*PresetWithTags, error) {
	p, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{
		Id:     i.PresetId,
		UserId: i.PublisherId,
	})
	if err != nil {
		return nil, err
	}
	if !p.IsPublished {
		return nil, fmt.Errorf(appError.ErrCannotTagUnpublished)
	}

	i.ClassId = p.ClassId
	_, err = s.tRepo.CreateTags(i)
	if err != nil {
		return nil, err
	}

	t, err := s.AttachTags(i.PublisherId, []repository.RoPreset{*p})
	if err != nil {
		return nil, err
	}

	return &t[0], nil
}

func (s presetTagService) DeleteTag(i DeleteTagInput) (*PresetWithTags, error) {
	p, err := s.ValidatePresetOwner(CheckPresetOwnerRequest{
		Id:     i.PresetId,
		UserId: i.UserId,
	})
	if err != nil {
		return nil, err
	}
	if !p.IsPublished {
		return nil, fmt.Errorf(appError.ErrCannotTagUnpublished)
	}

	err = s.tRepo.DeleteTag(i.TagId)
	if err != nil {
		return nil, err
	}

	t, err := s.AttachTags(i.UserId, []repository.RoPreset{*p})
	if err != nil {
		return nil, err
	}

	return &t[0], nil
}

func (s presetTagService) LikeTag(i repository.LikeTagInput) (*repository.PresetTag, error) {
	tag, err := s.tRepo.FindTagById(i.Id)
	if err != nil {
		return nil, err
	}
	if slices.Contains(tag.Likes, i.UserId) {
		return tag, nil
	}

	tag.TotalLike = tag.TotalLike + 1
	i.TotalLike = tag.TotalLike

	return tag, s.tRepo.LikeTag(i)
}

func (s presetTagService) PartialSearchTags(i repository.PartialSearchTagsInput, si PartialSearchMetaInput) (*PartialSearchTagsResult, error) {
	tags, err := s.tRepo.PartialSearchTags(i, si.Skip, si.Limit)
	if err != nil {
		return nil, err
	}

	presetIds := []string{}
	for _, v := range tags.Items {
		presetIds = append(presetIds, v.PresetId)
	}

	presets, err := s.pRepo.FindPresetByIds(presetIds)
	if err != nil {
		return nil, err
	}
	presetMap := map[string]repository.RoPreset{}
	for _, v := range presets {
		presetMap[v.Id] = v
	}

	tgs, err := s.tRepo.FindByPresetIds(presetIds)
	if err != nil {
		return nil, err
	}

	presetTagsMap := map[string]map[string]int{}
	for _, v := range tgs {
		if presetTagsMap[v.PresetId] == nil {
			presetTagsMap[v.PresetId] = map[string]int{}
		}
		presetTagsMap[v.PresetId][v.Tag] = len(v.Likes)
	}

	presetTags := []PresetTag{}
	for _, v := range tags.Items {
		presetTags = append(presetTags, PresetTag{
			RoPreset: presetMap[v.PresetId],
			TagId:    v.Id,
			Tags:     presetTagsMap[v.PresetId],
			Liked:    v.IsILike(si.UserId),
		})
	}

	return &PartialSearchTagsResult{
		Total: int64(tags.Total),
		Items: presetTags,
	}, nil
}

func (s presetTagService) UnLikeTag(i repository.LikeTagInput) (*repository.PresetTag, error) {
	tag, err := s.tRepo.FindTagById(i.Id)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(tag.Likes, i.UserId) {
		return tag, nil
	}

	tag.TotalLike = tag.TotalLike - 1
	i.TotalLike = tag.TotalLike

	return tag, s.tRepo.UnLikeTag(i)
}
