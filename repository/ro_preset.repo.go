package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRoPresetRepository(collection *mongo.Collection) RoPresetRepository {
	return roPresetRepo{collection: collection}
}

type roPresetRepo struct {
	collection *mongo.Collection
}

// FindPresetByIds implements RoPresetRepository.
func (r roPresetRepo) FindPresetByIds(ids []string) (*[]RoPreset, error) {
	r.collection.Find(context.Background(), bson.M{
		"id": bson.M{
			"$in": ids,
		},
	})

	return nil, nil
}

// PartialSearchPresets implements RoPresetRepository.
func (r roPresetRepo) PartialSearchPresets(i PartialSearchRoPresetInput) (*PartialSearchRoPresetResult, error) {
	filter := bson.M{}
	if i.ClassId != nil {
		filter["class_id"] = *i.ClassId
	}
	if i.Id != nil {
		filter["id"] = *i.Id
	}
	if i.Tag != nil {
		filter["tags"] = bson.M{
			"$in": []string{*i.Tag},
		}
	}
	if i.UserId != nil {
		filter["user_id"] = *i.UserId
	}

	total, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	skip := int64(0)
	if i.Skip != nil {
		skip = int64(*i.Skip)
	}

	take := total
	if i.Take != nil {
		take = int64(*i.Take)
	}

	var opt = options.FindOptions{}
	if !i.InCludeModel {
		opt.Projection = bson.M{
			"model": 0,
		}
	}

	cursor, err := r.collection.Find(context.Background(), filter, &options.FindOptions{
		Projection: opt.Projection,
		Skip:       &skip,
		Limit:      &take,
	})
	if err != nil {
		return nil, err
	}

	items := []RoPreset{}
	err = cursor.All(context.Background(), &items)
	if err != nil {
		return nil, err
	}

	return &PartialSearchRoPresetResult{
		Items: items,
		Total: total,
	}, nil
}

// DeletePresetById implements RoPresetRepository.
func (r roPresetRepo) DeletePresetById(id string) (*int, error) {
	res, err := r.collection.DeleteOne(context.Background(), PartialSearchRoPresetInput{Id: &id})
	if err != nil {
		return nil, err
	}

	i := int(res.DeletedCount)

	return &i, nil
}

// CreatePreset implements RoPresetRepository.
func (r roPresetRepo) CreatePreset(i CreatePresetInput) (*RoPreset, error) {
	id := uuid.NewString()
	_, err := r.collection.InsertOne(context.Background(), RoPreset{
		Id:        id,
		UserId:    i.UserId,
		Label:     i.Label,
		Model:     i.Model,
		ClassId:   i.Model.Class,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &RoPreset{
		Id:     id,
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
			Id:        uuid.NewString(),
			UserId:    ip.UserId,
			Label:     cur.Label,
			Model:     cur.Model,
			ClassId:   cur.Model.Class,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
func (r roPresetRepo) UpdatePreset(i UpdatePresetInput) error {
	i.UpdatedAt = time.Now()

	if i.Model != nil {
		i.ClassId = i.Model.Class
	}

	_, err := r.collection.UpdateOne(context.Background(), PartialSearchRoPresetInput{Id: &i.Id}, bson.M{
		"$set": i,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r roPresetRepo) FindPresetById(i FindPresetByIdInput) (*RoPreset, error) {
	var opt = options.FindOneOptions{
		Projection: bson.M{
			"model": 0,
		},
	}
	fmt.Println("i.InCludeModel", i.InCludeModel)
	if i.InCludeModel {
		opt.Projection = bson.M{
			"model": 1,
		}
	}

	var data RoPreset
	err := r.collection.FindOne(context.Background(), PartialSearchRoPresetInput{Id: &i.Id}, &opt).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
