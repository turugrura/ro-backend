package service

import (
	"fmt"
	"ro-backend/appError"
	"ro-backend/repository"
	"time"
)

func NewProductService(pRepo repository.ProductRepository, sRepo repository.StoreRepository) ProductService {
	return productService{pRepo: pRepo, sRepo: sRepo}
}

type productService struct {
	pRepo repository.ProductRepository
	sRepo repository.StoreRepository
}

func (p productService) PartialSearchProductList(input repository.PartialSearchProductsInput) (*repository.PartialSearchProductsOutput, error) {
	if input.Limit > 20 {
		input.Limit = 20
	}

	input.ExpDate = time.Now()

	return p.pRepo.PartialSearchProductList(input)
}

func (p productService) GetMyProductList(userId, role string, skip, limit int) (*repository.PartialSearchProductsOutput, error) {
	store, err := p.sRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, fmt.Errorf(appError.ErrStoreNotFound)
	}

	storeId := store.Id.Hex()

	return p.pRepo.PartialSearchProductList(repository.PartialSearchProductsInput{
		ProductFiltering: repository.ProductFiltering{
			StoreId: &storeId,
		},
		Skip:  skip,
		Limit: limit,
	})
}

func (p productService) CreateProductList(userId, role string, inputs []repository.Product) ([]repository.Product, error) {
	store, err := p.sRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, fmt.Errorf("store not found")
	}

	exp := time.Now()
	if role == repository.UserRole.Admin {
		exp = exp.Add(time.Hour * 24 * 7)
	} else {
		exp = exp.Add(time.Hour * 24 * 2)
	}

	for i := 0; i < len(inputs); i++ {
		v := &inputs[i]
		v.StoreId = store.Id
		v.IsPublished = true
		v.ExpDate = exp
	}

	return p.pRepo.CreateProductList(inputs)
}

func (p productService) UpdateProductList(userId, role string, inputs []repository.RawProductInput) ([]repository.Product, error) {
	ids := []string{}
	for _, p := range inputs {
		ids = append(ids, p.RawId)
	}

	err := p.validateProducts(ids)
	if err != nil {
		return nil, err
	}

	return p.pRepo.UpdateProductList(inputs)
}

func (p productService) DeleteProductList(userId string, ids []string) error {
	err := p.validateProducts(ids)
	if err != nil {
		return err
	}

	return p.pRepo.DeleteProductList(ids)
}

func (p productService) validateProducts(ids []string) error {
	myProducts, err := p.pRepo.FindByIds(ids)
	if err != nil {
		return err
	}

	if len(myProducts) != len(ids) {
		return fmt.Errorf(appError.ErrBadInput)
	}

	return nil
}
