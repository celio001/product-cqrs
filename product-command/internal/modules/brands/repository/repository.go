package brands_repository

import (
	"context"
	"errors"

	"github.com/celio001/product-command/internal/modules/brands"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrBrandNotFound = errors.New("brand not found")
)

type brandsRepo struct {
	PgPool *pgxpool.Pool
}

type BrandsRepoInterface interface {
	GetBrandByID(ctx context.Context, id uuid.UUID) (brands.Brand, error)
	CreateBrand(ctx context.Context, b brands.Brand) (brands.Brand, error)
	SoftDeleteBrand(ctx context.Context, id uuid.UUID) error
}

func NewBrandsRepository(pool *pgxpool.Pool) BrandsRepoInterface {
	return &brandsRepo{PgPool: pool}
}

func (b *brandsRepo) GetBrandByID(ctx context.Context, id uuid.UUID) (brands.Brand, error) {
	var brand brands.Brand
	err := b.PgPool.QueryRow(ctx, "SELECT id, name FROM brands WHERE id = $1 AND deleted_at IS NULL", id).Scan(&brand.ID, &brand.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return brands.Brand{}, ErrBrandNotFound
		} else {
			return brands.Brand{}, err
		}
	}
	return brand, nil
}

func (b *brandsRepo) CreateBrand(ctx context.Context, brand brands.Brand) (brands.Brand, error) {
	err := b.PgPool.QueryRow(ctx, "INSERT INTO brands (name) VALUES($1) RETURNING id", brand.Name).Scan(&brand.ID)
	if err != nil {
		return brands.Brand{}, err
	}
	return brand, nil
}

func (b *brandsRepo) SoftDeleteBrand(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE brands SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL"

	result, err := b.PgPool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrBrandNotFound
	}

	return nil
}
