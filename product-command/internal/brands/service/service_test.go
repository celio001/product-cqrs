package brands_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/celio001/product-command/internal/brands"
	brands_mocks "github.com/celio001/product-command/internal/brands/repository/mocks"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateBrands(t *testing.T) {
	logger.Init("product-command", "1.0.0", "development")
	tests := []struct {
		name             string
		brand            brands.Brand
		ctx              context.Context
		mockBrandsReturn brands.Brand
		mockError        error
		expected         brands.Brand
		expectedError    bool
	}{
		{
			name:             "Success",
			brand:            brands.Brand{Name: "Marca-teste"},
			ctx:              context.Background(),
			mockBrandsReturn: brands.Brand{ID: uuid.New(), Name: "Marca-teste", CreatedAt: time.Now()},
			mockError:        nil,
			expected:         brands.Brand{ID: uuid.New(), Name: "Marca-teste", CreatedAt: time.Now()},
			expectedError:    false,
		}, {
			name:             "Error",
			brand:            brands.Brand{Name: "Marca-teste"},
			ctx:              context.Background(),
			mockBrandsReturn: brands.Brand{},
			mockError:        errors.New("error create brand"),
			expected:         brands.Brand{},
			expectedError:    true,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := brands_mocks.NewMockBrandsRepoInterface(ctrl)
		mockRepo.EXPECT().
			CreateBrand(tt.ctx ,tt.brand).
			Return(tt.mockBrandsReturn, tt.mockError).
			Times(1)
		
		svc := NewBrandSvc(mockRepo)
		brand, err := svc.CreateBrandSvc(tt.ctx, tt.brand)

		if tt.expectedError {
			assert.Error(t, err)
			assert.Equal(t, brands.Brand{}, brand)
		} else if !tt.expectedError{
			assert.NoError(t, err)
			assert.IsType(t, brand, brands.Brand{})
			assert.Equal(t, brand.Name, brand.Name)
		}
	}
}
