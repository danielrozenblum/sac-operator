package model

import (
	"github.com/google/uuid"
	"time"
)

type Application struct {
	ID         uuid.UUID
	Name       string
	CreatedOn  time.Time
	ModifiedOn time.Time
}
