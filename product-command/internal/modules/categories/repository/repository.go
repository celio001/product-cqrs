package categories_repository

import (
	"context"
	"errors"

	"github.com/celio001/product-command/internal/modules/categories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type categoriesRepo struct {
	PgPool *pgxpool.Pool
}

type CategoriesInterface interface {
	GetCategoryByID(ctx context.Context, id uuid.UUID) (categories.Categories, error)
	CreateCategory(ctx context.Context, categorie categories.Categories) (categories.Categories, error)
	SoftDeleteCategory(ctx context.Context, id uuid.UUID) error
}

func NewCategoriesRepo(PgPool *pgxpool.Pool) CategoriesInterface {
	return &categoriesRepo{PgPool: PgPool}
}

func (c *categoriesRepo) GetCategoryByID(ctx context.Context, id uuid.UUID) (categories.Categories, error) {

	var category categories.Categories

	query := "SELECT id, name FROM categories WHERE id = $1 AND deleted_at IS NULL"
	err := c.PgPool.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return categories.Categories{}, ErrCategoryNotFound
		} else {
			return categories.Categories{}, err
		}
	}

	return category, nil
}

func (c *categoriesRepo) CreateCategory(ctx context.Context, categorie categories.Categories) (categories.Categories, error) {
	query := "INSERT INTO categories(name) VALUES($1) RETURNING id, parent_id"
	err := c.PgPool.QueryRow(ctx, query, categorie.Name).Scan(&categorie.ID, &categorie.ParentID)
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
