package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                  uuid.UUID
	BrandID             uuid.UUID
	CategoryID          uuid.UUID
	Name                string
	Sku                 string
	BarCodeEan          string
	ShortDescription    string
	DetailedDescription string
	UnitOfMeasure       string
	CostPrice           float64
	SalePrice           float64
	PromotionalPrice    float64
	GrossWeight         float64
	NetWeight           float64
	Height              float64
	Width               float64
	Length              float32
	Status              string
	CreatedAt           time.Time
	UpdatedAt            time.Time
}
