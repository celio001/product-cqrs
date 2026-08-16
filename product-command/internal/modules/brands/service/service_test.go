package brands_service

import (
	"context"
	"errors"
	"testing"
	"time"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	"github.com/celio001/product-command/internal/modules/brands"
	brandsRepo "github.com/celio001/product-command/internal/modules/brands/repository"
	"github.com/celio001/product-command/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type fakeProducer struct{}

func (f *fakeProducer) PublishProductCreated(ctx context.Context, p product_dto.CreateProductResponse) error {
	return nil
}

func (f *fakeProducer) PublishBrandCreated(ctx context.Context, b brands.Brand) error {
	return nil
}

type noopTx struct{ pgx.Tx }

func (n noopTx) Begin(context.Context) (pgx.Tx, error) { return n, nil }
func (n noopTx) Commit(context.Context) error          { return nil }
func (n noopTx) Rollback(context.Context) error        { return nil }
func (n noopTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (n noopTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults { return nil }
func (n noopTx) LargeObjects() pgx.LargeObjects                         { return pgx.LargeObjects{} }
func (n noopTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (n noopTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (n noopTx) Query(context.Context, string, ...any) (pgx.Rows, error) { return nil, nil }
func (n noopTx) QueryRow(context.Context, string, ...any) pgx.Row        { return nil }
func (n noopTx) Conn() *pgx.Conn                                         { return nil }

type fakeBrandRepo struct {
	createBrandFn     func(context.Context, brands.Brand) (brands.Brand, error)
	getBrandByIDFn    func(context.Context, uuid.UUID) (brands.Brand, error)
	softDeleteBrandFn func(context.Context, uuid.UUID) error
	beginTxCalled     bool
}

func (f *fakeBrandRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	f.beginTxCalled = true
	return noopTx{}, nil
}

func (f *fakeBrandRepo) WithTx(tx pgx.Tx) brandsRepo.BrandsRepoInterface {
	return f
}

func (f *fakeBrandRepo) GetBrandByID(ctx context.Context, id uuid.UUID) (brands.Brand, error) {
	if f.getBrandByIDFn != nil {
		return f.getBrandByIDFn(ctx, id)
	}
	return brands.Brand{}, nil
}

func (f *fakeBrandRepo) CreateBrand(ctx context.Context, b brands.Brand) (brands.Brand, error) {
	if f.createBrandFn != nil {
		return f.createBrandFn(ctx, b)
	}
	return b, nil
}

func (f *fakeBrandRepo) SoftDeleteBrand(ctx context.Context, id uuid.UUID) error {
	if f.softDeleteBrandFn != nil {
		return f.softDeleteBrandFn(ctx, id)
	}
	return nil
}

func TestCreateBrandsSvc(t *testing.T) {
	logger.Init("product-command", "1.0.0", "development")

	producer := &fakeProducer{}

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
		},
		{
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
			repo := &fakeBrandRepo{
				createBrandFn: func(ctx context.Context, b brands.Brand) (brands.Brand, error) {
					return tt.mockBrandsReturn, tt.mockError
				},
			}

			svc := NewBrandSvc(repo, producer)
			brand, err := svc.CreateBrandSvc(tt.ctx, tt.brand)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, brands.Brand{}, brand)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.mockBrandsReturn.ID, brand.ID)
			assert.Equal(t, tt.mockBrandsReturn.Name, brand.Name)
			assert.Equal(t, tt.mockBrandsReturn.CreatedAt, brand.CreatedAt)
		})
	}
}

func TestSoftDeleteBrandSvc(t *testing.T) {
	logger.Init("product-command", "1.0.0", "development")

	producer := &fakeProducer{}
	uuidValue := uuid.New()

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
			mockGetBrandByIDReturn:    brands.Brand{ID: uuidValue, Name: "Marca-teste", CreatedAt: time.Now()},
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
			mockGetBrandByIDReturn:    brands.Brand{ID: uuidValue, Name: "Marca-teste", CreatedAt: time.Now()},
			mockGetBrandByIDErrors:    nil,
			mockSoftDeleteBrandErrors: brandsRepo.ErrBrandNotFound,
			expectedError:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeBrandRepo{
				getBrandByIDFn: func(ctx context.Context, id uuid.UUID) (brands.Brand, error) {
					return tt.mockGetBrandByIDReturn, tt.mockGetBrandByIDErrors
				},
				softDeleteBrandFn: func(ctx context.Context, id uuid.UUID) error {
					return tt.mockSoftDeleteBrandErrors
				},
			}

			svc := NewBrandSvc(repo, producer)
			err := svc.SoftDeleteBrandSvc(tt.ctx, uuidValue)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
