package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	if input.ProductFiltering.ExpDate != nil {
		filter = append(filter, bson.E{
			Key: "exp_date", Value: bson.M{
				"$gte": *input.ProductFiltering.ExpDate,
			},
		})
	}
	if input.ProductFiltering.IsPublished != nil {
		filter = append(filter, bson.E{
			Key: "is_published", Value: *input.ProductFiltering.IsPublished,
		})
	}

	cursor, err := p.coll.Aggregate(context.Background(), bson.A{
		bson.D{
			{Key: "$match", Value: filter},
		},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "store"},
					{Key: "localField", Value: "store_id"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "store"},
				},
			},
		},
		bson.D{
			{Key: "$unwind",
				Value: bson.D{
					{Key: "path", Value: "$store"},
					{Key: "preserveNullAndEmptyArrays", Value: true},
				},
			},
		},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "store_id", Value: 1},
					{Key: "item_id", Value: 1},
					{Key: "desc", Value: 1},
					{Key: "refine", Value: 1},
					{Key: "enchant_ids", Value: 1},
					{Key: "card_ids", Value: 1},
					{Key: "opts", Value: 1},
					{Key: "baht", Value: 1},
					{Key: "zeny", Value: 1},
					{Key: "quantity", Value: 1},
					{Key: "is_published", Value: 1},
					{Key: "exp_date", Value: 1},
					{Key: "store.name", Value: 1},
					{Key: "store.description", Value: 1},
					{Key: "store.rating", Value: 1},
					{Key: "store.fb", Value: 1},
					{Key: "store.character_name", Value: 1},
				},
			},
		},
		bson.D{
			{Key: "$sort",
				Value: bson.D{
					{Key: "m", Value: 1},
					{Key: "baht", Value: 1},
				},
			},
		},
		bson.D{
			{Key: "$facet",
				Value: bson.D{
					{Key: "total_item",
						Value: bson.A{
							bson.D{{Key: `$count`, Value: "count"}},
						},
					},
					{Key: "items",
						Value: bson.A{
							bson.D{{Key: "$skip", Value: skip}},
							bson.D{{Key: "$limit", Value: limit}},
						},
					},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$total_item"}}}},
		bson.D{
			{Key: "$replaceWith",
				Value: bson.D{
					{Key: "total_item", Value: "$total_item.count"},
					{Key: "items", Value: "$items"},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	result := []PartialSearchProductsOutput{}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &PartialSearchProductsOutput{
			Items:     []SearchProductsOutput{},
			TotalItem: 0,
			Skip:      skip,
			Limit:     limit,
		}, nil
	}

	response := result[0]

	response.Skip = skip
	response.Limit = limit

	return &response, nil
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

func (p productRepo) UpdateProductList(storeId primitive.ObjectID, inputs []RawUpdateProductInput) ([]Product, error) {
	now := time.Now()
	updateModels := []mongo.WriteModel{}
	ids := []string{}
	for i := 0; i < len(inputs); i++ {
		var c = inputs[i]

		objId, err := primitive.ObjectIDFromHex(c.RawId)
		if err != nil {
			return nil, err
		}

		var p = c.toUpdateModel(now)
		ids = append(ids, c.RawId)
		updateOpe := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": objId, "store_id": storeId}).SetUpdate(bson.M{"$set": p})
		updateModels = append(updateModels, updateOpe)
	}

	_, err := p.coll.BulkWrite(context.Background(), updateModels)
	if err != nil {
		return nil, err
	}

	return p.FindByIds(ids)
}

func (p productRepo) PatchProductList(storeId primitive.ObjectID, inputs []RawPatchProductInput) ([]Product, error) {
	now := time.Now()
	updateModels := []mongo.WriteModel{}
	ids := []string{}
	for i := 0; i < len(inputs); i++ {
		var c = inputs[i]

		objId, err := primitive.ObjectIDFromHex(c.RawId)
		if err != nil {
			return nil, err
		}
		var p = c.toUpdateModel(now)
		ids = append(ids, c.RawId)
		updateOpe := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": objId, "store_id": storeId}).SetUpdate(bson.M{"$set": p})
		updateModels = append(updateModels, updateOpe)
	}

	_, err := p.coll.BulkWrite(context.Background(), updateModels)
	if err != nil {
		return nil, err
	}

	return p.FindByIds(ids)
}

func (p productRepo) DeleteProductList(storeId primitive.ObjectID, ids []string) error {
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
		"store_id": storeId,
	})

	return err
}
