package clickhouse

import (
	"gymnote/internal/entity"
	"time"

	"github.com/google/uuid"
)

type ExerciseRow struct {
	ID          uuid.UUID `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	Name        string    `db:"name"`
	MuscleGroup string    `db:"muscle_group"`
	Equipment   string    `db:"equipment"`
}

func (e *ExerciseRow) ToEntity() *entity.Exercise {
	return entity.NewExercise(
		entity.WithExerciseRestoreSpec(entity.ExerciseRestoreSpecification{
			ID:          e.ID,
			Name:        e.Name,
			MuscleGroup: e.MuscleGroup,
			Equipment:   e.Equipment,
			CreatedAt:   e.CreatedAt,
		}),
	)
}
