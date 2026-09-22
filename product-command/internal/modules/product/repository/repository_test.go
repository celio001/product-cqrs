package product_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/celio001/product-command/internal/database"
	product "github.com/celio001/product-command/internal/modules/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeDB struct {
	queryRowFn func(context.Context, string, ...any) pgx.Row
	execFn     func(context.Context, string, ...any) (pgconn.CommandTag, error)
}

func (f *fakeDB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return f.execFn(ctx, query, args...)
}

func (f *fakeDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unexpected Query call")
}

func (f *fakeDB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return f.queryRowFn(ctx, query, args...)
}

type fakeRow struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	err       error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	*dest[0].(*uuid.UUID) = r.id
	*dest[1].(*time.Time) = r.createdAt
	*dest[2].(*time.Time) = r.updatedAt
	return nil
}

func TestProductRepoCreateProductRepo(t *testing.T) {
	ctx := context.Background()
	brandID := uuid.New()
	categoryID := uuid.New()
	productID := uuid.New()
	createdAt := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(time.Minute)
	input := product.Product{
		BrandID:             brandID,
		CategoryID:          categoryID,
		Name:                "Produto teste",
		Sku:                 "SKU-001",
		BarCodeEan:          "7891234567890",
		ShortDescription:    "Descricao curta",
		DetailedDescription: "Descricao detalhada",
		UnitOfMeasure:       "UN",
		CostPrice:           10.5,
		SalePrice:           15.75,
		PromotionalPrice:    12.99,
		GrossWeight:         1.2,
		NetWeight:           1.1,
		Height:              10,
		Width:               20,
		Length:              30,
		Status:              "ACTIVE",
	}

	t.Run("success", func(t *testing.T) {
		called := false
		db := &fakeDB{
			queryRowFn: func(gotCtx context.Context, query string, args ...any) pgx.Row {
				called = true
				assert.Equal(t, ctx, gotCtx)
				assert.Contains(t, query, "INSERT INTO products")
				assert.Equal(t, []any{
					brandID, categoryID, input.Name, input.Sku, input.BarCodeEan,
					input.ShortDescription, input.DetailedDescription, input.UnitOfMeasure,
					input.CostPrice, input.SalePrice, input.PromotionalPrice, input.GrossWeight,
					input.NetWeight, input.Height, input.Width, input.Length, input.Status,
				}, args)
				return fakeRow{id: productID, createdAt: createdAt, updatedAt: updatedAt}
			},
		}
		repo := NewProductRepo(nil, database.New(db))

		got, err := repo.CreateProductRepo(ctx, input)

		require.NoError(t, err)
		assert.True(t, called)
		assert.Equal(t, productID, got.ID)
		assert.Equal(t, createdAt, got.CreatedAt)
		assert.Equal(t, updatedAt, got.UpdatedAt)
		assert.Equal(t, input.Name, got.Name)
	})

	t.Run("returns scan error", func(t *testing.T) {
		scanErr := errors.New("scan failed")
		db := &fakeDB{
			queryRowFn: func(context.Context, string, ...any) pgx.Row {
				return fakeRow{err: scanErr}
			},
		}
		repo := NewProductRepo(nil, database.New(db))

		got, err := repo.CreateProductRepo(ctx, input)

		require.ErrorIs(t, err, scanErr)
		assert.Equal(t, product.Product{}, got)
	})
}

func TestProductRepoSoftDeleteProduct(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()

	tests := []struct {
		name        string
		commandTag  pgconn.CommandTag
		execErr     error
		expectedErr error
	}{
		{name: "success", commandTag: pgconn.NewCommandTag("UPDATE 1")},
		{name: "not found", commandTag: pgconn.NewCommandTag("UPDATE 0"), expectedErr: ErrProductNotFound},
		{name: "database error", execErr: errors.New("database unavailable")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			db := &fakeDB{
				execFn: func(gotCtx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
					called = true
					assert.Equal(t, ctx, gotCtx)
					assert.Contains(t, query, "UPDATE products SET deleted_at = now()")
					assert.Equal(t, []any{productID}, args)
					return tt.commandTag, tt.execErr
				},
			}
			repo := NewProductRepo(nil, database.New(db))

			err := repo.SoftDeleteProduct(ctx, productID)

			assert.True(t, called)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				return
			}
			if tt.execErr != nil {
				assert.ErrorIs(t, err, tt.execErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
