package entity

import (
	"slices"
	"time"

	"github.com/google/uuid"

	"gymnote/internal/errs"
)

type TrainingSessionOption func(o *TrainingSession)

type TrainingSession struct {
	id        uuid.UUID
	userID    string
	date      time.Time
	exercises []SessionExercise
	notes     string
	createdAt time.Time
}

func (ts *TrainingSession) ID() uuid.UUID {
	return ts.id
}

func (ts *TrainingSession) UserID() string {
	return ts.userID
}

func (ts *TrainingSession) Date() time.Time {
	return ts.date
}

func (ts *TrainingSession) TotalVolume() float32 {
	tv := float32(0)
	for _, exc := range ts.exercises {
		tv += exc.TotalVolume()
	}
	return tv
}

func (ts *TrainingSession) ExerciseCount() uint8 {
	return uint8(len(ts.exercises))
}

func (ts *TrainingSession) SetCount() uint8 {
	sets := uint8(0)
	for _, exs := range ts.exercises {
		sets += exs.SetCount()
	}

	return sets
}

func (ts *TrainingSession) Exercises() []SessionExercise {
	return ts.exercises
}

func (ts *TrainingSession) Notes() string {
	return ts.notes
}

func (ts *TrainingSession) ActiveExercise() *SessionExercise {
	if len(ts.exercises) == 0 {
		return nil
	}
	return &ts.exercises[len(ts.exercises)-1]
}

func (ts *TrainingSession) DeleteLastExercise(exerciseID uuid.UUID) error {
	for i := len(ts.exercises) - 1; i >= 0; i-- {
		if ts.exercises[i].Exercise.ID() == exerciseID {
			ts.exercises = slices.Delete(ts.exercises, i, i+1)
			return nil
		}
	}
	return errs.ErrExerciseNotFound
}

func (ts *TrainingSession) CreatedAt() time.Time {
	return ts.createdAt
}

func (ts *TrainingSession) AddExercise(exercise *SessionExercise) {
	ts.exercises = append(ts.exercises, *exercise)
}

func NewTrainingSession(opts ...TrainingSessionOption) *TrainingSession {
	session := &TrainingSession{}

	for _, opt := range opts {
		opt(session)
	}

	return session
}

type TrainingSessionInitSpecification struct {
	UserID    string
	Date      time.Time
	Exercises []SessionExercise
	Notes     string
}

func WithTrainingSessionInitSpec(spec TrainingSessionInitSpecification) TrainingSessionOption {
	return func(ts *TrainingSession) {
		copiedExercises := make([]SessionExercise, len(spec.Exercises))
		copy(copiedExercises, spec.Exercises)

		ts.id = uuid.New()
		ts.userID = spec.UserID
		ts.date = spec.Date
		ts.exercises = copiedExercises
		ts.notes = spec.Notes
		ts.createdAt = time.Now()
	}
}

type TrainingSessionRestoreSpecification struct {
	ID        uuid.UUID
	UserID    string
	Date      time.Time
	Exercises []SessionExercise
	Notes     string
	CreatedAt time.Time
}

func WithTrainingSessionRestoreSpec(spec TrainingSessionRestoreSpecification) TrainingSessionOption {
	return func(ts *TrainingSession) {
		copiedExercises := make([]SessionExercise, len(spec.Exercises))
		copy(copiedExercises, spec.Exercises)

		ts.id = spec.ID
		ts.userID = spec.UserID
		ts.date = spec.Date
		ts.exercises = copiedExercises
		ts.notes = spec.Notes
		ts.createdAt = spec.CreatedAt
	}
}
