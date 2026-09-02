package product_handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/celio001/product-cqrs/product-query/internal/modules/product"
	product_service "github.com/celio001/product-cqrs/product-query/internal/modules/product/service"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type mockProductService struct{}

func (m mockProductService) GetProductByIDSvc(ctx context.Context, id uuid.UUID) (product.Product, error) {
	return product.Product{
		ID:   id.String(),
		Name: "Produto Teste",
	}, nil
}

var _ product_service.ProductServiceInterface = mockProductService{}

func TestGetProductByIDHandler(t *testing.T) {
	id := uuid.New()
	handler := NewProductHandler(mockProductService{})

	r := chi.NewRouter()
	r.Get("/product/{id}", handler.GetProductByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/product/"+id.String(), nil)
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}
