package categories

import (
	"time"

	"github.com/google/uuid"
)

type Categories struct {
	ID        uuid.UUID
	ParentID uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time
	DeletedAt time.Time
}
