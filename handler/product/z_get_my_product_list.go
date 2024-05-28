package product_handler

import (
	"encoding/json"
	"net/http"
	"ro-backend/core"
)

type MyProductRequest struct {
	Skip int `json:"skip"`
	Take int `json:"take"`
}

func (p productHandler) GetMyProductList(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	role := r.Header.Get("role")

	var d MyProductRequest
	json.NewDecoder(r.Body).Decode(&d)

	result, err := p.service.GetMyProductList(userId, role, d.Skip, d.Take)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, result)
}
