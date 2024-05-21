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
	models := []interface{}{}
	now := time.Now()
	for i := 0; i < len(inputs); i++ {
		var c = inputs[i]
		var p = Product{
			StoreId:     c.StoreId,
			ItemId:      c.ItemId,
			BundleId:    c.BundleId,
			Name:        c.Name,
			Desc:        c.Desc,
			EnchantIds:  c.EnchantIds,
			Opts:        c.Opts,
			Baht:        c.Baht,
			M:           c.M,
			Quantity:    c.Quantity,
			Type:        c.Type,
			SubType:     c.SubType,
			IsPublished: c.IsPublished,
			ExpDate:     c.ExpDate,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
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

func (p productRepo) UpdateProductList(inputs []PatchProductInput) ([]Product, error) {
	now := time.Now()

	updateModels := []mongo.WriteModel{}
	products := []Product{}
	for i := 0; i < len(inputs); i++ {
		var c = inputs[i]

		objId, err := primitive.ObjectIDFromHex(c.Id)
		if err != nil {
			return nil, err
		}
		var p = Product{
			Id:          objId,
			ItemId:      c.ItemId,
			BundleId:    c.BundleId,
			Name:        c.Name,
			Desc:        c.Desc,
			EnchantIds:  c.EnchantIds,
			Opts:        c.Opts,
			Baht:        c.Baht,
			M:           c.M,
			Quantity:    c.Quantity,
			Type:        c.Type,
			SubType:     c.SubType,
			IsPublished: c.IsPublished,
			ExpDate:     c.ExpDate,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		products = append(products, p)
		updateModels = append(updateModels, &mongo.UpdateOneModel{
			Filter: bson.M{"_id": objId},
			Update: p,
		})
	}

	_, err := p.coll.BulkWrite(context.Background(), updateModels)
	if err != nil {
		return nil, err
	}

	return products, nil
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
