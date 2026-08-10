package fiber

import (
	brandsSvc "github.com/celio001/product-command/internal/modules/brands/service"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	product_service "github.com/celio001/product-command/internal/modules/product/service"
	"github.com/gofiber/fiber/v3"
)

type HttpServer struct {
	app           *fiber.App
	brandsSvc     brandsSvc.BrandSvcInterface
	categoriesSvc categories_service.CategoriesSvcInterface
	productSvc    product_service.ProductSvcInterface
}

func CreateApp(brandsSvc brandsSvc.BrandSvcInterface, CategoriesSvc categories_service.CategoriesSvcInterface, productSvc product_service.ProductSvcInterface) HttpServer {
	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	httpServer := HttpServer{
		app:           app,
		brandsSvc:     brandsSvc,
		categoriesSvc: CategoriesSvc,
		productSvc: productSvc,
	}

	return httpServer
}
