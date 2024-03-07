package repository

import (
	"context"
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

func (r roPresetRepo) UpdateUserName(userId string, userName string) error {
	_, err := r.collection.UpdateMany(context.Background(), PartialSearchRoPresetForUpdateInput{
		UserId: userId,
	}, bson.M{
		"$set": UpdatePresetInput{
			UserName: userName,
		},
	})

	return err
}

func (r roPresetRepo) UnpublishedPreset(id string) error {
	_, err := r.collection.UpdateOne(context.Background(), IdSearchInput{Id: id}, bson.M{
		"$set": UnPublishPresetInput{
			IsPublished: false,
		},
	})

	return err
}

func (r roPresetRepo) FindPresetByIds(ids []string) ([]RoPreset, error) {
	cs, err := r.collection.Find(context.Background(), bson.M{
		"id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return nil, err
	}

	presets := []RoPreset{}
	err = cs.All(context.Background(), &presets)
	if err != nil {
		return nil, err
	}

	return presets, nil
}

func (r roPresetRepo) PartialSearchPresets(i PartialSearchRoPresetInput) (*PartialSearchRoPresetResult, error) {
	filter := bson.M{}
	if i.ClassId != nil {
		filter["class_id"] = *i.ClassId
	}
	if i.Id != nil {
		filter["id"] = *i.Id
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
		Sort: PresetListSorting{
			UpdatedAt: -1,
		},
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

func (r roPresetRepo) DeletePresetById(id string) (*int, error) {
	res, err := r.collection.DeleteOne(context.Background(), IdSearchInput{Id: id})
	if err != nil {
		return nil, err
	}

	i := int(res.DeletedCount)

	return &i, nil
}

func (r roPresetRepo) CreatePreset(i CreatePresetInput) (*RoPreset, error) {
	id := uuid.NewString()
	_, err := r.collection.InsertOne(context.Background(), RoPreset{
		Id:        id,
		UserId:    i.UserId,
		Label:     i.Label,
		Model:     i.Model,
		ClassId:   i.Model.Class,
		UserName:  i.UserName,
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

func (r roPresetRepo) CreatePresets(ip BulkCreatePresetInput) ([]RoPreset, error) {
	var models []interface{}
	now := time.Now()
	for i := 0; i < len(ip.BulkData); i++ {
		var cur = ip.BulkData[i]
		var p = RoPreset{
			Id:        uuid.NewString(),
			UserId:    ip.UserId,
			Label:     cur.Label,
			Model:     cur.Model,
			ClassId:   cur.Model.Class,
			UserName:  ip.UserName,
			CreatedAt: now,
			UpdatedAt: now,
		}
		models = append(models, p)
	}

	res, err := r.collection.InsertMany(context.Background(), models)
	if err != nil {
		return nil, err
	}

	inserted, err := r.collection.Find(context.Background(), bson.M{
		"_id": bson.M{
			"$in": res.InsertedIDs,
		},
	})
	if err != nil {
		return nil, err
	}

	var presets []RoPreset
	err = inserted.All(context.Background(), &presets)
	if err != nil {
		return nil, err
	}

	return presets, nil
}

func (r roPresetRepo) UpdatePreset(id string, i UpdatePresetInput) error {
	i.UpdatedAt = time.Now()

	if i.Model != nil {
		i.ClassId = i.Model.Class
	}

	_, err := r.collection.UpdateOne(context.Background(), IdSearchInput{Id: id}, bson.M{
		"$set": i,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r roPresetRepo) FindPresetById(i FindPresetByIdInput) (*RoPreset, error) {
	var opt = options.FindOneOptions{}
	if !i.InCludeModel {
		opt.Projection = bson.M{
			"model": 0,
		}
	}

	var data RoPreset
	err := r.collection.FindOne(context.Background(), IdSearchInput{Id: i.Id}, &opt).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
