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

func NewAuthenticationDataRepo(collection *mongo.Collection) AuthenticationDataRepository {
	return authenticationDataRepo{collection: collection}
}

func (repo authenticationDataRepo) GetAuthenticationById(id string) (*AuthenticationData, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var authenticationData AuthenticationData
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&authenticationData)
	if err != nil {
		return nil, err
	}

	return &authenticationData, nil
}

func (repo authenticationDataRepo) GetAuthenticationByEmail(email string) (*AuthenticationData, error) {
	var authenticationData AuthenticationData
	err := repo.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&authenticationData)
	if err != nil {
		return nil, err
	}

	return &authenticationData, nil
}

func (repo authenticationDataRepo) CreateAuthenticationData(params CreateAuthenticationDataInput) (*AuthenticationData, error) {
	var data = AuthenticationData{
		AuthReference: params.AuthReference,
		Email:         params.Email,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	result, err := repo.collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, err
	}

	data.Id = fmt.Sprintf("%v", result.InsertedID)

	return &data, nil
}
