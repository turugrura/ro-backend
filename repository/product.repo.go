package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewProductRepository(collection *mongo.Collection) ProductRepository {
	return productRepo{coll: collection}
}

type productRepo struct {
	coll *mongo.Collection
}

func (p productRepo) CreateProductList(inputs []Product) ([]Product, error) {
	now := time.Now()
	models := []interface{}{}

	for i := 0; i < len(inputs); i++ {
		var c = inputs[i]
		var p = c.toCreateProductModel(now)
		models = append(models, p)
	}

	res, err := p.coll.InsertMany(context.Background(), models)
	if err != nil {
		return nil, err
	}

	inserted, err := p.coll.Find(context.Background(), bson.M{
		"_id": bson.M{
			"$in": res.InsertedIDs,
		},
	})
	if err != nil {
		return nil, err
	}

	var products []Product
	err = inserted.All(context.Background(), &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p productRepo) PartialSearchProductList(input PartialSearchProductsInput) (*PartialSearchProductsOutput, error) {
	skip := input.Skip
	limit := input.Limit

	filter := bson.D{}

	if input.Name != nil {
		filter = append(filter, bson.E{Key: "name", Value: primitive.Regex{
			Pattern: fmt.Sprintf("^%v", *input.Name),
			Options: "i",
		}})
	}

	if input.StoreId != nil {
		storeId, err := primitive.ObjectIDFromHex(*input.StoreId)
		if err != nil {
			return nil, err
		}

		filter = append(filter, bson.E{
			Key: "store_id", Value: storeId,
		})
	}

	if input.ItemId != nil {
		filter = append(filter, bson.E{
			Key: "item_id", Value: *input.ItemId,
		})
	}
	if input.Type != nil {
		filter = append(filter, bson.E{
			Key: "type", Value: *input.Type,
		})
	}
	if input.SubType != nil {
		filter = append(filter, bson.E{
			Key: "sub_type", Value: *input.SubType,
		})
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{
		{Key: "m", Value: 1},
		{Key: "baht", Value: 1},
	})
	total, err := p.coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	cursor, err := p.coll.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	items := []Product{}
	err = cursor.All(context.Background(), &items)
	if err != nil {
		return nil, err
	}

	return &PartialSearchProductsOutput{
		Items:     items,
		TotalItem: int(total),
		Skip:      skip,
		Limit:     limit,
	}, nil
}

func (p productRepo) FindByIds(ids []string) ([]Product, error) {
	objIds := []primitive.ObjectID{}
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		objIds = append(objIds, objId)
	}

	cs, err := p.coll.Find(context.Background(), bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	})
	if err != nil {
		return nil, err
	}

	products := []Product{}
	err = cs.All(context.Background(), &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p productRepo) UpdateProductList(inputs []RawProductInput) ([]Product, error) {
	now := time.Now()

	updateModels := []mongo.WriteModel{}
	ids := []string{}
	for i := 0; i < len(inputs); i++ {
		var c = inputs[i]

		objId, err := primitive.ObjectIDFromHex(c.RawId)
		if err != nil {
			return nil, err
		}
		var p = c.toUpdateModel(objId, now)
		ids = append(ids, c.RawId)
		updateOpe := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": objId}).SetUpdate(bson.M{"$set": p})
		updateModels = append(updateModels, updateOpe)
	}

	_, err := p.coll.BulkWrite(context.Background(), updateModels)
	if err != nil {
		return nil, err
	}

	return p.FindByIds(ids)
}

func (p productRepo) DeleteProductList(ids []string) error {
	objIds := []primitive.ObjectID{}
	for _, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}
		objIds = append(objIds, objId)
	}

	_, err := p.coll.DeleteMany(context.Background(), bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	})

	return err
}
