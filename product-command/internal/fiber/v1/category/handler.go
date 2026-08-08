package category_handler

import (
	"errors"
	"net/http"

	"github.com/celio001/product-command/internal/modules/categories"
	categories_repository "github.com/celio001/product-command/internal/modules/categories/repository"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/celio001/product-command/pkg/response"
	validate_errors "github.com/celio001/product-command/pkg/validate"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoriesHandler struct {
	categoriesSvc categories_service.CategoriesSvcInterface
}

type CategoriesHandlerInterface interface {
	CreateCategoriesHandler(c fiber.Ctx) error
	SoftDeleteCategoriesHandler(c fiber.Ctx) error
}

func NewCategoriesHandler(categoriesSvc categories_service.CategoriesSvcInterface) CategoriesHandlerInterface {
	return &categoriesHandler{categoriesSvc: categoriesSvc}
}

func (h categoriesHandler) CreateCategoriesHandler(c fiber.Ctx) error {

	var categorie CreateCategoriesRequest

	if err := c.Bind().Body(&categorie); err != nil {
		logger.Error("invalid create category body",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_CREATE_CATEGORY_BODY"))
		return response.New().
			Status(http.StatusBadRequest).
			Message(err.Error()).
			Error("INVALID_CATEGORY_CREATE_CATEGORY").
			Send(c)
	}

	validate := validator.New()

	if err := validate.Struct(categorie); err != nil {
		logger.Error("invalid create category payload",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_CREATE_CATEGORY_PAYLOAD"))

		return response.New().
			Status(http.StatusBadRequest).
			Message("ïnvalid request data").
			Error(validate_errors.CategoriesValidateError(err)).
			Send(c)
	}

	category, err := h.categoriesSvc.CreateCategorySvc(c, categories.Categories{Name: categorie.Name})
	if err != nil {
		return response.New().
			Status(http.StatusInternalServerError).
			Error("ERROR_CREATE_CATEGORY").
			Send(c)
	}
	categoryReso := CreateCategoryResponse{
		ID:       category.ID.String(),
		ParentID: category.ParentID.String(),
		Name:     category.Name,
	}
	return response.New().
		Status(http.StatusOK).
		Data(categoryReso).
		Send(c)
}

func (h *categoriesHandler) SoftDeleteCategoriesHandler(c fiber.Ctx) error {
	uString := c.Params("id")

	uuid, err := uuid.Parse(uString)
	if err != nil {
		logger.Error("invalid id category",
			zap.String("error.type", "ValidateError"),
			zap.String("error.message", err.Error()),
			zap.String("error.code", "INVALID_UUID_SOFT_DELETE_CATEGORY"))

		return response.New().
			Status(http.StatusBadRequest).
			Message("invalid id brand").
			Error("INVALID_UUID_SOFT_DELETE_BRAND").
			Send(c)
	}

	err = h.categoriesSvc.SoftDeleteCategory(c, uuid)
	if err != nil {
		switch {
		case errors.Is(err, categories_repository.ErrCategoryNotFound):
			return response.New().
				Status(http.StatusBadRequest).
				Message("category id not found").
				Error("ERROR_CATEGORY_NOT_FOUND").
				Send(c)
		default:
			return response.New().
				Status(http.StatusInternalServerError).
				Message("error delete category").
				Error("ERROR_DELETE_CATEGORY").
				Send(c)
		}
	}

	return response.New().
		Status(http.StatusOK).
		Message("category successfully deleted").
		Send(c)
}
