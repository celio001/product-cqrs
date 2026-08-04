package categories_repository

import (
	"context"

	"github.com/celio001/product-command/internal/modules/categories"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoriesRepo struct {
	PgPool *pgxpool.Pool
}

type CategoriesInterface interface {
	CreateCategorie(ctx context.Context, categorie categories.Categories)(categories.Categories, error)
}

func NewCategoriesRepo(PgPool *pgxpool.Pool) CategoriesInterface{
	return &categoriesRepo{PgPool: PgPool}
}

func (c *categoriesRepo) CreateCategorie(ctx context.Context, categorie categories.Categories)(categories.Categories, error){
	err := c.PgPool.QueryRow(ctx, "INSERT INTO categories(name) VALUES($1) RETURNING id, parent_id", categorie.Name).Scan(&categorie.ID, &categorie.ParentID)
	if err != nil {
		return categories.Categories{}, err
	}

	return categorie, nil
}