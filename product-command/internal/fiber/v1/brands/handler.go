package brands_handler

import (
	"net/http"

	brandsSvc "github.com/celio001/product-command/internal/brands/service"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/celio001/product-command/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type brandsHandler struct {
	brandSvc brandsSvc.BrandSvcInterface
}

type BrandsHandlerInterface interface {
	CreateBrandHandler(c fiber.Ctx) error
}

var validate *validator.Validate

func NewBrandHandler(brandSvc brandsSvc.BrandSvcInterface) BrandsHandlerInterface {
	return &brandsHandler{brandSvc: brandSvc}
}

func (s *brandsHandler) CreateBrandHandler(c fiber.Ctx) error {

	createBrandRequest := CreateBrandRequest{}

	if err := c.Bind().Body(createBrandRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	validate := validator.New()

	if err := validate.Struct(createBrandRequest); err != nil {
		logger.Error("invalid create brand pauload",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_CREATE_BRAND_PAYLOAD"))
		return response.New().Status(http.StatusOK)
	}

	s.brandSvc.CreateBrandSvc()
}
