package product_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	"github.com/celio001/product-command/internal/modules/fiscal"
	"github.com/celio001/product-command/internal/modules/inventory"
	"github.com/celio001/product-command/internal/modules/product"
	product_repository "github.com/celio001/product-command/internal/modules/product/repository"
	product_service_mocks "github.com/celio001/product-command/internal/modules/product/service/mocks"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type handlerResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   json.RawMessage `json:"error"`
}

func TestMain(m *testing.M) {
	logger.Init("product-command", "test", "test")
	m.Run()
}

func newProductTestApp(service *product_service_mocks.MockProductSvcInterface) *fiber.App {
	app := fiber.New()
	handler := NewProductHandler(service)
	app.Post("/products", handler.CreateProductHandler)
	app.Delete("/products/:id", handler.SoftDeleteProductHandler)
	return app
}

func validCreateRequest() product_dto.CreateProductRequest {
	return product_dto.CreateProductRequest{
		Name:          "Produto teste",
		Sku:           "SKU-001",
		UnitOfMeasure: "UN",
		CostPrice:     10,
		SalePrice:     15,
		Status:        "ACTIVE",
		Stock: product_dto.InventoryProductRequest{
			LocationAisle:     "A1",
			QuantityAvailable: 10,
			MinimumStock:      2,
		},
		Fiscal: product_dto.FiscalProductRequest{
			NcmCode: "12345678",
		},
	}
}

func decodeHandlerResponse(t *testing.T, response *http.Response) handlerResponse {
	t.Helper()
	var body handlerResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	return body
}

func newJSONRequest(method, target string, body string) *http.Request {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func responseError(t *testing.T, body handlerResponse) string {
	t.Helper()
	var message string
	require.NoError(t, json.Unmarshal(body.Error, &message))
	return message
}

func TestCreateProductHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := product_service_mocks.NewMockProductSvcInterface(ctrl)
	app := newProductTestApp(service)
	request := validCreateRequest()
	brandID := uuid.New()
	categoryID := uuid.New()
	productID := uuid.New()
	createdAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	response := product_dto.CreateProductResponse{ID: productID, Name: request.Name, Sku: request.Sku, CreatedAt: createdAt}

	service.EXPECT().CreateProductSvc(
		gomock.Any(),
		product.Product{
			BrandID: brandID, CategoryID: categoryID, Name: request.Name, Sku: request.Sku,
			UnitOfMeasure: request.UnitOfMeasure, CostPrice: request.CostPrice, SalePrice: request.SalePrice,
			Status: request.Status,
		},
		inventory.Inventory{LocationAisle: "A1", QuantityAvailable: 10, MinimumStock: 2},
		fiscal.FiscalData{NcmCode: "12345678"},
	).Return(response, nil)

	request.BrandID = brandID
	request.CategoryID = categoryID
	body, err := json.Marshal(request)
	require.NoError(t, err)
	httpResponse, err := app.Test(newJSONRequest(http.MethodPost, "/products", string(body)))

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
	decoded := decodeHandlerResponse(t, httpResponse)
	assert.Equal(t, http.StatusCreated, decoded.Status)
	assert.Equal(t, "Product created successfully", decoded.Message)
	var decodedProduct product_dto.CreateProductResponse
	require.NoError(t, json.Unmarshal(decoded.Data, &decodedProduct))
	assert.Equal(t, productID, decodedProduct.ID)
	assert.Equal(t, request.Name, decodedProduct.Name)
	assert.Equal(t, request.Sku, decodedProduct.Sku)
	assert.Equal(t, createdAt, decodedProduct.CreatedAt)
}

func TestCreateProductHandlerInvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := product_service_mocks.NewMockProductSvcInterface(ctrl)
	app := newProductTestApp(service)

	httpResponse, err := app.Test(newJSONRequest(http.MethodPost, "/products", "{"))

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
	decoded := decodeHandlerResponse(t, httpResponse)
	assert.Equal(t, "INVALID_BODY_CREATE_PRODUCT", responseError(t, decoded))
}

func TestCreateProductHandlerInvalidPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := product_service_mocks.NewMockProductSvcInterface(ctrl)
	app := newProductTestApp(service)

	httpResponse, err := app.Test(newJSONRequest(http.MethodPost, "/products", `{"sku":"SKU-001"}`))

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
	decoded := decodeHandlerResponse(t, httpResponse)
	assert.Equal(t, http.StatusBadRequest, decoded.Status)
	assert.NotEmpty(t, decoded.Error)
}

func TestCreateProductHandlerServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := product_service_mocks.NewMockProductSvcInterface(ctrl)
	app := newProductTestApp(service)
	request := validCreateRequest()
	serviceErr := errors.New("service unavailable")
	service.EXPECT().CreateProductSvc(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(product_dto.CreateProductResponse{}, serviceErr)
	body, err := json.Marshal(request)
	require.NoError(t, err)

	httpResponse, err := app.Test(newJSONRequest(http.MethodPost, "/products", string(body)))

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, httpResponse.StatusCode)
	decoded := decodeHandlerResponse(t, httpResponse)
	assert.Equal(t, "ERROR_CREATE_PRODUCT", responseError(t, decoded))
	assert.Equal(t, serviceErr.Error(), decoded.Message)
}

func TestSoftDeleteProductHandler(t *testing.T) {
	productID := uuid.New()
	tests := []struct {
		name       string
		serviceErr error
		status     int
		errorCode  string
		message    string
	}{
		{name: "success", status: http.StatusOK, message: "product successfully deleted"},
		{name: "not found", serviceErr: product_repository.ErrProductNotFound, status: http.StatusBadRequest, errorCode: "ERROR_PRODUCT_NOT_FOUND", message: "product id not found"},
		{name: "service error", serviceErr: errors.New("delete failed"), status: http.StatusInternalServerError, errorCode: "ERROR_DELETE_PRODUCT", message: "error delete product"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := product_service_mocks.NewMockProductSvcInterface(ctrl)
			app := newProductTestApp(service)
			service.EXPECT().SoftDeleteProductSvc(gomock.Any(), productID).Return(tt.serviceErr)

			httpResponse, err := app.Test(httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil))

			require.NoError(t, err)
			assert.Equal(t, tt.status, httpResponse.StatusCode)
			decoded := decodeHandlerResponse(t, httpResponse)
			assert.Equal(t, tt.status, decoded.Status)
			assert.Equal(t, tt.message, decoded.Message)
			if tt.errorCode != "" {
				assert.Equal(t, tt.errorCode, responseError(t, decoded))
			}
		})
	}
}

func TestSoftDeleteProductHandlerInvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := product_service_mocks.NewMockProductSvcInterface(ctrl)
	app := newProductTestApp(service)

	httpResponse, err := app.Test(httptest.NewRequest(http.MethodDelete, "/products/not-a-uuid", nil))

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
	decoded := decodeHandlerResponse(t, httpResponse)
	assert.Equal(t, "INVALID_UUID_SOFT_DELETE_PRODUCT", responseError(t, decoded))
}
