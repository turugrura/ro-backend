package product_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
)

type CreateProductRequest struct {
	ItemId      int      `json:"itemId,omitempty"`
	BundleId    string   `json:"bundleId,omitempty"`
	Name        string   `json:"name,omitempty"`
	Desc        string   `json:"desc,omitempty"`
	EnchantIds  []int    `json:"enchantIds,omitempty"`
	Opts        []string `json:"opts,omitempty"`
	Baht        float64  `json:"baht,omitempty"`
	M           float64  `json:"m,omitempty"`
	Quantity    int      `json:"quantity,omitempty"`
	Type        int      `json:"type,omitempty"`
	SubType     int      `json:"subType,omitempty"`
	IsPublished bool     `json:"isPublished,omitempty"`
}

func (r CreateProductRequest) toCreateInput() repository.Product {
	return repository.Product{
		ItemId:      r.ItemId,
		BundleId:    r.BundleId,
		Name:        r.Name,
		Desc:        r.Desc,
		EnchantIds:  r.EnchantIds,
		Opts:        r.Opts,
		Baht:        r.Baht,
		M:           r.M,
		Quantity:    r.Quantity,
		Type:        r.Type,
		SubType:     r.SubType,
		IsPublished: r.IsPublished,
	}
}

func (r CreateProductRequest) verify() error {
	if r.ItemId == 0 {
		return fmt.Errorf("itemId should be > 0")
	}
	if r.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if r.Type == 0 {
		return fmt.Errorf("type should be > 0")
	}
	if r.SubType == 0 {
		return fmt.Errorf("subType should be > 0")
	}

	return nil
}

func (p productHandler) CreateProductList(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	role := r.Header.Get("role")

	var d []CreateProductRequest
	json.NewDecoder(r.Body).Decode(&d)

	inputs := []repository.Product{}
	for i, product := range d {
		err := product.verify()
		if err != nil {
			core.WriteErr(w, fmt.Sprintf("No %v, %v", i, err.Error()))
			return
		}

		inputs = append(inputs, product.toCreateInput())
	}

	store, err := p.service.CreateProductList(userId, role, inputs)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
