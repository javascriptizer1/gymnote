package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gymnote/internal/entity"
)

func (c *clickHouse) InsertTrainingLogs(ctx context.Context, req entity.TrainingSession) error {
	batch, err := c.conn.PrepareBatch(ctx, `
		INSERT INTO training_logs (id, user_id, session_id, exercise_id, session_date, exercise_name, exercise_number, set_number, weight, reps, difficulty, notes, muscle_group, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, exs := range req.Exercises() {
		for _, log := range exs.Sets() {
			if err := batch.Append(log.ID(), log.UserID(), req.ID(), log.ExerciseID(), req.Date(), exs.Exercise.Name(), exs.Number(),
				log.Number(), log.Weight(), log.Reps(), log.Difficulty(), log.Notes(), exs.Exercise.MuscleGroup(), log.CreatedAt()); err != nil {
				return fmt.Errorf("failed to append training log: %w", err)
			}
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

func (c *clickHouse) GetExerciseProgression(ctx context.Context, userID string, exerciseID uuid.UUID, fromDate, toDate time.Time) ([]entity.ExerciseProgression, error) {
	var result []entity.ExerciseProgression

	query := `
		SELECT exercise_name, session_date, MAX(weight) AS weight, MAX(reps) AS reps
		FROM training_logs
		WHERE user_id = ? AND exercise_id = ? AND session_date BETWEEN ? AND ?
		GROUP BY exercise_name, session_date
		ORDER BY session_date;
	`

	rows, err := c.conn.Query(ctx, query, userID, exerciseID, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise progression: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name   string
			date   time.Time
			weight float32
			reps   uint8
		)
		if err := rows.Scan(&name, &date, &weight, &reps); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		result = append(result, entity.ExerciseProgression{
			ExerciseName: name,
			SessionDate:  date,
			Weight:       weight,
			Reps:         reps,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return result, nil
}

func (c *clickHouse) GetLastSetsForExercise(ctx context.Context, userID string, exerciseID uuid.UUID, limitDays int64) ([]entity.ExerciseProgression, error) {
	var result []entity.ExerciseProgression

	query := `
        SELECT session_date, weight, reps
        FROM training_logs
        WHERE user_id = ? AND exercise_id = ? AND session_date IN (
            SELECT DISTINCT session_date
            FROM training_logs
            WHERE user_id = ? AND exercise_id = ?
            ORDER BY session_date DESC
            LIMIT ?
        )
        ORDER BY session_date ASC, set_number ASC
    `

	rows, err := c.conn.Query(ctx, query, userID, exerciseID, userID, exerciseID, limitDays)
	if err != nil {
		return nil, fmt.Errorf("failed to get last sets for exercise: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			date   time.Time
			weight float32
			reps   uint8
		)
		if err := rows.Scan(&date, &weight, &reps); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}

		result = append(result, entity.ExerciseProgression{
			SessionDate: date,
			Weight:      weight,
			Reps:        reps,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return result, nil
}
