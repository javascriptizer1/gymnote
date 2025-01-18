package entity

import (
	"time"

	"github.com/google/uuid"
)

type Exercise struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Name        string
	MuscleGroup string
	Equipment   string
	Sets        []Set
}
