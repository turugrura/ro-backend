package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StoreId     primitive.ObjectID `bson:"store_id" json:"storeId"`
	ItemId      int                `bson:"item_id" json:"itemId"`
	BundleId    string             `bson:"bundle_id" json:"bundleId"`
	Name        string             `bson:"name" json:"name"`
	Desc        string             `bson:"desc" json:"desc"`
	Refine      int                `bson:"refine" json:"refine"`
	EnchantIds  []int              `bson:"enchant_ids" json:"enchantIds"`
	CardIds     []int              `bson:"card_ids" json:"cardIds"`
	Opts        []string           `bson:"opts" json:"opts"`
	Baht        float64            `bson:"baht" json:"baht"`
	Zeny        float64            `bson:"zeny" json:"zeny"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Type        int                `bson:"type" json:"type"`
	SubType     int                `bson:"sub_type" json:"subType"`
	IsPublished bool               `bson:"is_published" json:"isPublished"`
	ExpDate     time.Time          `bson:"exp_date" json:"expDate"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

func (p Product) toCreateProductModel(createdTime time.Time) Product {
	return Product{
		StoreId:     p.StoreId,
		ItemId:      p.ItemId,
		BundleId:    p.BundleId,
		Name:        p.Name,
		Desc:        p.Desc,
		Refine:      p.Refine,
		EnchantIds:  p.EnchantIds,
		CardIds:     p.CardIds,
		Opts:        p.Opts,
		Baht:        p.Baht,
		Zeny:        p.Zeny,
		Quantity:    p.Quantity,
		Type:        p.Type,
		SubType:     p.SubType,
		IsPublished: p.IsPublished,
		ExpDate:     p.ExpDate,
		CreatedAt:   createdTime,
		UpdatedAt:   createdTime,
	}
}

type ProductFiltering struct {
	StoreId  *string
	ItemId   *int
	BundleId *string
	Type     *int
	SubType  *int
	Name     *string
}

type ProductSorting struct {
	Baht    int
	M       int
	ExpDate time.Time
}

type PartialSearchProductsInput struct {
	ProductFiltering
	ProductSorting
	Skip  int
	Limit int
}

type PartialSearchProductsOutput struct {
	Items     []Product
	TotalItem int
	Skip      int
	Limit     int
}

type PatchProductInput struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	ItemId      int                `bson:"item_id,omitempty"`
	BundleId    string             `bson:"bundle_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Desc        string             `bson:"desc,omitempty"`
	Refine      int                `bson:"refine,omitempty"`
	EnchantIds  []int              `bson:"enchant_ids,omitempty"`
	CardIds     []int              `bson:"card_ids,omitempty"`
	Opts        []string           `bson:"opts,omitempty"`
	Baht        float64            `bson:"baht,omitempty"`
	Zeny        float64            `bson:"zeny,omitempty"`
	Quantity    int                `bson:"quantity,omitempty"`
	Type        int                `bson:"type,omitempty"`
	SubType     int                `bson:"sub_type,omitempty"`
	IsPublished bool               `bson:"is_published,omitempty"`
	ExpDate     time.Time          `bson:"exp_date,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

type RawProductInput struct {
	RawId string
	PatchProductInput
}

func (p RawProductInput) toUpdateModel(objId primitive.ObjectID, updatedItem time.Time) PatchProductInput {
	return PatchProductInput{
		Id:          objId,
		ItemId:      p.ItemId,
		BundleId:    p.BundleId,
		Name:        p.Name,
		Desc:        p.Desc,
		Refine:      p.Refine,
		EnchantIds:  p.EnchantIds,
		CardIds:     p.CardIds,
		Opts:        p.Opts,
		Baht:        p.Baht,
		Zeny:        p.Zeny,
		Quantity:    p.Quantity,
		Type:        p.Type,
		SubType:     p.SubType,
		IsPublished: p.IsPublished,
		ExpDate:     p.ExpDate,
		UpdatedAt:   updatedItem,
	}
}

type ProductRepository interface {
	PartialSearchProductList(input PartialSearchProductsInput) (*PartialSearchProductsOutput, error)
	FindByIds(ids []string) ([]Product, error)
	CreateProductList(inputs []Product) ([]Product, error)
	UpdateProductList(inputs []RawProductInput) ([]Product, error)
	DeleteProductList(ids []string) error
}
