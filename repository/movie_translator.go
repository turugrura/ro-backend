package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Sentence struct {
	Raw      string  `bson:"raw" json:"raw"`
	Index    float32 `bson:"index" json:"index"`
	Speaker  string  `bson:"speaker" json:"speaker"`
	Sentence string  `bson:"sentence" json:"sentence"`
	Th       string  `bson:"th" json:"th"`
}

type MovieTranslator struct {
	Id        string     `bson:"_id,omitempty" json:"_id,omitempty"`
	Season    float32    `bson:"season" json:"season"`
	Episode   float32    `bson:"episode" json:"episode"`
	Title     string     `bson:"title" json:"title"`
	Sentences []Sentence `bson:"sentences" json:"sentences"`
	CreatedAt time.Time  `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updatedAt"`
	UpdatedBy string     `bson:"updated_by" json:"updatedBy"`
}

type MovieInfo struct {
	Id        string    `bson:"_id" json:"_id"`
	Season    float32   `bson:"season" json:"season,omitempty"`
	Episode   float32   `bson:"episode" json:"episode,omitempty"`
	Title     string    `bson:"title" json:"title"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt,omitempty"`
	UpdatedBy string    `bson:"updated_by" json:"updatedBy,omitempty"`
}

type MovieTranslatorIdentifier struct {
	Season  float32 `bson:"season"`
	Episode float32 `bson:"episode"`
}

type PatchSentenceInput struct {
	Index     float32 `bson:"index" json:"index"`
	Speaker   *string `bson:"speaker" json:"speaker"`
	Sentence  *string `bson:"sentence" json:"sentence"`
	Th        *string `bson:"th" json:"th"`
	UpdatedBy string  `bson:"updated_by" json:"updatedBy"`
}

func (p PatchSentenceInput) toIdentifier(ss, ep float32) bson.M {
	return bson.M{
		"season":          ss,
		"episode":         ep,
		"sentences.index": p.Index,
	}
}

func (p PatchSentenceInput) toUpdateModel() bson.M {
	var updateModel = bson.M{
		"updated_by": p.UpdatedBy,
		"updated_at": time.Now(),
	}

	if p.Speaker != nil {
		updateModel["sentences.$.speaker"] = *p.Speaker
	}
	if p.Sentence != nil {
		updateModel["sentences.$.sentence"] = *p.Sentence
	}
	if p.Th != nil {
		updateModel["sentences.$.th"] = *p.Th
	}

	return bson.M{
		"$set": updateModel,
	}
}

type PatchSentenceIdentifier struct {
	Season  float32 `bson:"season" `
	Episode float32 `bson:"episode" `
	Index   float32 `bson:"index" `
}

type MovieTranslatorRepository interface {
	GetAllEpisodes() ([]MovieInfo, error)
	GetEpisode(ss, ep float32) (*MovieTranslator, error)
	PatchSentence(ss, ep float32, sentence PatchSentenceInput) (*MovieTranslator, error)
}
