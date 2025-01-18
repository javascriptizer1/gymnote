package clickhouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (c *clickHouse) GetExerciseIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	var id string

	query := "SELECT id FROM exercises WHERE name = ? LIMIT 1"

	err := c.conn.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get exercise ID by name: %w", err)
	}

	guid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("received id is not uuid: %w", err)
	}

	return guid, nil
}
