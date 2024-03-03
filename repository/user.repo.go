package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewUserRepo(collection *mongo.Collection) UserRepository {
	return userRepo{collection: collection}
}

type userRepo struct {
	collection *mongo.Collection
}

func (repo userRepo) FindUsersByIds(ids []string) ([]User, error) {
	objIds := []primitive.ObjectID{}
	for _, v := range ids {
		objId, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, err
		}
		objIds = append(objIds, objId)
	}

	cur, err := repo.collection.Find(context.Background(), bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	})
	if err != nil {
		return nil, err
	}

	users := []User{}
	err = cur.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo userRepo) CreateUser(input CreateUserInput) (*User, error) {
	var newUser = User{
		Name:            input.Name,
		Email:           input.Email,
		Role:            input.Role,
		RegisterChannel: input.RegisterChannel,
		Status:          UserStatus.Active,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	result, err := repo.collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return nil, err
	}

	newUser.Id = fmt.Sprintf("%v", result.InsertedID)

	return &newUser, nil
}

func (r userRepo) PatchUser(id string, u UpdateUserInput) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	u.UpdatedAt = time.Now()

	_, err = r.collection.UpdateByID(context.Background(), objId, bson.M{
		"$set": u,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r userRepo) FindUserById(id string) (*User, error) {
	var user = User{}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r userRepo) FindUserByEmail(email string) (*User, error) {
	var user = User{}
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
