package service

import "ro-backend/repository"

type ProductService interface {
	PartialSearchProductList(input repository.PartialSearchProductsInput) (*repository.PartialSearchProductsOutput, error)
	GetMyProductList(userId, role string, skip, limit int) (*repository.PartialSearchProductsOutput, error)
	CreateProductList(userId, role string, inputs []repository.Product) ([]repository.Product, error)
	UpdateProductList(userId, role string, inputs []repository.RawProductInput) ([]repository.Product, error)
	DeleteProductList(userId string, ids []string) error
}
