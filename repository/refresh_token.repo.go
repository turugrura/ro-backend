package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type refreshTokenRepo struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepo(collection *mongo.Collection) RefreshTokenRepository {
	return refreshTokenRepo{collection: collection}
}

func (repo refreshTokenRepo) GetRefreshTokenById(id string) (*RefreshToken, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var refreshToken RefreshToken
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&refreshToken)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (repo refreshTokenRepo) DeleteRefreshTokenByUserId(userId string) error {
	_, err := repo.collection.DeleteMany(context.Background(), bson.M{"user_id": userId})

	return err
}

func (repo refreshTokenRepo) CreateRefreshToken(input CreateRefreshTokenInput) (*RefreshToken, error) {
	now := time.Now()
	var newRefreshToken = RefreshToken{
		UserId:    input.UserId,
		Count:     input.Count,
		UserAgent: input.UserAgent,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	result, err := repo.collection.InsertOne(context.Background(), newRefreshToken)
	if err != nil {
		return nil, err
	}

	var newOne RefreshToken
	repo.collection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&newOne)

	newRefreshToken.Id = result.InsertedID.(primitive.ObjectID).Hex()

	return &newRefreshToken, nil
}

func (repo refreshTokenRepo) UpdateRefreshToken(input UpdateRefreshTokenInput) (*RefreshToken, error) {
	objectId, err := primitive.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, err
	}

	var updatedRefreshToken RefreshToken
	err = repo.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectId}, bson.M{
		"$set": PatchRefreshToken{
			Count:     input.Count,
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}).Decode(&updatedRefreshToken)
	if err != nil {
		return nil, err
	}

	// res, err := repo.collection.UpdateByID(context.Background(), objectId, RefreshToken{
	// 	Count:     input.Count,
	// 	UpdatedAt: time.Now().Format(time.RFC3339),
	// })
	// fmt.Println("ModifiedCount", res.ModifiedCount)

	// if result.ModifiedCount != 1 {
	// 	return nil, fmt.Errorf("%v not found", input.Id)
	// }

	return &updatedRefreshToken, nil
}
