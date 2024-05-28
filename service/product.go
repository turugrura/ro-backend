package service

import "ro-backend/repository"

type ProductService interface {
	PartialSearchProductList(input repository.PartialSearchProductsInput) (*repository.PartialSearchProductsOutput, error)
	GetMyProductList(userId, role string, skip, limit int) (*repository.PartialSearchProductsOutput, error)
	CreateProductList(userId, role string, inputs []repository.Product) ([]repository.Product, error)
	UpdateProductList(userId, role string, inputs []repository.RawUpdateProductInput) ([]repository.Product, error)
	PatchProductList(userId, role string, inputs []repository.RawPatchProductInput) ([]repository.Product, error)
	RenewExpDateProductList(userId, role string, ids []string) ([]repository.Product, error)
	DeleteProductList(userId string, ids []string) error
}
