package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExerciseOption func(o *Exercise)

type Exercise struct {
	id          uuid.UUID
	createdAt   time.Time
	name        string
	muscleGroup string
	equipment   string
}

func (e *Exercise) ID() uuid.UUID {
	return e.id
}

func (e *Exercise) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Exercise) Name() string {
	return e.name
}

func (e *Exercise) MuscleGroup() string {
	return e.muscleGroup
}

func (e *Exercise) Equipment() string {
	return e.equipment
}

func NewExercise(opts ...ExerciseOption) *Exercise {
	exercise := &Exercise{}

	for _, opt := range opts {
		opt(exercise)
	}

	return exercise
}

type ExerciseInitSpecification struct {
	Name        string
	MuscleGroup string
	Equipment   string
}

func WithExerciseInitSpec(e ExerciseInitSpecification) ExerciseOption {
	return func(o *Exercise) {
		o.id = uuid.New()
		o.createdAt = time.Now()
		o.name = e.Name
		o.muscleGroup = e.MuscleGroup
		o.equipment = e.Equipment
	}
}

type ExerciseRestoreSpecification struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Name        string
	MuscleGroup string
	Equipment   string
}

func WithExerciseRestoreSpec(e ExerciseRestoreSpecification) ExerciseOption {
	return func(o *Exercise) {
		o.id = e.ID
		o.createdAt = e.CreatedAt
		o.name = e.Name
		o.muscleGroup = e.MuscleGroup
		o.equipment = e.Equipment
	}
}
