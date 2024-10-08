package product_handler

import (
	"net/http"
	"ro-backend/service"
)

func NewProductHandler(s service.ProductService) ProductHandler {
	return productHandler{service: s}
}

type ProductHandler interface {
	SearchProductList(w http.ResponseWriter, r *http.Request)
	GetMyProductList(w http.ResponseWriter, r *http.Request)
	CreateProductList(w http.ResponseWriter, r *http.Request)
	UpdateProductList(w http.ResponseWriter, r *http.Request)
	PatchProductList(w http.ResponseWriter, r *http.Request)
	RenewExpDateProductList(w http.ResponseWriter, r *http.Request)
	DeleteProductList(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	service service.ProductService
}
