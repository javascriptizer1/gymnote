package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gymnote/internal/entity"
	"gymnote/internal/errs"

	"github.com/google/uuid"
)

func (c *clickHouse) InsertExercise(ctx context.Context, req entity.Exercise) error {
	query := "INSERT INTO exercises (id, name, muscle_group, equipment, created_at) VALUES (?, ?, ?, ?, ?)"
	if err := c.conn.Exec(ctx, query, req.ID(), req.Name(), req.MuscleGroup(), req.Equipment(), req.CreatedAt()); err != nil {
		return fmt.Errorf("failed to insert exercise: %w", err)
	}

	return nil
}

func (c *clickHouse) GetExerciseByName(ctx context.Context, name string) (entity.Exercise, error) {
	var row ExerciseRow

	query := "SELECT id, name, muscle_group, equipment, created_at FROM exercises WHERE name = ? LIMIT 1"

	err := c.conn.QueryRow(ctx, query, name).Scan(
		&row.ID,
		&row.Name,
		&row.MuscleGroup,
		&row.Equipment,
		&row.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Exercise{}, errs.ErrExerciseNotFound
		}
		return entity.Exercise{}, fmt.Errorf("failed to get exercise by name: %w", err)
	}

	return *row.ToEntity(), nil
}

func (c *clickHouse) GetExerciseByID(ctx context.Context, id uuid.UUID) (entity.Exercise, error) {
	var row ExerciseRow

	query := "SELECT id, name, muscle_group, equipment, created_at FROM exercises WHERE id = ? LIMIT 1"

	err := c.conn.QueryRow(ctx, query, id).Scan(
		&row.ID,
		&row.Name,
		&row.MuscleGroup,
		&row.Equipment,
		&row.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Exercise{}, errs.ErrExerciseNotFound
		}
		return entity.Exercise{}, fmt.Errorf("failed to get exercise by id: %w", err)
	}

	return *row.ToEntity(), nil
}

func (c *clickHouse) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error) {
	var result []entity.Exercise

	query := "SELECT id, name, muscle_group, equipment, created_at FROM exercises WHERE muscle_group = ? ORDER BY created_at DESC"

	rows, err := c.conn.Query(ctx, query, muscleGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises by muscle_group: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row ExerciseRow
		if err := rows.Scan(&row.ID, &row.Name, &row.MuscleGroup, &row.Equipment, &row.CreatedAt); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		result = append(result, *row.ToEntity())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return result, nil
}
