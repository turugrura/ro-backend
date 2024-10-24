package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection
var authDataCollection *mongo.Collection
var refreshTokenCollection *mongo.Collection
var roPresetCollection *mongo.Collection
var roPresetForSummaryCollection *mongo.Collection
var roTagCollection *mongo.Collection

// var storeCollection *mongo.Collection
// var productCollection *mongo.Collection
var friendTranslatorCollection *mongo.Collection

func connectMongoDB() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.Mongodb.ConnectionStr))
	if err != nil {
		panic(fmt.Errorf("fatal error connect DB: %w", err))
	}

	mongoDb := client.Database(appConfig.Mongodb.DbName)
	userCollection = mongoDb.Collection("users")
	_, err = userCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"name": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		panic(fmt.Errorf("index users: %w", err))
	}

	authDataCollection = mongoDb.Collection("authorization_codes")
	_, err = authDataCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"code": 1,
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index authorization_codes: %w", err))
	}

	refreshTokenCollection = mongoDb.Collection("refresh_tokens")

	roPresetCollection = mongoDb.Collection("ro_presets")
	roPresetForSummaryCollection = mongoDb.Collection("authorization_codes")
	_, err = roPresetCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{
				"id": 1,
			},
		},
		{
			Keys: bson.M{
				"user_id": 1,
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index ro_presets: %w", err))
	}

	roTagCollection = mongoDb.Collection("preset_tags")
	_, err = roTagCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tag", Value: 1},
				{Key: "class_id", Value: 1},
			},
		},
		{
			Keys: bson.M{
				"preset_id": 1,
			},
		},
		{
			Keys: bson.D{
				{Key: "preset_id", Value: 1},
				{Key: "tag", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "total_like", Value: 1},
				{Key: "created_at", Value: 1},
			},
		},
	})
	if err != nil {
		panic(fmt.Errorf("index preset_tags: %w", err))
	}

	// storeCollection = mongoDb.Collection("store")
	// _, err = storeCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
	// 	{
	// 		Keys: bson.M{
	// 			"owner_id": 1,
	// 		},
	// 		Options: options.Index().SetUnique(true),
	// 	},
	// 	{
	// 		Keys: bson.M{
	// 			"name": 1,
	// 		},
	// 		Options: options.Index().SetUnique(true),
	// 	},
	// })
	// if err != nil {
	// 	panic(fmt.Errorf("index store: %w", err))
	// }

	// productCollection = mongoDb.Collection("product")
	// _, err = productCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "exp_date", Value: 1},
	// 		},
	// 	},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "baht", Value: 1},
	// 		},
	// 	},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "m", Value: 1},
	// 		},
	// 	},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "m", Value: 1},
	// 			{Key: "baht", Value: 1},
	// 		},
	// 	},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "m", Value: 1},
	// 			{Key: "baht", Value: 1},
	// 			{Key: "exp_date", Value: 1},
	// 		},
	// 	},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "m", Value: 1},
	// 			{Key: "baht", Value: 1},
	// 			{Key: "exp_date", Value: 1},
	// 			{Key: "is_published", Value: -1},
	// 		},
	// 	},
	// 	{
	// 		Keys: bson.D{
	// 			{Key: "m", Value: 1},
	// 			{Key: "baht", Value: 1},
	// 			{Key: "exp_date", Value: 1},
	// 			{Key: "is_published", Value: -1},
	// 			{Key: "name", Value: 1},
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	panic(fmt.Errorf("index product: %w", err))
	// }

	friendTranslatorCollection = mongoDb.Collection("friends")
	_, err = friendTranslatorCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "episode", Value: 1},
				{Key: "season", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		panic(fmt.Errorf("index friends: %w", err))
	}

	return
}
