package clickhouse

import (
	"context"
	"fmt"

	"gymnote/internal/entity"
)

func (c *clickHouse) InsertTrainingLogs(ctx context.Context, req entity.TrainingSession) error {
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO training_logs (id, user_id, session_date, exercise_id, set_number, weight, reps, difficulty, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, exs := range req.Exercises {
		for _, log := range exs.Sets {
			if err := batch.Append(log.ID, log.UserID, req.Date, log.ExerciseID, log.SetNumber, log.Weight, log.Reps, log.Difficulty, log.Notes); err != nil {
				return fmt.Errorf("failed to append training log: %w", err)
			}
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}
