package product_repository

import (
	"context"

	"github.com/celio001/product-command/internal/modules/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	PgPool *pgxpool.Pool
}

type ProductRepoInterface interface {
	CreateProductRepo(ctx context.Context, p product.Product) (product.Product, error)
}

func NewProductRepo(PgPool *pgxpool.Pool) ProductRepoInterface {
	return &productRepo{PgPool: PgPool}
}

func (r productRepo) CreateProductRepo(ctx context.Context, p product.Product) (product.Product, error) {
	query := `INSERT INTO products 
	(brand_id, category_id, name, sku, barcode_ean13, short_description, detailed_description, unit_of_measure, cost_price, sale_price, promotional_price, gross_weight, net_weight, height, width, length, status) 
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) 
	RETURNING id"`

	err := r.PgPool.QueryRow(ctx, query, p.BrandID, p.CategoryID, p.Name, p.Sku, p.BarCodeEan, p.ShortDescription, p.DetailedDescription, p.UnitOfMeasure, p.CostPrice, p.SalePrice, p.PromotionalPrice, p.GrossWeifht, p.NetWeight, p.Height, p.Width, p.Length, p.Status).Scan(&p.ID)
	if err != nil {
		return product.Product{}, err
	}

	return p, nil 
}