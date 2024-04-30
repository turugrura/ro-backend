package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewStoreRepository(collection *mongo.Collection) StoreRepository {
	return storeRepo{collection: collection}
}

type storeRepo struct {
	collection *mongo.Collection
}

func (repo storeRepo) FindStoreById(id string) (*Store, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var store Store
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&store)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (repo storeRepo) CreateStore(input CreateStoreInput) (*Store, error) {
	newStore := input.toModelForCreate()
	result, err := repo.collection.InsertOne(context.Background(), newStore)
	if err != nil {
		return nil, err
	}

	newStore.Id = result.InsertedID.(primitive.ObjectID).Hex()

	return &newStore, nil
}

func (repo storeRepo) UpdateStore(id string, input PatchStoreInput) (*Store, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	input.UpdatedAt = time.Now()

	_, err = repo.collection.UpdateByID(context.Background(), objId, bson.M{
		"$set": input,
	})
	if err != nil {
		return nil, err
	}

	return repo.FindStoreById(id)
}

func (repo storeRepo) UpdateRatingStore(storeId string, input UpdateRatingInput) (*Store, error) {
	store, err := repo.FindStoreById(storeId)
	store.Review[input.ReviewerId] = Review{
		Rating:  input.Rating,
		Comment: input.Comment,
	}
	if err != nil {
		return nil, err
	}

	totalRating := 0
	for _, review := range store.Review {
		totalRating += review.Rating
	}

	rating := float64(0)
	if len(store.Review) > 0 {
		rating = float64(totalRating) / float64(len(store.Review))
	}

	store.Rating = rating

	objId, err := primitive.ObjectIDFromHex(storeId)
	if err != nil {
		return nil, err
	}

	_, err = repo.collection.UpdateByID(context.Background(), objId, bson.M{
		"$set": bson.M{
			"review." + input.ReviewerId: Review{
				Rating:  input.Rating,
				Comment: input.Comment,
			},
			"rating": rating,
		},
	})
	if err != nil {
		return nil, err
	}

	return store, nil
}
