package product_handler

import (
	"errors"
	"net/http"

	product_service "github.com/celio001/product-cqrs/product-query/internal/modules/product/service"
	"github.com/celio001/product-cqrs/product-query/pkg/logger"
	"github.com/celio001/product-cqrs/product-query/pkg/response"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type productHandler struct {
	productSvc product_service.ProductServiceInterface
}

type ProductHandlerInterface interface {
	GetProductByIDHandler(w http.ResponseWriter, r *http.Request)
}

func NewProductHandler(productSvc product_service.ProductServiceInterface) ProductHandlerInterface {
	return &productHandler{productSvc: productSvc}
}

func (h *productHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.New().
			Status(http.StatusBadRequest).
			Message("invalid id product").
			Error("INVALID_UUID_PRODUCT").
			Send(w)
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Error("invalid id product",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_UUID_PRODUCT"),
		)
		response.New().
			Status(http.StatusBadRequest).
			Message("invalid id product").
			Error("INVALID_UUID_PRODUCT").
			Send(w)
		return
	}

	product, err := h.productSvc.GetProductByIDSvc(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			response.New().
				Status(http.StatusNotFound).
				Message("product not found").
				Error("PRODUCT_NOT_FOUND").
				Send(w)
			return
		}

		logger.Error("failed to get product",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "ERROR_GET_PRODUCT"),
		)
		response.New().
			Status(http.StatusInternalServerError).
			Message("error get product").
			Error("ERROR_GET_PRODUCT").
			Send(w)
		return
	}

	response.New().
		Status(http.StatusOK).
		Message("product retrieved successfully").
		Data(product).
		Send(w)
}
