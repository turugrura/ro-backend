package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRoPresetRepository(collection *mongo.Collection) RoPresetRepository {
	return roPresetRepo{collection: collection}
}

type roPresetRepo struct {
	collection *mongo.Collection
}

// DeletePresetById implements RoPresetRepository.
func (r roPresetRepo) DeletePresetById(id string) (*int, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := r.collection.DeleteOne(context.Background(), objectId)
	if err != nil {
		return nil, err
	}

	i := int(res.DeletedCount)

	return &i, nil
}

// CreatePreset implements RoPresetRepository.
func (r roPresetRepo) CreatePreset(i CreatePresetInput) (*RoPreset, error) {
	res, err := r.collection.InsertOne(context.Background(), RoPreset{
		UserId: i.UserId,
		Label:  i.Label,
		Model:  i.Model,
	})
	if err != nil {
		return nil, err
	}

	return &RoPreset{
		Id:     res.InsertedID.(primitive.ObjectID).Hex(),
		UserId: i.UserId,
		Label:  i.Label,
		Model:  i.Model,
	}, nil
}

// CreatePresets implements RoPresetRepository.
func (r roPresetRepo) CreatePresets(ip BulkCreatePresetInput) (*[]RoPreset, error) {
	var models []mongo.WriteModel
	for i := 0; i < len(ip.BulkData); i++ {
		var cur = ip.BulkData[i]
		var p = RoPreset{
			UserId:    ip.UserId,
			Label:     cur.Label,
			Model:     cur.Model,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}
		models = append(models, mongo.NewInsertOneModel().SetDocument(p))
	}

	res, err := r.collection.BulkWrite(context.Background(), models)
	if err != nil {
		return nil, err
	}

	var presets = []RoPreset{}
	for i := 0; i < int(res.InsertedCount); i++ {
		presets = append(presets, RoPreset{
			UserId: ip.UserId,
			Label:  ip.BulkData[i].Label,
			Model:  ip.BulkData[i].Model,
		})
	}

	return &presets, nil
}

// UpdatePreset implements RoPresetRepository.
func (r roPresetRepo) UpdatePreset(i UpdatePresetInput) (*RoPreset, error) {
	objectId, err := primitive.ObjectIDFromHex(i.Id)
	if err != nil {
		return nil, err
	}

	_, err = r.collection.UpdateByID(context.Background(), objectId, UpdatePresetInput{
		UserId: i.UserId,
		Label:  i.Label,
		Model:  i.Model,
	})
	if err != nil {
		return nil, err
	}

	return &RoPreset{
		Id:     i.Id,
		UserId: i.UserId,
		Label:  i.Label,
		Model:  i.Model,
	}, nil
}

// FindPresetById implements RoPresetRepository.
func (r roPresetRepo) FindPresetById(id string) (*RoPreset, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var data RoPreset
	err = r.collection.FindOne(context.Background(), PartialSearchRoPreset{Id: objectId}).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// FindPresetsByUserId implements RoPresetRepository.
func (r roPresetRepo) FindPresetsByUserId(userId string) (*[]FindPreset, error) {
	cursor, err := r.collection.Find(context.Background(), PartialSearchRoPreset{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var presets = []FindPreset{}

	for cursor.Next(context.Background()) {
		var res RoPreset
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}
		presets = append(presets, FindPreset{
			Id:    res.Id,
			Label: res.Label,
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &presets, nil
}
