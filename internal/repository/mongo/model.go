package mongodb

import (
	"time"

	"github.com/google/uuid"

	"gymnote/internal/entity"
)

type ExerciseOption func(o *ExerciseRow)
type TrainingSessionOption func(o *TrainingSessionRow)
type SetOption func(o *SetRow)

type ExerciseRow struct {
	ID          string    `bson:"id"`
	CreatedAt   time.Time `bson:"created_at"`
	Name        string    `bson:"name"`
	MuscleGroup string    `bson:"muscle_group"`
	Equipment   string    `bson:"equipment"`
}

func (e *ExerciseRow) ToEntity() *entity.Exercise {
	id, _ := uuid.Parse(e.ID)
	return entity.NewExercise(
		entity.WithExerciseRestoreSpec(entity.ExerciseRestoreSpecification{
			ID:          id,
			Name:        e.Name,
			MuscleGroup: e.MuscleGroup,
			Equipment:   e.Equipment,
			CreatedAt:   e.CreatedAt,
		}),
	)
}

func NewExerciseRow(opts ...ExerciseOption) *ExerciseRow {
	exercise := &ExerciseRow{}

	for _, opt := range opts {
		opt(exercise)
	}

	return exercise
}

type ExerciseRowRestoreSpecification struct {
	ID          string
	CreatedAt   time.Time
	Name        string
	MuscleGroup string
	Equipment   string
}

func WithExerciseRowRestoreSpec(e ExerciseRowRestoreSpecification) ExerciseOption {
	return func(o *ExerciseRow) {
		o.ID = e.ID
		o.CreatedAt = e.CreatedAt
		o.Name = e.Name
		o.MuscleGroup = e.MuscleGroup
		o.Equipment = e.Equipment
	}
}

type TrainingSessionRow struct {
	ID        string    `bson:"id"`
	UserID    string    `bson:"user_id"`
	Date      time.Time `bson:"date"`
	Notes     string    `bson:"notes"`
	CreatedAt time.Time `bson:"created_at"`
}

func NewTrainingSessionRow(opts ...TrainingSessionOption) *TrainingSessionRow {
	session := &TrainingSessionRow{}

	for _, opt := range opts {
		opt(session)
	}

	return session
}

type TrainingSessionRowRestoreSpecification struct {
	ID        string
	UserID    string
	Date      time.Time
	Notes     string
	CreatedAt time.Time
}

func WithTrainingSessionRowRestoreSpec(e TrainingSessionRowRestoreSpecification) TrainingSessionOption {
	return func(o *TrainingSessionRow) {
		o.ID = e.ID
		o.UserID = e.UserID
		o.Date = e.Date
		o.CreatedAt = e.CreatedAt
		o.Notes = e.Notes
	}
}

type trainingSessionRowWithLogs struct {
	ID        string    `bson:"id"`
	UserID    string    `bson:"user_id"`
	Date      time.Time `bson:"date"`
	Notes     string    `bson:"notes,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	Logs      []SetRow  `bson:"logs"`
}

type SetRow struct {
	ID             string    `bson:"id"`
	UserID         string    `bson:"user_id"`
	SessionID      string    `bson:"session_id"`
	SessionDate    time.Time `bson:"session_date"`
	ExerciseID     string    `bson:"exercise_id"`
	ExerciseName   string    `bson:"exercise_name"`
	ExerciseNumber uint8     `bson:"exercise_number"`
	SetNumber      uint8     `bson:"set_number"`
	Weight         float32   `bson:"weight"`
	Reps           uint8     `bson:"reps"`
	Difficulty     string    `bson:"difficulty"`
	Notes          string    `bson:"notes"`
	MuscleGroup    string    `bson:"muscle_group"`
	CreatedAt      time.Time `bson:"created_at"`
}

func NewSetRow(opts ...SetOption) *SetRow {
	set := &SetRow{}

	for _, opt := range opts {
		opt(set)
	}

	return set
}

type SetRowRestoreSpecification struct {
	ID             string
	UserID         string
	SessionID      string
	SessionDate    time.Time
	ExerciseID     string
	ExerciseName   string
	ExerciseNumber uint8
	SetNumber      uint8
	Number         uint8
	Weight         float32
	Reps           uint8
	Difficulty     string
	Notes          string
	MuscleGroup    string
	CreatedAt      time.Time
}

func WithSetRestoreSpec(s SetRowRestoreSpecification) SetOption {
	return func(o *SetRow) {
		o.ID = s.ID
		o.UserID = s.UserID
		o.SessionID = s.SessionID
		o.SessionDate = s.SessionDate
		o.ExerciseID = s.ExerciseID
		o.ExerciseName = s.ExerciseName
		o.ExerciseNumber = s.ExerciseNumber
		o.SetNumber = s.SetNumber
		o.Weight = s.Weight
		o.Reps = s.Reps
		o.Difficulty = s.Difficulty
		o.Notes = s.Notes
		o.MuscleGroup = s.MuscleGroup
		o.CreatedAt = s.CreatedAt
	}
}
