package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepo(collection *mongo.Collection) UserRepository {
	return userRepo{collection: collection}
}

func (repo userRepo) CreateUser(input CreateUserInput) (*User, error) {
	var newUser = User{
		Name:      input.Name,
		Email:     input.Email,
		Role:      input.Role,
		Status:    UserStatus.Active,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	result, err := repo.collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return nil, err
	}

	newUser.Id = fmt.Sprintf("%v", result.InsertedID)

	return &newUser, nil
}

func (repo userRepo) FindUserById(id string) (*User, error) {
	var user = User{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = repo.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo userRepo) FindUserByEmail(email string) (*User, error) {
	var user = User{}
	err := repo.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
