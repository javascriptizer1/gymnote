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
}

type SessionExercise struct {
	ID             uuid.UUID
	Exercise       Exercise
	ExerciseNumber uint8
	Sets           []Set
}

func (se *SessionExercise) SetCount() uint8 {
	return uint8(len(se.Sets))
}

func (se *SessionExercise) TotalVolume() float32 {
	totalVolume := float32(0)
	for _, set := range se.Sets {
		totalVolume += set.Weight * float32(set.Reps)
	}

	return totalVolume
}
