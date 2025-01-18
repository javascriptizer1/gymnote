package entity

import (
	"github.com/google/uuid"
)

type Set struct {
	ID         uuid.UUID
	UserID     string
	ExerciseID uuid.UUID
	SetNumber  uint8
	Weight     float32
	Reps       uint8
	Difficulty string
	Notes      string
}
