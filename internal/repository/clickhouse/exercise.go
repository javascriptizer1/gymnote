package clickhouse

import (
	"context"
	"fmt"
	"gymnote/internal/entity"
)

func (c *clickHouse) GetExerciseByName(ctx context.Context, name string) (entity.Exercise, error) {
	var exercise entity.Exercise

	query := "SELECT id, name, muscle_group, equipment, createdAt FROM exercises WHERE name = ? LIMIT 1"

	err := c.conn.QueryRow(ctx, query, name).Scan(&exercise.ID, &exercise.Name, &exercise.MuscleGroup, &exercise.Equipment, &exercise.CreatedAt)
	if err != nil {
		return exercise, fmt.Errorf("failed to get exercise by name: %w", err)
	}

	return exercise, nil
}
