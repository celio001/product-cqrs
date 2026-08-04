package fiber

import (
	brandsSvc "github.com/celio001/product-command/internal/modules/brands/service"
	"github.com/gofiber/fiber/v3"
)

type HttpServer struct {
	app       *fiber.App
	brandsSvc brandsSvc.BrandSvcInterface
}

func CreateApp(brandsSvc brandsSvc.BrandSvcInterface) HttpServer {
	app := fiber.New(fiber.Config{
		StrictRouting:         true,
	})

	httpServer := HttpServer{
		app:       app,
		brandsSvc: brandsSvc,
	}

	return httpServer
}
