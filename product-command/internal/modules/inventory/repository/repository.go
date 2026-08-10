package inventory_repository

import (
	"context"

	"github.com/celio001/product-command/internal/database"
	"github.com/celio001/product-command/internal/modules/inventory"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type inventoryRepo struct {
	PgPool *pgxpool.Pool
	Tx     *database.Queries
}

type InventoryRepoInterface interface {
	CreateInventoryRepo(ctx context.Context, i inventory.Inventory) (inventory.Inventory, error)

	//init transaction
	WithTx(tx pgx.Tx) InventoryRepoInterface
	//close transection
	BeginTx(context.Context) (pgx.Tx, error)
}

func NewInventoryRepo(Pg *pgxpool.Pool, Tx *database.Queries) InventoryRepoInterface {
	return &inventoryRepo{
		PgPool: Pg,
		Tx:     Tx,
	}
}

func (r *inventoryRepo) WithTx(tx pgx.Tx) InventoryRepoInterface {
	return &inventoryRepo{
		PgPool: r.PgPool,
		Tx:     r.Tx.WithTx(tx),
	}
}

func (r *inventoryRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.PgPool.Begin(ctx)
}

func (r *inventoryRepo) CreateInventoryRepo(ctx context.Context, i inventory.Inventory) (inventory.Inventory, error) {
	query := `INSERT INTO product_inventory 
	(product_id, location_aisle, quantity_available, minimum_stock, maximum_stock) 
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, updated_at`
	err := r.Tx.DB.QueryRow(ctx, query, i.ProductID, i.LocationAisle, i.QuantityAvailable, i.MinimumStock, i.MaximumStock).Scan(&i.ID, &i.UpdatedAt)
	if err != nil {
		return inventory.Inventory{}, err
	}

	return i, nil
}
