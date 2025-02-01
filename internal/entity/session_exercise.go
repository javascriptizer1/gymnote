package entity

import "github.com/google/uuid"

type SessionExerciseOption func(o *SessionExercise)

type SessionExercise struct {
	*Exercise
	id     uuid.UUID
	number uint8
	sets   []Set
}

func (se *SessionExercise) ID() uuid.UUID {
	return se.id
}

func (se *SessionExercise) Number() uint8 {
	return se.number
}

func (se *SessionExercise) SetCount() uint8 {
	return uint8(len(se.sets))
}

func (se *SessionExercise) Sets() []Set {
	return se.sets
}

func (se *SessionExercise) ActiveSet() *Set {
	if len(se.sets) == 0 {
		return nil
	}
	return &se.sets[len(se.sets)-1]
}

func (se *SessionExercise) TotalVolume() float32 {
	totalVolume := float32(0)
	for _, set := range se.sets {
		totalVolume += set.Weight() * float32(set.Reps())
	}

	return totalVolume
}

func NewSessionExercise(exercise *Exercise, sets []Set, opts ...SessionExerciseOption) *SessionExercise {
	sessionExercise := &SessionExercise{
		Exercise: exercise,
		sets:     sets,
	}

	for _, opt := range opts {
		opt(sessionExercise)
	}

	return sessionExercise
}

type SessionExerciseInitSpecification struct {
	Number uint8
}

func WithSessionExerciseInitSpec(s SessionExerciseInitSpecification) SessionExerciseOption {
	return func(o *SessionExercise) {
		o.id = uuid.New()
		o.number = s.Number
	}
}

type SessionExerciseRestoreSpecification struct {
	ID     uuid.UUID
	Number uint8
}

func WithSessionExerciseRestoreSpec(s SessionExerciseRestoreSpecification) SessionExerciseOption {
	return func(o *SessionExercise) {
		o.id = s.ID
		o.number = s.Number
	}
}
