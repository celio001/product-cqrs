package product_router

import (
	product_handler "github.com/celio001/product-cqrs/product-query/internal/modules/product/handler"
	product_service "github.com/celio001/product-cqrs/product-query/internal/modules/product/service"
	"github.com/go-chi/chi"
)

const (
	HandlerPath = "/product"
)

func RegisterRouter(router chi.Router, productSvc product_service.ProductServiceInterface) {
	handler := product_handler.NewProductHandler(productSvc)
	router.Get("/{id}", handler.GetProductByIDHandler)
}
