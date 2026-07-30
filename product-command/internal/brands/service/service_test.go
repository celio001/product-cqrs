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

func TestCreateBrands_Success(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBrandsRepo := brands_mocks.NewMockBrandsRepoInterface(ctrl)
	uuid := uuid.New()

	brand := brands.Brand{Name: "Marca-teste"}

	mockBrandsRepo.EXPECT().CreateBrand(ctx, brand).Return(brands.Brand{ID: uuid, Name: "Marca-teste", CreatedAt: time.Now()}, nil)

	svc := NewBrandSvc(mockBrandsRepo)
	brand, err := svc.CreateBrandSvc(ctx, brand)

	assert.NoError(t, err)
	assert.IsType(t, brand, brands.Brand{})
	assert.Equal(t, brand.Name, brand.Name)
}

func TestCreateBrands_Error(t *testing.T) {
	logger.Init("product-command", "1.0.0", "development")
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBrandsRepo := brands_mocks.NewMockBrandsRepoInterface(ctrl)

	brand := brands.Brand{Name: "Marca-teste"}

	mockBrandsRepo.EXPECT().CreateBrand(ctx, brand).Return(brands.Brand{}, errors.New("error create bands in database"))

	svc := NewBrandSvc(mockBrandsRepo)
	brand, err := svc.CreateBrandSvc(ctx, brand)

	assert.Error(t, err)
	assert.Equal(t, brands.Brand{}, brand)
}
