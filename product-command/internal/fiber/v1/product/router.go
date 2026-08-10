package product_handler

import (
	product_service "github.com/celio001/product-command/internal/modules/product/service"
	"github.com/gofiber/fiber/v3"
)

const (
	HandlerPath = "/product"
)

func RegisterRouter(router fiber.Router, productSvc product_service.ProductSvcInterface){
	
	handler := NewProductHandler(productSvc)
	router.Post("/", handler.CreateProductHandler)
}