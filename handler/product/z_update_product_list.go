package product_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
)

type UpdateProductRequest struct {
	Id string `json:"id"`
	CreateProductRequest
}

func (r UpdateProductRequest) toUpdateInput() repository.RawUpdateProductInput {
	return repository.RawUpdateProductInput{
		RawId: r.Id,
		UpdateProductInput: repository.UpdateProductInput{
			ItemId:      r.ItemId,
			BundleId:    r.BundleId,
			Name:        r.Name,
			Refine:      r.Refine,
			CardIds:     r.CardIds,
			Desc:        r.Desc,
			EnchantIds:  r.EnchantIds,
			Opts:        r.Opts,
			Baht:        r.Baht,
			Zeny:        r.Zeny,
			Quantity:    r.Quantity,
			Type:        r.Type,
			SubType:     r.SubType,
			IsPublished: r.IsPublished,
		},
	}
}

func (p productHandler) UpdateProductList(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	role := r.Header.Get("role")

	var d []UpdateProductRequest
	json.NewDecoder(r.Body).Decode(&d)

	inputs := []repository.RawUpdateProductInput{}
	for i, product := range d {
		err := product.verify()
		if err != nil {
			core.WriteErr(w, fmt.Sprintf("No %v, %v", i, err.Error()))
			return
		}

		inputs = append(inputs, product.toUpdateInput())
	}

	store, err := p.service.UpdateProductList(userId, role, inputs)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
