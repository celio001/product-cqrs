package router

import (
	product_service "github.com/celio001/product-cqrs/product-query/internal/modules/product/service"
	"github.com/go-chi/chi"
)

func NewSetupRouters(router chi.Router, productSvc product_service.ProductServiceInterface) chi.Router {
	return SetupRoutes(router, productSvc)
}
