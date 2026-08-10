package product_repository

import (
	"context"

	"github.com/celio001/product-command/internal/database"
	"github.com/celio001/product-command/internal/modules/product"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepo struct {
	PgPool *pgxpool.Pool
	Tx     *database.Queries
}

type ProductRepoInterface interface {
	CreateProductRepo(ctx context.Context, p product.Product) (product.Product, error)

	//init transaction
	WithTx(tx pgx.Tx) ProductRepoInterface
	//close transection
	BeginTx(context.Context) (pgx.Tx, error)
}

func NewProductRepo(PgPool *pgxpool.Pool, tx *database.Queries) ProductRepoInterface {
	return &productRepo{PgPool: PgPool, Tx: tx}
}

func (r *productRepo) WithTx(tx pgx.Tx) ProductRepoInterface {
	return &productRepo{
		PgPool: r.PgPool,
		Tx:     r.Tx.WithTx(tx),
	}
}

func (r *productRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.PgPool.Begin(ctx)
}

func (r *productRepo) CreateProductRepo(ctx context.Context, p product.Product) (product.Product, error) {

	query := `INSERT INTO products 
	(brand_id, category_id, name, sku, barcode_ean13, short_description, detailed_description, unit_of_measure, cost_price, sale_price, promotional_price, gross_weight, net_weight, height, width, length, status) 
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) 
	RETURNING id, created_at, updated_at`

	err := r.Tx.DB.QueryRow(ctx, query, p.BrandID, p.CategoryID, p.Name, p.Sku, p.BarCodeEan, p.ShortDescription, p.DetailedDescription, p.UnitOfMeasure, p.CostPrice, p.SalePrice, p.PromotionalPrice, p.GrossWeight, p.NetWeight, p.Height, p.Width, p.Length, p.Status).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return product.Product{}, err
	}
	return p, nil
}
