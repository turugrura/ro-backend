package product_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ro-backend/core"
)

type DeleteProductListRequest struct {
	Ids []string `json:"ids"`
}

func (r DeleteProductListRequest) verify() error {
	if len(r.Ids) == 0 {
		return fmt.Errorf("length > 0")
	}

	for _, v := range r.Ids {
		if v == "" {
			return fmt.Errorf("element should not be empty")
		}
	}

	return nil
}

func (p productHandler) DeleteProductList(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")

	var d DeleteProductListRequest
	json.NewDecoder(r.Body).Decode(&d)

	err := d.verify()
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	err = p.service.DeleteProductList(userId, d.Ids)
	if err != nil {
		core.WriteErr(w, err.Error())
		return
	}

	core.WriteOK(w, nil)
}
