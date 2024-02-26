package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPresetTagRepository(c *mongo.Collection) PresetTagRepository {
	return presetTagRepo{c: c}
}

type presetTagRepo struct {
	c *mongo.Collection
}

func (r presetTagRepo) FindTagById(id string) (*PresetTag, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var p PresetTag
	err = r.c.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r presetTagRepo) FindByPresetIds(ids []string) ([]PresetTag, error) {
	cursor, err := r.c.Find(context.Background(), bson.M{
		"preset_id": bson.M{
			"$in": ids,
		},
	})
	if err != nil {
		return nil, err
	}

	tags := []PresetTag{}
	err = cursor.All(context.Background(), &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r presetTagRepo) PartialSearchTags(i PartialSearchTagsInput, skip, limit int) (*PartialSearchTagsResult, error) {
	total, err := r.c.CountDocuments(context.Background(), i)
	if err != nil {
		return nil, err
	}

	fOpts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{
		{Key: "total_like", Value: -1},
		{Key: "created_at", Value: -1},
	})
	res, err := r.c.Find(context.Background(), i, fOpts)
	if err != nil {
		return nil, err
	}

	items := []PresetTag{}
	err = res.All(context.Background(), &items)
	if err != nil {
		return nil, err
	}

	return &PartialSearchTagsResult{
		Items: items,
		Total: int(total),
	}, nil
}

func (r presetTagRepo) CreateTags(i CreateTagInput) ([]string, error) {
	var tags = []interface{}{}
	now := time.Now()
	for _, tag := range i.Tags {
		tags = append(tags, PresetTag{
			PublisherId: i.PublisherId,
			Tag:         tag,
			ClassId:     i.ClassId,
			PresetId:    i.PresetId,
			Likes:       []string{},
			TotalLike:   0,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	res, err := r.c.InsertMany(context.Background(), tags)
	if err != nil {
		return nil, err
	}

	insertedIDs := []string{}
	for _, id := range res.InsertedIDs {
		insertedIDs = append(insertedIDs, id.(primitive.ObjectID).Hex())
	}

	return insertedIDs, nil
}

func (r presetTagRepo) DeleteTag(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.c.DeleteOne(context.Background(), bson.M{"_id": objId})

	return err
}

func (r presetTagRepo) LikeTag(i LikeTagInput) error {
	objId, err := primitive.ObjectIDFromHex(i.Id)
	if err != nil {
		return err
	}

	_, err = r.c.UpdateByID(context.Background(), objId, bson.M{
		"$addToSet": bson.M{
			"likes": i.UserId,
		},
		"$set": PartialUpdateTagInput{
			TotalLike: i.TotalLike,
			UpdatedAt: time.Now(),
		},
	})

	return err
}

func (r presetTagRepo) UnLikeTag(i LikeTagInput) error {
	objId, err := primitive.ObjectIDFromHex(i.Id)
	if err != nil {
		return err
	}

	_, err = r.c.UpdateByID(context.Background(), objId, bson.M{
		"$pull": bson.M{
			"likes": i.UserId,
		},
		"$set": PartialUpdateTagInput{
			TotalLike: i.TotalLike,
			UpdatedAt: time.Now(),
		},
	})

	return err
}
