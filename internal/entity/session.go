package entity

import (
	"time"

	"github.com/google/uuid"
)

type TrainingSession struct {
	ID            uuid.UUID
	UserID        string
	Date          time.Time
	TotalVolume   float32
	ExerciseCount uint8
	Notes         string
	Exercises     []Exercise
}
