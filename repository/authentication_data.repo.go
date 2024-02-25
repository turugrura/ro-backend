package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAuthenticationDataRepo(collection *mongo.Collection) AuthenticationDataRepository {
	return authenticationDataRepo{collection: collection}
}

type authenticationDataRepo struct {
	collection *mongo.Collection
}

func (r authenticationDataRepo) DeleteAuthenticationDataByEmail(email string) error {
	_, err := r.collection.DeleteMany(context.Background(), DeleteAuthDataInput{
		Email: email,
	})

	return err
}

func (r authenticationDataRepo) DeleteAuthenticationDataById(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	return nil
}

func (r authenticationDataRepo) PartialSearchAuthData(i PartialSearchAuthDataInput) (*AuthenticationData, error) {
	var authenticationData AuthenticationData
	err := r.collection.FindOne(context.Background(), i).Decode(&authenticationData)
	if err != nil {
		return nil, err
	}

	return &authenticationData, nil
}

func (repo authenticationDataRepo) CreateAuthenticationData(i CreateAuthDataInput) (*string, error) {
	i.CreatedAt = time.Now()
	result, err := repo.collection.InsertOne(context.Background(), i)
	if err != nil {
		return nil, err
	}

	insertedID := fmt.Sprintf("%v", result.InsertedID)

	return &insertedID, nil
}
