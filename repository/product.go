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
	StoreId     *string
	ItemId      *int
	BundleId    *string
	IsPublished *bool
	Type        *int
	SubType     *int
	Name        *string
	ExpDate     *time.Time
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

type ProductStore struct {
	Name           string `bson:"name" json:"name"`
	Description    string `bson:"description" json:"description"`
	Rating         int    `bson:"rating" json:"rating"`
	Fb             string `bson:"fb" json:"fb"`
	Character_name string `bson:"character_name" json:"characterName"`
}

type SearchProductsOutput struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StoreId     primitive.ObjectID `bson:"store_id" json:"storeId"`
	ItemId      int                `bson:"item_id" json:"itemId"`
	Desc        string             `bson:"desc" json:"desc"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Refine      int                `bson:"refine" json:"refine"`
	EnchantIds  []int              `bson:"enchant_ids" json:"enchantIds"`
	CardIds     []int              `bson:"card_ids" json:"cardIds"`
	Opts        []string           `bson:"opts" json:"opts"`
	Baht        float64            `bson:"baht" json:"baht"`
	Zeny        float64            `bson:"zeny" json:"zeny"`
	ExpDate     time.Time          `bson:"exp_date" json:"expDate"`
	IsPublished bool               `bson:"is_published" json:"isPublished"`
	Store       ProductStore       `bson:"store" json:"store"`
}

type PartialSearchProductsOutput struct {
	Items     []SearchProductsOutput `bson:"items" json:"items"`
	TotalItem int                    `bson:"total_item" json:"totalItem"`
	Skip      int                    `json:"skip"`
	Limit     int                    `json:"limit"`
}

type UpdateProductInput struct {
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

type RawUpdateProductInput struct {
	RawId string
	UpdateProductInput
}

func (p RawUpdateProductInput) toUpdateModel(updatedItem time.Time) UpdateProductInput {
	return UpdateProductInput{
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

type PatchProductInput struct {
	Zeny        *float64   `bson:"zeny,omitempty"`
	Quantity    *int       `bson:"quantity,omitempty"`
	IsPublished *bool      `bson:"is_published,omitempty"`
	ExpDate     *time.Time `bson:"exp_date,omitempty"`
	UpdatedAt   time.Time  `bson:"updated_at,omitempty"`
}

type RawPatchProductInput struct {
	RawId string
	PatchProductInput
}

func (p RawPatchProductInput) toUpdateModel(updatedItem time.Time) PatchProductInput {
	return PatchProductInput{
		Zeny:        p.Zeny,
		Quantity:    p.Quantity,
		IsPublished: p.IsPublished,
		ExpDate:     p.ExpDate,
		UpdatedAt:   updatedItem,
	}
}

type ProductRepository interface {
	PartialSearchProductList(input PartialSearchProductsInput) (*PartialSearchProductsOutput, error)
	FindByIds(ids []string) ([]Product, error)
	CreateProductList(inputs []Product) ([]Product, error)
	UpdateProductList(storeId primitive.ObjectID, inputs []RawUpdateProductInput) ([]Product, error)
	PatchProductList(storeId primitive.ObjectID, inputs []RawPatchProductInput) ([]Product, error)
	DeleteProductList(storeId primitive.ObjectID, ids []string) error
}
