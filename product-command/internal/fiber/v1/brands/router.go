package brands_handler

import (
	brandsSvc "github.com/celio001/product-command/internal/brands/service"
	"github.com/gofiber/fiber/v3"
)

const (
	HandlerPath = "/brands"
)

func RegisterRouter(router fiber.Router, brandSvc brandsSvc.BrandSvcInterface) {
	
	handler := NewBrandHandler(brandSvc)
	router.Post("/", handler.CreateBrandHandler)
}
