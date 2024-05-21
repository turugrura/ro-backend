package product_handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/core"
	"ro-backend/repository"
)

type ProductFiltering struct {
	StoreId  *string `json:"storeId,omitempty"`
	ItemId   *int    `json:"itemId,omitempty"`
	BundleId *string `json:"bundleId,omitempty"`
	Type     *int    `json:"type,omitempty"`
	SubType  *int    `json:"subType,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type SearchRequest struct {
	ProductFiltering
	Skip int `json:"skip"`
	Take int `json:"take"`
}

func (s SearchRequest) toSearchInput() repository.PartialSearchProductsInput {
	return repository.PartialSearchProductsInput{
		ProductFiltering: repository.ProductFiltering{
			StoreId:  s.StoreId,
			Name:     s.Name,
			ItemId:   s.ItemId,
			BundleId: s.BundleId,
			Type:     s.Type,
			SubType:  s.SubType,
		},
		Skip:  s.Skip,
		Limit: s.Take,
	}
}

type SearchPresetTagsResponse struct {
	Items      []repository.Product `json:"items"`
	TotalItems int                  `json:"totalItem"`
	Skip       int                  `json:"skip"`
	Take       int                  `json:"take"`
}

func (p productHandler) SearchProductList(w http.ResponseWriter, r *http.Request) {
	var d SearchRequest
	json.NewDecoder(r.Body).Decode(&d)

	result, err := p.service.PartialSearchProductList(d.toSearchInput())
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	response := SearchPresetTagsResponse{
		Items:      result.Items,
		TotalItems: result.TotalItem,
		Skip:       result.Skip,
		Take:       result.Limit,
	}

	core.WriteOK(w, response)
}
