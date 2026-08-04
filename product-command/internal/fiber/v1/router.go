package v1

import (
	brands_handler "github.com/celio001/product-command/internal/fiber/v1/brands"
	category_handler "github.com/celio001/product-command/internal/fiber/v1/category"
	brandsSvc "github.com/celio001/product-command/internal/modules/brands/service"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	"github.com/gofiber/fiber/v3"
)

const (
	HandlerPath = "/v1"
)

func RegisterRouter(router fiber.Router, brandSvc brandsSvc.BrandSvcInterface, CategoriesSvc categories_service.CategoriesSvcInterface) {

	brandsRouters := router.Group(brands_handler.HandlerPath)
	brands_handler.RegisterRouter(brandsRouters, brandSvc)

	categoriesRouters := router.Group(category_handler.HandlerPath)
	category_handler.RegisterRouter(categoriesRouters, CategoriesSvc)
}
