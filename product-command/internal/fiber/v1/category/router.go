package category_handler

import (
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	"github.com/gofiber/fiber/v3"
)

const (
	HandlerPath = "/categories"
)

func RegisterRouter(router fiber.Router, CategoriesSvc categories_service.CategoriesSvcInterface) {
	
	handler := NewCategoriesHandler(CategoriesSvc)
	router.Post("/", handler.CreateCategoriesHandler)
}
