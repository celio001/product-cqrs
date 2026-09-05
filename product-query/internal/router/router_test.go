package router

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
	return product.Product{ID: id.String()}, nil
}

var _ product_service.ProductServiceInterface = mockProductService{}

func TestSetupRoutes(t *testing.T) {
	router := SetupRoutes(chi.NewRouter(), mockProductService{})

	tests := []struct {
		name string
		path string
		want int
	}{
		{name: "health", path: "/health", want: http.StatusOK},
		{name: "product", path: "/api/v1/product/" + uuid.NewString(), want: http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			if res.Code != test.want {
				t.Fatalf("expected status %d, got %d", test.want, res.Code)
			}
		})
	}
}
