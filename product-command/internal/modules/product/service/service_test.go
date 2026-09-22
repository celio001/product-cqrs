package product_service

import (
	"context"
	"errors"
	"testing"
	"time"

	product_dto "github.com/celio001/product-command/internal/fiber/v1/product/dto"
	"github.com/celio001/product-command/internal/modules/brands"
	brands_repository "github.com/celio001/product-command/internal/modules/brands/repository"
	"github.com/celio001/product-command/internal/modules/categories"
	categories_repository "github.com/celio001/product-command/internal/modules/categories/repository"
	"github.com/celio001/product-command/internal/modules/fiscal"
	fiscal_repository "github.com/celio001/product-command/internal/modules/fiscal/repository"
	"github.com/celio001/product-command/internal/modules/inventory"
	inventory_repository "github.com/celio001/product-command/internal/modules/inventory/repository"
	"github.com/celio001/product-command/internal/modules/product"
	product_publisher "github.com/celio001/product-command/internal/modules/product/publisher"
	product_repository "github.com/celio001/product-command/internal/modules/product/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

type serviceTx struct {
	pgx.Tx
	commitErr   error
	rollbackErr error
	commits     int
	rollbacks   int
}

func (t *serviceTx) Begin(context.Context) (pgx.Tx, error) { return t, nil }
func (t *serviceTx) Commit(context.Context) error {
	t.commits++
	return t.commitErr
}
func (t *serviceTx) Rollback(context.Context) error {
	t.rollbacks++
	return t.rollbackErr
}
func (t *serviceTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
func (t *serviceTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults { return nil }
func (t *serviceTx) LargeObjects() pgx.LargeObjects                         { return pgx.LargeObjects{} }
func (t *serviceTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (t *serviceTx) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (t *serviceTx) Query(context.Context, string, ...any) (pgx.Rows, error) { return nil, nil }
func (t *serviceTx) QueryRow(context.Context, string, ...any) pgx.Row        { return nil }
func (t *serviceTx) Conn() *pgx.Conn                                         { return nil }

type serviceProductRepo struct {
	tx                  pgx.Tx
	beginErr            error
	createProductFn     func(context.Context, product.Product) (product.Product, error)
	softDeleteProductFn func(context.Context, uuid.UUID) error
}

func (r *serviceProductRepo) BeginTx(context.Context) (pgx.Tx, error) { return r.tx, r.beginErr }
func (r *serviceProductRepo) WithTx(pgx.Tx) product_repository.ProductRepoInterface {
	return r
}
func (r *serviceProductRepo) CreateProductRepo(ctx context.Context, p product.Product) (product.Product, error) {
	if r.createProductFn != nil {
		return r.createProductFn(ctx, p)
	}
	return p, nil
}
func (r *serviceProductRepo) SoftDeleteProduct(ctx context.Context, id uuid.UUID) error {
	if r.softDeleteProductFn != nil {
		return r.softDeleteProductFn(ctx, id)
	}
	return nil
}

type serviceInventoryRepo struct {
	createFn func(context.Context, inventory.Inventory) (inventory.Inventory, error)
}

func (r *serviceInventoryRepo) WithTx(pgx.Tx) inventory_repository.InventoryRepoInterface { return r }
func (r *serviceInventoryRepo) BeginTx(context.Context) (pgx.Tx, error)                   { return nil, nil }
func (r *serviceInventoryRepo) CreateInventoryRepo(ctx context.Context, i inventory.Inventory) (inventory.Inventory, error) {
	if r.createFn != nil {
		return r.createFn(ctx, i)
	}
	return i, nil
}

type serviceFiscalRepo struct {
	createFn func(context.Context, fiscal.FiscalData) (fiscal.FiscalData, error)
}

func (r *serviceFiscalRepo) WithTx(pgx.Tx) fiscal_repository.FiscalRepositoryInterface { return r }
func (r *serviceFiscalRepo) BeginTx(context.Context) (pgx.Tx, error)                   { return nil, nil }
func (r *serviceFiscalRepo) CreateFiscalData(ctx context.Context, f fiscal.FiscalData) (fiscal.FiscalData, error) {
	if r.createFn != nil {
		return r.createFn(ctx, f)
	}
	return f, nil
}

type serviceCategoryRepo struct {
	category categories.Categories
	err      error
}

func (r *serviceCategoryRepo) GetCategoryByID(context.Context, uuid.UUID) (categories.Categories, error) {
	return r.category, r.err
}
func (r *serviceCategoryRepo) CreateCategory(context.Context, categories.Categories) (categories.Categories, error) {
	return categories.Categories{}, nil
}
func (r *serviceCategoryRepo) SoftDeleteCategory(context.Context, uuid.UUID) error     { return nil }
func (r *serviceCategoryRepo) WithTx(pgx.Tx) categories_repository.CategoriesInterface { return r }
func (r *serviceCategoryRepo) BeginTx(context.Context) (pgx.Tx, error)                 { return nil, nil }

type serviceBrandRepo struct {
	brand brands.Brand
	err   error
}

func (r *serviceBrandRepo) GetBrandByID(context.Context, uuid.UUID) (brands.Brand, error) {
	return r.brand, r.err
}
func (r *serviceBrandRepo) CreateBrand(context.Context, brands.Brand) (brands.Brand, error) {
	return brands.Brand{}, nil
}
func (r *serviceBrandRepo) SoftDeleteBrand(context.Context, uuid.UUID) error    { return nil }
func (r *serviceBrandRepo) WithTx(pgx.Tx) brands_repository.BrandsRepoInterface { return r }
func (r *serviceBrandRepo) BeginTx(context.Context) (pgx.Tx, error)             { return nil, nil }

type servicePublisher struct {
	createdFn func(context.Context, product_dto.CreateProductResponse) error
	deletedFn func(context.Context, uuid.UUID) error
}

func (p *servicePublisher) PublishProductCreated(ctx context.Context, response product_dto.CreateProductResponse) error {
	if p.createdFn != nil {
		return p.createdFn(ctx, response)
	}
	return nil
}
func (p *servicePublisher) PublishProductDeleted(ctx context.Context, id uuid.UUID) error {
	if p.deletedFn != nil {
		return p.deletedFn(ctx, id)
	}
	return nil
}

var (
	_ product_repository.ProductRepoInterface     = (*serviceProductRepo)(nil)
	_ inventory_repository.InventoryRepoInterface = (*serviceInventoryRepo)(nil)
	_ fiscal_repository.FiscalRepositoryInterface = (*serviceFiscalRepo)(nil)
	_ categories_repository.CategoriesInterface   = (*serviceCategoryRepo)(nil)
	_ brands_repository.BrandsRepoInterface       = (*serviceBrandRepo)(nil)
	_ product_publisher.ProductPublisherInterface = (*servicePublisher)(nil)
)

func TestCreateProductSvc(t *testing.T) {
	ctx := context.Background()
	brandID := uuid.New()
	categoryID := uuid.New()
	productID := uuid.New()
	inventoryID := uuid.New()
	fiscalID := uuid.New()
	timestamp := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	tx := &serviceTx{}
	inputProduct := product.Product{BrandID: brandID, CategoryID: categoryID, Name: "Produto", Sku: "SKU-1"}
	inputInventory := inventory.Inventory{LocationAisle: "A1", QuantityAvailable: 10}
	inputFiscal := fiscal.FiscalData{NcmCode: "12345678", OriginCode: 0}
	createdProduct := inputProduct
	createdProduct.ID = productID
	createdProduct.CreatedAt = timestamp
	createdProduct.UpdatedAt = timestamp

	var receivedInventory inventory.Inventory
	var receivedFiscal fiscal.FiscalData
	var published product_dto.CreateProductResponse
	svc := NewProductSvc(
		&serviceProductRepo{
			tx: tx,
			createProductFn: func(_ context.Context, p product.Product) (product.Product, error) {
				assert.Equal(t, categoryID, p.CategoryID)
				assert.Equal(t, brandID, p.BrandID)
				return createdProduct, nil
			},
		},
		&serviceFiscalRepo{createFn: func(_ context.Context, f fiscal.FiscalData) (fiscal.FiscalData, error) {
			receivedFiscal = f
			f.ID = fiscalID
			f.UpdatedAt = timestamp
			return f, nil
		}},
		&serviceInventoryRepo{createFn: func(_ context.Context, i inventory.Inventory) (inventory.Inventory, error) {
			receivedInventory = i
			i.ID = inventoryID
			i.UpdatedAt = timestamp
			return i, nil
		}},
		&serviceCategoryRepo{category: categories.Categories{ID: categoryID}},
		&serviceBrandRepo{brand: brands.Brand{ID: brandID}},
		&servicePublisher{createdFn: func(_ context.Context, response product_dto.CreateProductResponse) error {
			published = response
			return nil
		}},
		noop.NewTracerProvider().Tracer("test"),
	)

	response, err := svc.CreateProductSvc(ctx, inputProduct, inputInventory, inputFiscal)

	require.NoError(t, err)
	assert.Equal(t, productID, response.ID)
	assert.Equal(t, inventoryID, response.Stock.ID)
	assert.Equal(t, fiscalID, response.Fiscal.ID)
	assert.Equal(t, productID, receivedInventory.ProductID)
	assert.Equal(t, productID, receivedFiscal.ProductId)
	assert.Equal(t, response, published)
	assert.Equal(t, 1, tx.commits)
	assert.Equal(t, 0, tx.rollbacks)
}

func TestCreateProductSvcReturnsCategoryError(t *testing.T) {
	categoryErr := errors.New("category unavailable")
	tx := &serviceTx{}
	svc := NewProductSvc(
		&serviceProductRepo{tx: tx},
		&serviceFiscalRepo{},
		&serviceInventoryRepo{},
		&serviceCategoryRepo{err: categoryErr},
		&serviceBrandRepo{},
		&servicePublisher{},
		noop.NewTracerProvider().Tracer("test"),
	)

	_, err := svc.CreateProductSvc(context.Background(), product.Product{CategoryID: uuid.New()}, inventory.Inventory{}, fiscal.FiscalData{})

	assert.ErrorIs(t, err, categoryErr)
	assert.Equal(t, 0, tx.commits)
	assert.Equal(t, 0, tx.rollbacks)
}

func TestCreateProductSvcRollsBackWhenProductCreationFails(t *testing.T) {
	productErr := errors.New("product creation failed")
	tx := &serviceTx{}
	svc := NewProductSvc(
		&serviceProductRepo{tx: tx, createProductFn: func(context.Context, product.Product) (product.Product, error) {
			return product.Product{}, productErr
		}},
		&serviceFiscalRepo{},
		&serviceInventoryRepo{},
		&serviceCategoryRepo{category: categories.Categories{ID: uuid.New()}},
		&serviceBrandRepo{brand: brands.Brand{ID: uuid.New()}},
		&servicePublisher{},
		noop.NewTracerProvider().Tracer("test"),
	)

	_, err := svc.CreateProductSvc(context.Background(), product.Product{}, inventory.Inventory{}, fiscal.FiscalData{})

	assert.ErrorIs(t, err, productErr)
	assert.Equal(t, 0, tx.commits)
	assert.Equal(t, 1, tx.rollbacks)
}

func TestSoftDeleteProductSvc(t *testing.T) {
	ctx := context.Background()
	productID := uuid.New()
	tx := &serviceTx{}
	var publishedID uuid.UUID
	svc := NewProductSvc(
		&serviceProductRepo{tx: tx},
		&serviceFiscalRepo{},
		&serviceInventoryRepo{},
		&serviceCategoryRepo{},
		&serviceBrandRepo{},
		&servicePublisher{deletedFn: func(_ context.Context, id uuid.UUID) error {
			publishedID = id
			return nil
		}},
		noop.NewTracerProvider().Tracer("test"),
	)

	err := svc.SoftDeleteProductSvc(ctx, productID)

	require.NoError(t, err)
	assert.Equal(t, productID, publishedID)
	assert.Equal(t, 1, tx.commits)
	assert.Equal(t, 0, tx.rollbacks)
}

func TestSoftDeleteProductSvcRollsBackWhenPublishingFails(t *testing.T) {
	publishErr := errors.New("publish failed")
	tx := &serviceTx{}
	svc := NewProductSvc(
		&serviceProductRepo{tx: tx},
		&serviceFiscalRepo{},
		&serviceInventoryRepo{},
		&serviceCategoryRepo{},
		&serviceBrandRepo{},
		&servicePublisher{deletedFn: func(context.Context, uuid.UUID) error { return publishErr }},
		noop.NewTracerProvider().Tracer("test"),
	)

	err := svc.SoftDeleteProductSvc(context.Background(), uuid.New())

	assert.ErrorIs(t, err, publishErr)
	assert.Equal(t, 0, tx.commits)
	assert.Equal(t, 1, tx.rollbacks)
}
