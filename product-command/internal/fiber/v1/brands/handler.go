package brands_handler

import (
	"net/http"

	"github.com/celio001/product-command/internal/brands"
	brandsSvc "github.com/celio001/product-command/internal/brands/service"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/celio001/product-command/pkg/response"
	validate_errors "github.com/celio001/product-command/pkg/validate"
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

func NewBrandHandler(brandSvc brandsSvc.BrandSvcInterface) BrandsHandlerInterface {
	return &brandsHandler{brandSvc: brandSvc}
}

func (s *brandsHandler) CreateBrandHandler(c fiber.Ctx) error {

	createBrandRequest := CreateBrandRequest{}

	if err := c.Bind().Body(&createBrandRequest); err != nil {
		return response.New().
			Status(http.StatusBadRequest).
			Message(err.Error()).
			Error("INVALID_BODY_CREATE_BRAND").
			Send(c)
	}

	validate := validator.New()

	if err := validate.Struct(createBrandRequest); err != nil {
		logger.Error("invalid create brand payload",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_CREATE_BRAND_PAYLOAD"))

		return response.New().
			Status(http.StatusBadRequest).
			Message("ïnvalid request data").
			Error(validate_errors.BrandsValidateError(err)).
			Send(c)
	}

	b, err := s.brandSvc.CreateBrandSvc(c, brands.Brand{Name: createBrandRequest.Name})
	if err != nil {
		return response.New().
			Status(http.StatusInternalServerError).
			Message(err.Error()).
			Error("ERROR_CREATE_BRAND").
			Send(c)
	}

	return response.New().
		Status(http.StatusOK).
		Message("brand successfully created").
		Data(b).
		Send(c)
}
