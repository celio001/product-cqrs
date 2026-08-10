package product_handler

import (
	"net/http"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	"github.com/celio001/product-command/internal/modules/fiscal"
	"github.com/celio001/product-command/internal/modules/inventory"
	"github.com/celio001/product-command/internal/modules/product"
	product_service "github.com/celio001/product-command/internal/modules/product/service"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/celio001/product-command/pkg/response"
	validate_errors "github.com/celio001/product-command/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type productHandler struct {
	productSvc product_service.ProductSvcInterface
}

type ProductHandlerInterface interface {
	CreateProductHandler(c fiber.Ctx) error
}

func NewProductHandler(productSvc product_service.ProductSvcInterface) ProductHandlerInterface {
	return &productHandler{productSvc: productSvc}
}

var validate = validator.New()

func (h *productHandler) CreateProductHandler(c fiber.Ctx) error {
	var request product_dto.CreateProductRequest

	if err := c.Bind().Body(&request); err != nil {
		logger.Error("invalid create product body",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_CREATE_PRODUCT_BODY"))
		return response.New().
			Status(http.StatusBadRequest).
			Message(err.Error()).
			Error("INVALID_BODY_CREATE_PRODUCT").
			Send(c)
	}

	if err := validate.Struct(request); err != nil {
		logger.Error("invalid create product payload",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_CREATE_PRODUCT_PAYLOAD"))

		return response.New().
			Status(http.StatusBadRequest).
			Message("ïnvalid request data").
			Error(validate_errors.ProductValidateError(err)).
			Send(c)
	}
	p := requestToProduct(request)
	i := requestToInventory(request)
	f := responseToFiscal(request)
	resp, err := h.productSvc.CreateProductSvc(c, p, i, f)

	if err != nil {
		logger.Error("failed to create product",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_CREATE_PRODUCT"),
			zap.String("product.name", resp.Name),
		)
		return response.New().
			Status(http.StatusInternalServerError).
			Message(err.Error()).
			Error("ERROR_CREATE_PRODUCT").
			Send(c)
	}

	logger.Info("product created successfully",
		zap.String("product.name", resp.Name),
		zap.String("product.id", resp.ID.String()),
		zap.String("event.action", "PRODUCT_CREATED_SUCCESS"))

	return response.New().
		Status(http.StatusCreated).
		Message("Product created successfully").
		Data(resp).
		Send(c)
}

func requestToProduct(r product_dto.CreateProductRequest) product.Product {
	return product.Product{
		BrandID:          r.BrandID,
		CategoryID:       r.CategoryID,
		Name:             r.Name,
		Sku:              r.Sku,
		BarCodeEan:       r.BarCodeEan,
		ShortDescription: r.ShortDescription,
		UnitOfMeasure:    r.UnitOfMeasure,
		CostPrice:        r.CostPrice,
		SalePrice:        r.SalePrice,
		PromotionalPrice: r.PromotionalPrice,
		GrossWeight:      r.GrossWeight,
		NetWeight:        r.NetWeight,
		Height:           r.Height,
		Width:            r.Width,
		Length:           r.Length,
		Status:           r.Status,
	}
}

func requestToInventory(r product_dto.CreateProductRequest) inventory.Inventory {
	return inventory.Inventory{
		LocationAisle:     r.Stock.LocationAisle,
		QuantityAvailable: r.Stock.QuantityAvailable,
		MinimumStock:      r.Stock.MinimumStock,
		MaximumStock:      r.Stock.MaximumStock,
	}
}

func responseToFiscal(r product_dto.CreateProductRequest) fiscal.FiscalData {
	return fiscal.FiscalData{
		NcmCode:    r.Fiscal.NcmCode,
		CestCode:   r.Fiscal.CestCode,
		OriginCode: r.Fiscal.OriginCode,
		IcmsRate:   r.Fiscal.IcmsRate,
		PisRate:    r.Fiscal.PisRate,
		CofinsRate: r.Fiscal.CofinsRate,
		IpiRate:    r.Fiscal.IpiRate,
	}
}
