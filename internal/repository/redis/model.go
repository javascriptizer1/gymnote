package redis

import (
	"time"

	"gymnote/internal/entity"

	"github.com/google/uuid"
)

type TrainingSessionRow struct {
	ID        uuid.UUID            `json:"id"`
	UserID    string               `json:"user_id"`
	Date      time.Time            `json:"date"`
	Exercises []SessionExerciseRow `json:"exercises"`
	Notes     string               `json:"notes"`
	CreatedAt time.Time            `json:"created_at"`
}

func (ts *TrainingSessionRow) ToEntity() *entity.TrainingSession {
	exercises := make([]entity.SessionExercise, 0, len(ts.Exercises))
	for _, exc := range ts.Exercises {
		exercises = append(exercises, *exc.ToEntity())
	}
	return entity.NewTrainingSession(entity.WithTrainingSessionRestoreSpec(
		entity.TrainingSessionRestoreSpecification{
			ID:        ts.ID,
			UserID:    ts.UserID,
			Date:      ts.Date,
			Notes:     ts.Notes,
			Exercises: exercises,
			CreatedAt: ts.CreatedAt,
		},
	))
}

func NewTrainingSessionRow(session *entity.TrainingSession) *TrainingSessionRow {
	exercises := make([]SessionExerciseRow, 0, len(session.Exercises()))
	for _, exc := range session.Exercises() {
		exercises = append(exercises, *NewSessionExerciseRow(&exc))
	}
	return &TrainingSessionRow{
		ID:        session.ID(),
		UserID:    session.UserID(),
		Date:      session.Date(),
		Exercises: exercises,
		Notes:     session.Notes(),
		CreatedAt: session.CreatedAt(),
	}
}

type SessionExerciseRow struct {
	ID                  uuid.UUID `json:"id"`
	ExerciseID          uuid.UUID `json:"exercise_id"`
	ExerciseName        string    `json:"exercise_name"`
	ExerciseMuscleGroup string    `json:"exercise_muscle_group"`
	ExerciseEquipment   string    `json:"exercise_equipment"`
	ExerciseCreatedAt   time.Time `json:"exercise_created_at"`
	Number              uint8     `json:"number"`
	Sets                []SetRow  `json:"sets"`
}

func (e *SessionExerciseRow) ToEntity() *entity.SessionExercise {
	sets := make([]entity.Set, 0, len(e.Sets))
	for _, s := range e.Sets {
		sets = append(sets, *s.ToEntity())
	}
	return entity.NewSessionExercise(
		entity.NewExercise(entity.WithExerciseRestoreSpec(entity.ExerciseRestoreSpecification{
			ID:          e.ExerciseID,
			Name:        e.ExerciseName,
			MuscleGroup: e.ExerciseMuscleGroup,
			Equipment:   e.ExerciseEquipment,
			CreatedAt:   e.ExerciseCreatedAt,
		})),
		sets,
		entity.WithSessionExerciseRestoreSpec(entity.SessionExerciseRestoreSpecification{
			ID:     e.ID,
			Number: e.Number,
		}))
}

func NewSessionExerciseRow(exercise *entity.SessionExercise) *SessionExerciseRow {
	sets := make([]SetRow, 0, len(exercise.Sets()))
	for _, s := range exercise.Sets() {
		sets = append(sets, *NewSetRow(&s))
	}
	return &SessionExerciseRow{
		ID:                  exercise.ID(),
		Number:              exercise.Number(),
		ExerciseID:          exercise.Exercise.ID(),
		ExerciseName:        exercise.Exercise.Name(),
		ExerciseMuscleGroup: exercise.Exercise.MuscleGroup(),
		ExerciseEquipment:   exercise.Exercise.Equipment(),
		ExerciseCreatedAt:   exercise.Exercise.CreatedAt(),
		Sets:                sets,
	}
}

type SetRow struct {
	ID         uuid.UUID `json:"id"`
	UserID     string    `json:"user_id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
	Number     uint8     `json:"number"`
	Weight     float32   `json:"weight"`
	Reps       uint8     `json:"reps"`
	Difficulty string    `json:"difficulty"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *SetRow) ToEntity() *entity.Set {
	return entity.NewSet(entity.WithSetRestoreSpec(entity.SetRestoreSpecification{
		ID:         s.ID,
		UserID:     s.UserID,
		ExerciseID: s.ExerciseID,
		Number:     s.Number,
		Weight:     s.Weight,
		Reps:       s.Reps,
		Difficulty: s.Difficulty,
		Notes:      s.Notes,
		CreatedAt:  s.CreatedAt,
	}))
}

func NewSetRow(set *entity.Set) *SetRow {
	return &SetRow{
		ID:         set.ID(),
		UserID:     set.UserID(),
		ExerciseID: set.ExerciseID(),
		Number:     set.Number(),
		Weight:     set.Weight(),
		Reps:       set.Reps(),
		Difficulty: set.Difficulty(),
		Notes:      set.Notes(),
		CreatedAt:  set.CreatedAt(),
	}
}
