package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMovieTranslatorRepository(collection *mongo.Collection) MovieTranslatorRepository {
	return movieTranslatorRepo{c: collection}
}

type movieTranslatorRepo struct {
	c *mongo.Collection
}

func (r movieTranslatorRepo) GetAllEpisodes() ([]MovieInfo, error) {
	cursor, err := r.c.Aggregate(context.Background(), bson.A{
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "title", Value: 1},
					{Key: "season", Value: 1},
					{Key: "episode", Value: 1},
					{Key: "updated_at", Value: 1},
					{Key: "updated_by", Value: 1},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	result := []MovieInfo{}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r movieTranslatorRepo) GetEpisode(ss float32, ep float32) (*MovieTranslator, error) {
	identifier := MovieTranslatorIdentifier{
		Season:  ss,
		Episode: ep,
	}

	var mt MovieTranslator
	err := r.c.FindOne(context.Background(), identifier).Decode(&mt)
	if err != nil {
		return nil, err
	}

	return &mt, nil
}

func (r movieTranslatorRepo) PatchSentence(ss float32, ep float32, sentence PatchSentenceInput) (*MovieTranslator, error) {
	identifier := sentence.toIdentifier(ss, ep)
	updateModel := sentence.toUpdateModel()

	result, err := r.c.UpdateOne(context.Background(), identifier, updateModel)

	if result.ModifiedCount <= 0 || err != nil {
		return nil, err
	}

	var updated MovieTranslator
	err = r.c.FindOne(context.Background(), identifier).Decode(&updated)
	if err != nil {
		return nil, err
	}

	for _, s := range updated.Sentences {
		if s.Index == sentence.Index {
			updated.Sentences = []Sentence{s}
			break
		}
	}

	return &updated, err
}
