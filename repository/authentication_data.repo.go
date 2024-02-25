package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type authenticationDataRepo struct {
	collection *mongo.Collection
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

func NewAuthenticationDataRepo(collection *mongo.Collection) AuthenticationDataRepository {
	return authenticationDataRepo{collection: collection}
}

func (r authenticationDataRepo) GetAuthenticationById(id string) (*AuthenticationData, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var authenticationData AuthenticationData
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&authenticationData)
	if err != nil {
		return nil, err
	}

	return &authenticationData, nil
}

func (r authenticationDataRepo) GetAuthenticationByEmail(email string) (*AuthenticationData, error) {
	var authenticationData AuthenticationData
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&authenticationData)
	if err != nil {
		return nil, err
	}

	return &authenticationData, nil
}

func (repo authenticationDataRepo) CreateAuthenticationData(params CreateAuthenticationDataInput) (*AuthenticationData, error) {
	var data = AuthenticationData{
		AuthReference: params.AuthReference,
		Email:         params.Email,
		CreatedAt:     time.Now(),
	}

	result, err := repo.collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	data.Id = fmt.Sprintf("%v", result.InsertedID)

	return &data, nil
}
