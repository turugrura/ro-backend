package product_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
)

type PatchProductRequest struct {
	Id          string   `json:"id"`
	Baht        *float64 `json:"baht,omitempty"`
	Zeny        *float64 `json:"zeny,omitempty"`
	Quantity    *int     `json:"quantity,omitempty"`
	IsPublished *bool    `json:"isPublished,omitempty"`
}

func (r PatchProductRequest) verify() error {
	if r.Baht != nil && *r.Baht < 0 {
		return fmt.Errorf("price Baht should be >= 0")
	}
	if r.Zeny != nil && *r.Zeny < 0 {
		return fmt.Errorf("price Zeny should be >= 0")
	}
	if r.Quantity != nil && *r.Quantity < 0 {
		return fmt.Errorf("quantity should be >= 0")
	}

	return nil
}

func (r PatchProductRequest) toUpdateInput() repository.RawPatchProductInput {
	return repository.RawPatchProductInput{
		RawId: r.Id,
		PatchProductInput: repository.PatchProductInput{
			Zeny:        r.Zeny,
			Quantity:    r.Quantity,
			IsPublished: r.IsPublished,
		},
	}
}

func (p productHandler) PatchProductList(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	role := r.Header.Get("role")

	var d []PatchProductRequest
	json.NewDecoder(r.Body).Decode(&d)

	inputs := []repository.RawPatchProductInput{}
	for i, product := range d {
		err := product.verify()
		if err != nil {
			core.WriteErr(w, fmt.Sprintf("No %v, %v", i, err.Error()))
			return
		}

		inputs = append(inputs, product.toUpdateInput())
	}

	store, err := p.service.PatchProductList(userId, role, inputs)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
