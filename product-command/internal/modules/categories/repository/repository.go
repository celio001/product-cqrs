package categories_repository

import (
	"context"
	"errors"

	"github.com/celio001/product-command/internal/modules/categories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type categoriesRepo struct {
	PgPool *pgxpool.Pool
}

type CategoriesInterface interface {
	CreateCategory(ctx context.Context, categorie categories.Categories) (categories.Categories, error)
	SoftDeleteCategory(ctx context.Context, id uuid.UUID) error
}

func NewCategoriesRepo(PgPool *pgxpool.Pool) CategoriesInterface {
	return &categoriesRepo{PgPool: PgPool}
}

func (c *categoriesRepo) CreateCategory(ctx context.Context, categorie categories.Categories) (categories.Categories, error) {
	err := c.PgPool.QueryRow(ctx, "INSERT INTO categories(name) VALUES($1) RETURNING id, parent_id", categorie.Name).Scan(&categorie.ID, &categorie.ParentID)
	if err != nil {
		return categories.Categories{}, err
	}

	return categorie, nil
}

func (c *categoriesRepo) SoftDeleteCategory(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE categories SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL"
	result, err := c.PgPool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}
	return nil
}
