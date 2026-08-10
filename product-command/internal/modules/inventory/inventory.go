package inventory

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID                uuid.UUID
	ProductID         uuid.UUID
	LocationAisle     string
	QuantityAvailable float64
	MinimumStock      float64
	MaximumStock      float64
	CreatedAt         time.Time
	UpdatedAt          time.Time
}
