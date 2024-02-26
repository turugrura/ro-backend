package service

import (
	"ro-backend/repository"
	"time"
)

type CreateTagResult struct {
	Id        string
	Label     string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeleteTagInput struct {
	TagId    string
	UserId   string
	PresetId string
}

type RoPresetTag struct {
	repository.RoPreset
	TagId string
	Tags  map[string]int
	Liked bool
}

type PartialSearchTagsResult struct {
	Items []RoPresetTag
	Total int64
}

type PartialSearchMetaInput struct {
	UserId string
	Skip   int
	Limit  int
}

type PresetTagService interface {
	CreateTags(repository.CreateTagInput) (*CreateTagResult, error)
	DeleteTag(DeleteTagInput) (*CreateTagResult, error)
	LikeTag(repository.LikeTagInput) (*repository.PresetTag, error)
	UnLikeTag(repository.LikeTagInput) (*repository.PresetTag, error)
	PartialSearchTags(repository.PartialSearchTagsInput, PartialSearchMetaInput) (*PartialSearchTagsResult, error)
}
