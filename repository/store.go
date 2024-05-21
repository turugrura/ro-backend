package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	Rating  int    `bson:"rating" json:"rating"`
	Comment string `bson:"comment" json:"comment"`
}

type Store struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	OwnerId       primitive.ObjectID `bson:"owner_id" json:"ownerId"`
	Review        map[string]Review  `bson:"review" json:"review"`
	Rating        float64            `bson:"rating" json:"rating"`
	Fb            string             `bson:"fb" json:"fb"`
	CharacterName string             `bson:"character_name" json:"characterName"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

type PatchStoreInput struct {
	Name          string    `bson:"name,omitempty"`
	Description   string    `bson:"description,omitempty"`
	Fb            string    `bson:"fb,omitempty"`
	CharacterName string    `bson:"character_name,omitempty"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

type CreateStoreInput struct {
	Name          string
	Description   string
	OwnerId       string
	Fb            string
	CharacterName string
}

func (c CreateStoreInput) toModelForCreate(ownerId primitive.ObjectID) Store {
	now := time.Now()

	return Store{
		Name:          c.Name,
		Description:   c.Description,
		OwnerId:       ownerId,
		Review:        map[string]Review{},
		Fb:            c.Fb,
		CharacterName: c.CharacterName,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

type UpdateRatingInput struct {
	ReviewerId string
	Rating     int
	Comment    string
}

type StoreRepository interface {
	FindStoreById(storeId string) (*Store, error)
	FindStoreByOwnerId(ownerId string) (*Store, error)
	CreateStore(input CreateStoreInput) (*Store, error)
	UpdateStore(storeId string, input PatchStoreInput) (*Store, error)
	UpdateRatingStore(storeId string, input UpdateRatingInput) (*Store, error)
}
