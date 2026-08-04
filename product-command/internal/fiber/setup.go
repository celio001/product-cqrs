package fiber

import (
	brandsSvc "github.com/celio001/product-command/internal/modules/brands/service"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	"github.com/gofiber/fiber/v3"
)

type HttpServer struct {
	app           *fiber.App
	brandsSvc     brandsSvc.BrandSvcInterface
	categoriesSvc categories_service.CategoriesSvcInterface
}

func CreateApp(brandsSvc brandsSvc.BrandSvcInterface, CategoriesSvc categories_service.CategoriesSvcInterface) HttpServer {
	app := fiber.New(fiber.Config{
		StrictRouting: true,
	})

	httpServer := HttpServer{
		app:           app,
		brandsSvc:     brandsSvc,
		categoriesSvc: CategoriesSvc,
	}

	return httpServer
}
