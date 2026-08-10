package product_dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	BrandID             uuid.UUID        `json:"brand_id,omitempty" validate:"omitempty,uuid"`
	CategoryID          uuid.UUID        `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Name                string           `json:"name" validate:"required,min=1,max=255"`
	Sku                 string           `json:"sku" validate:"required,max=50"`
	BarCodeEan          string           `json:"barcode_ean13,omitempty" validate:"omitempty,max=13"`
	ShortDescription    string           `json:"short_description,omitempty" validate:"omitempty,max=255"`
	DetailedDescription string           `json:"detailed_description,omitempty"`
	UnitOfMeasure       string           `json:"unit_of_measure" validate:"required,max=10"`
	CostPrice           float64          `json:"cost_price" validate:"required,min=0"`
	SalePrice           float64          `json:"sale_price" validate:"required,min=0"`
	PromotionalPrice    float64          `json:"promotional_price,omitempty" validate:"omitempty,min=0"`
	GrossWeight         float64          `json:"gross_weight,omitempty" validate:"omitempty,min=0"`
	NetWeight           float64          `json:"net_weight,omitempty" validate:"omitempty,min=0"`
	Height              float64          `json:"height,omitempty" validate:"omitempty,min=0"`
	Width               float64          `json:"width,omitempty" validate:"omitempty,min=0"`
	Length              float32          `json:"length,omitempty" validate:"omitempty,min=0"`
	Status              string           `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE INACTIVE"`
	Stock               InventoryProductRequest `json:"stock" validate:"required"`
	Fiscal              FiscalProductRequest    `json:"fiscal" validate:"required"`
}

type InventoryProductRequest struct {
	LocationAisle     string  `json:"location_aisle,omitempty" validate:"omitempty,max=50"`
	QuantityAvailable float64 `json:"quantity_available" validate:"min=0"`
	MinimumStock      float64 `json:"minimum_stock" validate:"min=0"`
	MaximumStock      float64 `json:"maximum_stock,omitempty" validate:"omitempty,min=0"`
}

type FiscalProductRequest struct {
	NcmCode    string  `json:"ncm_code,omitempty" validate:"omitempty,max=8"`
	CestCode   string  `json:"cest_code,omitempty" validate:"omitempty,max=7"`
	OriginCode int     `json:"origin_code,omitempty" validate:"omitempty,min=0"`
	IcmsRate   float64 `json:"icms_rate,omitempty" validate:"omitempty,min=0,max=100"`
	PisRate    float64 `json:"pis_rate,omitempty" validate:"omitempty,min=0,max=100"`
	CofinsRate float64 `json:"cofins_rate,omitempty" validate:"omitempty,min=0,max=100"`
	IpiRate    float64 `json:"ipi_rate,omitempty" validate:"omitempty,min=0,max=100"`
}

type CreateProductResponse struct {
	ID                  uuid.UUID                `json:"id"`
	BrandID             uuid.UUID                `json:"brand_id,omitempty"`
	CategoryID          uuid.UUID                `json:"category_id,omitempty"`
	Name                string                   `json:"name"`
	Sku                 string                   `json:"sku"`
	BarCodeEan          string                   `json:"barcode_ean13,omitempty"`
	ShortDescription    string                   `json:"short_description,omitempty"`
	DetailedDescription string                   `json:"detailed_description,omitempty"`
	UnitOfMeasure       string                   `json:"unit_of_measure"`
	CostPrice           float64                  `json:"cost_price"`
	SalePrice           float64                  `json:"sale_price"`
	PromotionalPrice    float64                  `json:"promotional_price,omitempty"`
	GrossWeight         float64                  `json:"gross_weight,omitempty"`
	NetWeight           float64                  `json:"net_weight,omitempty"`
	Height              float64                  `json:"height,omitempty"`
	Width               float64                  `json:"width,omitempty"`
	Length              float32                  `json:"length,omitempty"`
	Status              string                   `json:"status"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	Stock               InventoryProductResponse `json:"stock"`
	Fiscal              FiscalProductResponse    `json:"fiscal"`
}

type InventoryProductResponse struct {
	ID                uuid.UUID `json:"id"`
	ProductID         uuid.UUID `json:"product_id"`
	LocationAisle     string    `json:"location_aisle,omitempty"`
	QuantityAvailable float64   `json:"quantity_available"`
	MinimumStock      float64   `json:"minimum_stock"`
	MaximumStock      float64   `json:"maximum_stock,omitempty"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type FiscalProductResponse struct {
	ID         uuid.UUID `json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	NcmCode    string    `json:"ncm_code,omitempty"`
	CestCode   string    `json:"cest_code,omitempty"`
	OriginCode int     `json:"origin_code,omitempty"`
	IcmsRate   float64   `json:"icms_rate"`
	PisRate    float64   `json:"pis_rate"`
	CofinsRate float64   `json:"cofins_rate"`
	IpiRate    float64   `json:"ipi_rate"`
	UpdatedAt  time.Time `json:"updated_at"`
}
