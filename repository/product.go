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
	EnchantIds  []int              `bson:"enchant_ids" json:"enchantIds"`
	Opts        []string           `bson:"opts" json:"opts"`
	Baht        float64            `bson:"baht" json:"baht"`
	M           float64            `bson:"m" json:"m"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Type        int                `bson:"type" json:"type"`
	SubType     int                `bson:"sub_type" json:"subType"`
	IsPublished bool               `bson:"is_published" json:"isPublished"`
	ExpDate     time.Time          `bson:"exp_date" json:"expDate"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
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
	Id          string    `bson:"_id,omitempty"`
	ItemId      int       `bson:"item_id,omitempty"`
	BundleId    string    `bson:"bundle_id,omitempty"`
	Name        string    `bson:"name,omitempty"`
	Desc        string    `bson:"desc,omitempty"`
	EnchantIds  []int     `bson:"enchant_ids,omitempty"`
	Opts        []string  `bson:"opts,omitempty"`
	Baht        float64   `bson:"baht,omitempty"`
	M           float64   `bson:"m,omitempty"`
	Quantity    int       `bson:"quantity,omitempty"`
	Type        int       `bson:"type,omitempty"`
	SubType     int       `bson:"sub_type,omitempty"`
	IsPublished bool      `bson:"is_published,omitempty"`
	ExpDate     time.Time `bson:"exp_date,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty"`
}

type ProductRepository interface {
	PartialSearchProductList(input PartialSearchProductsInput) (*PartialSearchProductsOutput, error)
	CreateProductList(inputs []Product) ([]Product, error)
	UpdateProductList(inputs []PatchProductInput) ([]Product, error)
	DeleteProductList(ids []string) error
}
