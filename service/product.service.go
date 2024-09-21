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

	now := time.Now()
	input.ProductFiltering.ExpDate = &now

	published := true
	input.IsPublished = &published

	closelyExpired := 1
	input.ProductSorting.ExpDate = &closelyExpired

	lowestM := 1
	input.ProductSorting.M = &lowestM

	lowestZeny := 1
	input.ProductSorting.Baht = &lowestZeny

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

	closelyExpired := 1

	return p.pRepo.PartialSearchProductList(repository.PartialSearchProductsInput{
		ProductFiltering: repository.ProductFiltering{
			StoreId: &storeId,
		},
		ProductSorting: repository.ProductSorting{
			ExpDate: &closelyExpired,
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

	exp := p.getNextExpDate(role)

	for i := 0; i < len(inputs); i++ {
		v := &inputs[i]
		v.StoreId = store.Id
		v.IsPublished = true
		v.ExpDate = exp
	}

	return p.pRepo.CreateProductList(inputs)
}

func (p productService) UpdateProductList(userId, role string, inputs []repository.RawUpdateProductInput) ([]repository.Product, error) {
	store, err := p.sRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	return p.pRepo.UpdateProductList(store.Id, inputs)
}

func (p productService) PatchProductList(userId, role string, inputs []repository.RawPatchProductInput) ([]repository.Product, error) {
	patchInputs := []repository.RawPatchProductInput{}
	for _, p := range inputs {
		patchInputs = append(patchInputs, repository.RawPatchProductInput{
			RawId: p.RawId,
			PatchProductInput: repository.PatchProductInput{
				Zeny:        p.Zeny,
				Quantity:    p.Quantity,
				IsPublished: p.IsPublished,
			},
		})
	}

	store, err := p.sRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	return p.pRepo.PatchProductList(store.Id, patchInputs)
}

func (p productService) RenewExpDateProductList(userId, role string, ids []string) ([]repository.Product, error) {
	patchInputs := []repository.RawPatchProductInput{}
	nexExpDate := p.getNextExpDate(role)

	for _, id := range ids {
		patchInputs = append(patchInputs, repository.RawPatchProductInput{
			RawId: id,
			PatchProductInput: repository.PatchProductInput{
				ExpDate: &nexExpDate,
			},
		})
	}

	store, err := p.sRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return nil, err
	}

	return p.pRepo.PatchProductList(store.Id, patchInputs)
}

func (p productService) DeleteProductList(userId string, ids []string) error {
	store, err := p.sRepo.FindStoreByOwnerId(userId)
	if err != nil {
		return err
	}

	return p.pRepo.DeleteProductList(store.Id, ids)
}

func (p productService) getNextExpDate(role string) time.Time {
	exp := time.Now()
	if role == repository.UserRole.Admin {
		exp = exp.Add(time.Hour * 24 * 7)
	} else {
		exp = exp.Add(time.Hour * 24 * 2)
	}

	return exp
}
