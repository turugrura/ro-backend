package product_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/core"
)

type RenewExpDateProductRequest struct {
	Ids []string `json:"ids"`
}

func (r RenewExpDateProductRequest) verify() error {
	if len(r.Ids) == 0 {
		return fmt.Errorf("ids should not be empty list")
	}

	for _, id := range r.Ids {
		if id == "" {
			return fmt.Errorf("id should not be empty")
		}
	}

	return nil
}

func (p productHandler) RenewExpDateProductList(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	role := r.Header.Get("role")

	var d RenewExpDateProductRequest
	json.NewDecoder(r.Body).Decode(&d)

	err := d.verify()
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	store, err := p.service.RenewExpDateProductList(userId, role, d.Ids)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, store)
}
