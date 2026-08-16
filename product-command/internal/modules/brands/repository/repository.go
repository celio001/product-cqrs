package brands_repository

import (
	"context"
	"errors"

	"github.com/celio001/product-command/internal/database"
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
	Tx     *database.Queries
}

type BrandsRepoInterface interface {
	GetBrandByID(ctx context.Context, id uuid.UUID) (brands.Brand, error)
	CreateBrand(ctx context.Context, b brands.Brand) (brands.Brand, error)
	SoftDeleteBrand(ctx context.Context, id uuid.UUID) error

	WithTx(tx pgx.Tx) BrandsRepoInterface
	BeginTx(context.Context) (pgx.Tx, error)
}

func NewBrandsRepository(pool *pgxpool.Pool, Tx *database.Queries) BrandsRepoInterface {
	return &brandsRepo{PgPool: pool, Tx: Tx}
}

func (r *brandsRepo) WithTx(tx pgx.Tx) BrandsRepoInterface {
	return &brandsRepo{
		PgPool: r.PgPool,
		Tx:     r.Tx.WithTx(tx),
	}
}

func (r *brandsRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.PgPool.Begin(ctx)
}

func (b *brandsRepo) GetBrandByID(ctx context.Context, id uuid.UUID) (brands.Brand, error) {
	var brand brands.Brand
	err := b.Tx.DB.QueryRow(ctx, "SELECT id, name FROM brands WHERE id = $1 AND deleted_at IS NULL", id).Scan(&brand.ID, &brand.Name)
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
	err := b.Tx.DB.QueryRow(ctx, "INSERT INTO brands (name) VALUES($1) RETURNING id, created_at", brand.Name).Scan(&brand.ID, &brand.CreatedAt)
	if err != nil {
		return brands.Brand{}, err
	}
	return brand, nil
}

func (b *brandsRepo) SoftDeleteBrand(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE brands SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL"

	result, err := b.Tx.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrBrandNotFound
	}

	return nil
}
