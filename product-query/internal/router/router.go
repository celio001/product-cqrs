package router

import (
	"net/http"

	product_router "github.com/celio001/product-cqrs/product-query/internal/modules/product/router"
	product_service "github.com/celio001/product-cqrs/product-query/internal/modules/product/service"
	"github.com/go-chi/chi"
)

func SetupRoutes(router chi.Router, productSvc product_service.ProductServiceInterface) chi.Router {
	router.Get("/health", healthCheck)

	//v1 
	router.Route("/api/v1", func(route chi.Router) {
		route.Route(product_router.HandlerPath, func(productRoute chi.Router) {
			product_router.SetupProductRoutes(productRoute, productSvc)
		})
	})

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
