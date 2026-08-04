package v1

import (
	brandsSvc "github.com/celio001/product-command/internal/modules/brands/service"
	brands_handler "github.com/celio001/product-command/internal/fiber/v1/brands"
	"github.com/gofiber/fiber/v3"
)

const (
	HandlerPath = "/v1"
)

func RegisterRouter(router fiber.Router, brandSvc brandsSvc.BrandSvcInterface) {
	
	brandsRouters := router.Group(brands_handler.HandlerPath)
	brands_handler.RegisterRouter(brandsRouters, brandSvc)
}
