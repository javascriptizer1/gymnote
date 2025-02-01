package entity

import (
	"time"

	"github.com/google/uuid"
)

type SetOption func(o *Set)

type Set struct {
	id         uuid.UUID
	userID     string
	exerciseID uuid.UUID
	setNumber  uint8
	weight     float32
	reps       uint8
	difficulty string
	notes      string
	createdAt  time.Time
}

func (s *Set) ID() uuid.UUID {
	return s.id
}

func (s *Set) UserID() string {
	return s.userID
}

func (s *Set) ExerciseID() uuid.UUID {
	return s.exerciseID
}

func (s *Set) Number() uint8 {
	return s.setNumber
}

func (s *Set) Weight() float32 {
	return s.weight
}

func (s *Set) Reps() uint8 {
	return s.reps
}

func (s *Set) Difficulty() string {
	return s.difficulty
}

func (s *Set) Notes() string {
	return s.notes
}

func (s *Set) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Set) SetWeight(weight float32) {
	s.weight = weight
}

func (s *Set) SetReps(reps uint8) {
	s.reps = reps
}

func NewSet(opts ...SetOption) *Set {
	set := &Set{}

	for _, opt := range opts {
		opt(set)
	}

	return set
}

type SetInitSpecification struct {
	UserID     string
	ExerciseID uuid.UUID
	Number     uint8
	Weight     float32
	Reps       uint8
	Difficulty string
	Notes      string
}

func WithSetInitSpec(s SetInitSpecification) SetOption {
	return func(o *Set) {
		o.id = uuid.New()
		o.userID = s.UserID
		o.exerciseID = s.ExerciseID
		o.setNumber = s.Number
		o.weight = s.Weight
		o.reps = s.Reps
		o.difficulty = s.Difficulty
		o.notes = s.Notes
		o.createdAt = time.Now()
	}
}

type SetRestoreSpecification struct {
	ID         uuid.UUID
	UserID     string
	ExerciseID uuid.UUID
	Number     uint8
	Weight     float32
	Reps       uint8
	Difficulty string
	Notes      string
	CreatedAt  time.Time
}

func WithSetRestoreSpec(s SetRestoreSpecification) SetOption {
	return func(o *Set) {
		o.id = s.ID
		o.userID = s.UserID
		o.exerciseID = s.ExerciseID
		o.setNumber = s.Number
		o.weight = s.Weight
		o.reps = s.Reps
		o.difficulty = s.Difficulty
		o.notes = s.Notes
		o.createdAt = s.CreatedAt
	}
}
