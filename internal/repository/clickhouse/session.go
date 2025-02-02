package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gymnote/internal/entity"
)

func (c *clickHouse) InsertTrainingSession(ctx context.Context, req entity.TrainingSession) error {
	query := "INSERT INTO training_sessions (id, user_id, date, total_volume, exercise_count, notes) VALUES (?, ?, ?, ?, ?, ?)"
	if err := c.conn.Exec(ctx, query, req.ID(), req.UserID(), req.Date(), req.TotalVolume(), req.ExerciseCount(), req.Notes()); err != nil {
		return fmt.Errorf("failed to insert training session: %w", err)
	}

	return nil
}

func (c *clickHouse) GetTrainingSessions(ctx context.Context, userID string, fromDate, toDate time.Time) ([]entity.TrainingSession, error) {
	var query string
	var args []interface{}

	query = `
		SELECT id, user_id, session_id, session_date, exercise_id, exercise_name, exercise_number, set_number, weight, reps, difficulty, notes, muscle_group, created_at
		FROM training_logs
		WHERE user_id = ? AND session_date BETWEEN ? AND ?
		ORDER BY session_date, exercise_number, set_number
		`
	args = append(args, userID, fromDate, toDate)

	rows, err := c.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query training logs: %w", err)
	}
	defer rows.Close()

	sessions := make(map[uuid.UUID]entity.TrainingSession)

	for rows.Next() {
		var (
			id             uuid.UUID
			sessionID      uuid.UUID
			sessionDate    time.Time
			userID         string
			exerciseID     uuid.UUID
			exerciseName   string
			exerciseNumber uint8
			setNumber      uint8
			weight         float32
			reps           uint8
			difficulty     string
			notes          string
			muscleGroup    string
			createdAt      time.Time
		)

		if err := rows.Scan(&id, &userID, &sessionID, &sessionDate, &exerciseID, &exerciseName, &exerciseNumber, &setNumber, &weight, &reps, &difficulty, &notes, &muscleGroup, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		session, exists := sessions[sessionID]
		if !exists {
			session = *entity.NewTrainingSession(entity.WithTrainingSessionRestoreSpec(
				entity.TrainingSessionRestoreSpecification{
					ID:        sessionID,
					UserID:    userID,
					Date:      sessionDate,
					Exercises: []entity.SessionExercise{},
					CreatedAt: createdAt,
				},
			))
		}

		found := false
		for i, ex := range session.Exercises() {
			if ex.Exercise.ID() == exerciseID && ex.Number() == exerciseNumber {
				session.Exercises()[i].AddSet(entity.NewSet(entity.WithSetRestoreSpec(entity.SetRestoreSpecification{
					ID:         id,
					UserID:     userID,
					ExerciseID: exerciseID,
					Number:     setNumber,
					Weight:     weight,
					Reps:       reps,
					Difficulty: difficulty,
					Notes:      notes,
					CreatedAt:  createdAt,
				})))
				found = true
				break
			}
		}

		if !found {
			session.AddExercise(entity.NewSessionExercise(
				entity.NewExercise(entity.WithExerciseRestoreSpec(
					entity.ExerciseRestoreSpecification{
						ID:          exerciseID,
						Name:        exerciseName,
						MuscleGroup: muscleGroup,
						CreatedAt:   createdAt,
					},
				)),
				[]entity.Set{*entity.NewSet(entity.WithSetRestoreSpec(entity.SetRestoreSpecification{
					ID:         id,
					UserID:     userID,
					ExerciseID: exerciseID,
					Number:     setNumber,
					Weight:     weight,
					Reps:       reps,
					Difficulty: difficulty,
					Notes:      notes,
					CreatedAt:  createdAt,
				}))},
				entity.WithSessionExerciseRestoreSpec(
					entity.SessionExerciseRestoreSpecification{
						Number: exerciseNumber,
					},
				)),
			)
		}

		sessions[sessionID] = session
	}

	var result []entity.TrainingSession
	for _, s := range sessions {
		result = append(result, s)
	}

	return result, nil
}
