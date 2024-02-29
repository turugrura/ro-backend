package service

import (
	"ro-backend/repository"
)

type DeleteTagInput struct {
	TagId    string
	UserId   string
	PresetId string
}

type PresetTag struct {
	repository.RoPreset
	TagId string
	Tags  map[string]int
	Liked bool
}

type TagWithLiked struct {
	repository.PresetTag
	Liked bool
}

type PresetWithTags struct {
	repository.RoPreset
	Tags []TagWithLiked
}

type PartialSearchTagsResult struct {
	Items []PresetTag
	Total int64
}

type PartialSearchMetaInput struct {
	UserId string
	Skip   int
	Limit  int
}

type BulkOperationInput struct {
	PublisherId string
	ClassId     int
	PresetId    string
	CreateTags  []string
	DeleteTags  []string
}

type PresetTagService interface {
	CreateTags(repository.CreateTagInput) (*PresetWithTags, error)
	BulkOperationTags(BulkOperationInput) (*PresetWithTags, error)
	DeleteTag(DeleteTagInput) (*PresetWithTags, error)
	LikeTag(repository.LikeTagInput) (*repository.PresetTag, error)
	UnLikeTag(repository.LikeTagInput) (*repository.PresetTag, error)
	PartialSearchTags(repository.PartialSearchTagsInput, PartialSearchMetaInput) (*PartialSearchTagsResult, error)
	AttachTags(userId string, p []repository.RoPreset) ([]PresetWithTags, error)
}
