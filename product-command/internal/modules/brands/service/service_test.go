package brands_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/celio001/product-command/internal/modules/brands"
	brandsRepo "github.com/celio001/product-command/internal/modules/brands/repository"
	brands_mocks "github.com/celio001/product-command/internal/modules/brands/repository/mocks"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateBrandsSvc(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := brands_mocks.NewMockBrandsRepoInterface(ctrl)
			mockRepo.EXPECT().
				CreateBrand(tt.ctx, tt.brand).
				Return(tt.mockBrandsReturn, tt.mockError).
				Times(1)

			svc := NewBrandSvc(mockRepo)
			brand, err := svc.CreateBrandSvc(tt.ctx, tt.brand)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, brands.Brand{}, brand)
			} else if !tt.expectedError {
				assert.NoError(t, err)
				assert.IsType(t, brand, brands.Brand{})
				assert.Equal(t, brand.Name, brand.Name)
			}
		})
	}
}

func TestSoftDeleteBrandSvc(t *testing.T) {
	logger.Init("product-command", "1.0.0", "development")
	uuid := uuid.New()
	tests := []struct {
		name                      string
		ctx                       context.Context
		mockGetBrandByIDReturn    brands.Brand
		mockGetBrandByIDErrors    error
		mockSoftDeleteBrandErrors error
		expectedError             bool
	}{
		{
			name:                      "Success_SoftDeleteBrandSvc",
			ctx:                       context.Background(),
			mockGetBrandByIDReturn:    brands.Brand{ID: uuid, Name: "Marca-teste", CreatedAt: time.Now()},
			mockGetBrandByIDErrors:    nil,
			mockSoftDeleteBrandErrors: nil,
			expectedError:             false,
		},
		{
			name:                      "Error_SoftDeleteBrandSvc_GetBrandByID",
			ctx:                       context.Background(),
			mockGetBrandByIDReturn:    brands.Brand{},
			mockGetBrandByIDErrors:    brandsRepo.ErrBrandNotFound,
			mockSoftDeleteBrandErrors: nil,
			expectedError:             true,
		},
		{
			name:                      "Error_SoftDeleteBrandRepo",
			ctx:                       context.Background(),
			mockGetBrandByIDReturn:    brands.Brand{ID: uuid, Name: "Marca-teste", CreatedAt: time.Now()},
			mockGetBrandByIDErrors:    nil,
			mockSoftDeleteBrandErrors: brandsRepo.ErrBrandNotFound,
			expectedError:             true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := brands_mocks.NewMockBrandsRepoInterface(ctrl)

			if tt.mockGetBrandByIDErrors != nil {
				mockRepo.EXPECT().
					GetBrandByID(tt.ctx, uuid).
					Return(tt.mockGetBrandByIDReturn, tt.mockGetBrandByIDErrors).
					Times(1)

			} else if tt.mockSoftDeleteBrandErrors != nil {
				mockRepo.EXPECT().
					GetBrandByID(tt.ctx, uuid).
					Return(tt.mockGetBrandByIDReturn, tt.mockGetBrandByIDErrors)
				mockRepo.EXPECT().
					SoftDeleteBrand(tt.ctx, uuid).
					Return(tt.mockSoftDeleteBrandErrors).
					Times(1)
					
			} else {
				mockRepo.EXPECT().
					GetBrandByID(tt.ctx, uuid).
					Return(tt.mockGetBrandByIDReturn, tt.mockGetBrandByIDErrors)
				mockRepo.EXPECT().
					SoftDeleteBrand(tt.ctx, uuid).
					Return(tt.mockSoftDeleteBrandErrors).
					Times(1)
			}

			svc := NewBrandSvc(mockRepo)
			err := svc.SoftDeleteBrandSvc(tt.ctx, uuid)

			if tt.expectedError {
				assert.Error(t, err)
			} else if !tt.expectedError {
				assert.NoError(t, err)
			}
		})
	}
}
