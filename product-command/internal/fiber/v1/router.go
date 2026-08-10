package v1

import (
	brands_handler "github.com/celio001/product-command/internal/fiber/v1/brands"
	category_handler "github.com/celio001/product-command/internal/fiber/v1/category"
	product_handler "github.com/celio001/product-command/internal/fiber/v1/product"
	brandsSvc "github.com/celio001/product-command/internal/modules/brands/service"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	product_service "github.com/celio001/product-command/internal/modules/product/service"
	"github.com/gofiber/fiber/v3"
)

const (
	HandlerPath = "/v1"
)

func RegisterRouter(router fiber.Router, brandSvc brandsSvc.BrandSvcInterface, CategoriesSvc categories_service.CategoriesSvcInterface, productSvc product_service.ProductSvcInterface) {

	brandsRouters := router.Group(brands_handler.HandlerPath)
	brands_handler.RegisterRouter(brandsRouters, brandSvc)

	categoriesRouters := router.Group(category_handler.HandlerPath)
	category_handler.RegisterRouter(categoriesRouters, CategoriesSvc)

	productRouters := router.Group(product_handler.HandlerPath)
	product_handler.RegisterRouter(productRouters, productSvc)
}
