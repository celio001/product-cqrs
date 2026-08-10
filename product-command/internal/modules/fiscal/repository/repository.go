package fiscal_repository

import (
	"context"

	"github.com/celio001/product-command/internal/database"
	"github.com/celio001/product-command/internal/modules/fiscal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type fiscalRepository struct {
	PgPool *pgxpool.Pool
	Tx     *database.Queries
}

type FiscalRepositoryInterface interface {
	CreateFiscalData(ctx context.Context, f fiscal.FiscalData) (fiscal.FiscalData, error)

	//init transaction
	WithTx(tx pgx.Tx) FiscalRepositoryInterface
	//close transection
	BeginTx(context.Context) (pgx.Tx, error)
}

func NewFiscalRepo(PgPool *pgxpool.Pool, Tx *database.Queries) FiscalRepositoryInterface {
	return &fiscalRepository{
		PgPool: PgPool,
		Tx:     Tx,
	}
}

func (r *fiscalRepository) WithTx(tx pgx.Tx) FiscalRepositoryInterface {
	return &fiscalRepository{
		PgPool: r.PgPool,
		Tx:     r.Tx.WithTx(tx),
	}
}

func (r *fiscalRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.PgPool.Begin(ctx)
}

func (r *fiscalRepository) CreateFiscalData(ctx context.Context, f fiscal.FiscalData) (fiscal.FiscalData, error) {
	query := `INSERT INTO product_fiscal_data 
	(product_id, ncm_code, cest_code, origin_code, icms_rate, pis_rate, cofins_rate, ipi_rate)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING id, updated_at`

	err := r.Tx.DB.QueryRow(ctx, query, f.ProductId, f.NcmCode, f.CestCode, f.OriginCode, f.IcmsRate, f.PisRate, f.CofinsRate, f.IpiRate).Scan(&f.ID, &f.UpdatedAt)
	if err != nil {
		return fiscal.FiscalData{}, err
	}

	return f, nil
}
