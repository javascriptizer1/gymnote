package clickhouse

import (
	"context"
	"fmt"

	"gymnote/internal/entity"
)

func (c *clickHouse) InsertTrainingSession(ctx context.Context, req entity.TrainingSession) error {
	query := "INSERT INTO training_sessions (id, user_id, date, total_volume, exercise_count, notes) VALUES (?, ?, ?, ?, ?, ?)"
	if err := c.conn.Exec(ctx, query, req.ID, req.UserID, req.Date, req.TotalVolume, req.ExerciseCount, req.Notes); err != nil {
		return fmt.Errorf("failed to insert training session: %w", err)
	}

	return nil
}
