package fiscal

import (
	"time"

	"github.com/google/uuid"
)

type FiscalData struct {
	ID         uuid.UUID
	ProductId  uuid.UUID
	NcmCode    string
	CestCode   string
	OriginCode int
	IcmsRate   float64
	PisRate    float64
	CofinsRate float64
	IpiRate    float64
	UpdatedAt  time.Time
}
