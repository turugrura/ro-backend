package repository

import (
	"slices"
	"time"
)

type PresetTag struct {
	Id          string    `bson:"_id,omitempty"`
	PublisherId string    `bson:"publisher_id"`
	Tag         string    `bson:"tag"`
	ClassId     int       `bson:"class_id"`
	PresetId    string    `bson:"preset_id"`
	Likes       []string  `bson:"likes"`
	TotalLike   int       `bson:"total_like"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func (p *PresetTag) IsILike(userId string) bool {
	return slices.Contains(p.Likes, userId)
}

type CreateTagInput struct {
	PublisherId string
	Tags        []string
	ClassId     int
	PresetId    string
}

type PartialUpdateTagInput struct {
	TotalLike int       `bson:"total_like"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type LikeTagInput struct {
	Id        string
	UserId    string
	TotalLike int
}

type PartialSearchTagsInput struct {
	PublisherId string `bson:"publisher_id,omitempty"`
	Tag         string `bson:"tag,omitempty"`
	ClassId     int    `bson:"class_id,omitempty"`
	PresetId    string `bson:"preset_id,omitempty"`
}

type PartialSearchSorting struct {
	TotalLike int       `bson:"total_like"`
	CreatedAt time.Time `bson:"created_at"`
}

type PartialSearchTagsResult struct {
	Items []PresetTag
	Total int
}

type PresetTagRepository interface {
	FindTagById(string) (*PresetTag, error)
	FindTagsByPresetId(string) ([]PresetTag, error)
	CreateTags(CreateTagInput) ([]string, error)
	BulkOperationTags(createInput CreateTagInput, createIds []string) error
	DeleteTag(id string) error
	DeleteTagsByPresetId(presetId string) error
	LikeTag(LikeTagInput) error
	UnLikeTag(LikeTagInput) error
	PartialSearchTags(i PartialSearchTagsInput, skip, limit int) (*PartialSearchTagsResult, error)
	FindByPresetIds([]string) ([]PresetTag, error)
}
