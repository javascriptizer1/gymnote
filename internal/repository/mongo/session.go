package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"gymnote/internal/entity"
)

func (m *mongodb) InsertTrainingSession(ctx context.Context, req entity.TrainingSession) error {
	row := NewTrainingSessionRow(WithTrainingSessionRowRestoreSpec(TrainingSessionRowRestoreSpecification{
		ID:        req.ID().String(),
		UserID:    req.UserID(),
		Date:      req.Date(),
		Notes:     req.Notes(),
		CreatedAt: req.CreatedAt(),
	}))

	_, err := m.sessionColl.InsertOne(ctx, row)
	if err != nil {
		return fmt.Errorf("failed to insert training session: %w", err)
	}

	return nil
}

func (m *mongodb) GetTrainingSessions(ctx context.Context, userID string, fromDate, toDate time.Time) ([]entity.TrainingSession, error) {
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "date", Value: bson.D{
				{Key: "$gte", Value: fromDate},
				{Key: "$lte", Value: toDate},
			}},
		}},
	}

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: colLogs},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "session_id"},
			{Key: "as", Value: "logs"},
		}},
	}

	cursor, err := m.sessionColl.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage})
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(ctx)

	var rawSessions []trainingSessionRowWithLogs
	if err := cursor.All(ctx, &rawSessions); err != nil {
		return nil, fmt.Errorf("cursor decode error: %w", err)
	}

	var sessions []entity.TrainingSession

	for _, raw := range rawSessions {
		sessionID, _ := uuid.Parse(raw.ID)
		exercisesMap := make(map[string]*entity.SessionExercise)

		for _, log := range raw.Logs {
			exID, _ := uuid.Parse(log.ExerciseID)
			setID, _ := uuid.Parse(log.ID)

			exKey := fmt.Sprintf("%s-%d", exID.String(), log.ExerciseNumber)

			if _, ok := exercisesMap[exKey]; !ok {
				ex := entity.NewExercise(entity.WithExerciseRestoreSpec(entity.ExerciseRestoreSpecification{
					ID:          exID,
					Name:        log.ExerciseName,
					MuscleGroup: log.MuscleGroup,
					CreatedAt:   log.CreatedAt,
				}))

				exercisesMap[exKey] = entity.NewSessionExercise(ex, nil, entity.WithSessionExerciseRestoreSpec(
					entity.SessionExerciseRestoreSpecification{
						Number: uint8(log.ExerciseNumber),
					},
				))
			}

			set := entity.NewSet(entity.WithSetRestoreSpec(entity.SetRestoreSpecification{
				ID:         setID,
				UserID:     raw.UserID,
				ExerciseID: exID,
				Number:     log.SetNumber,
				Weight:     log.Weight,
				Reps:       log.Reps,
				Difficulty: log.Difficulty,
				Notes:      log.Notes,
				CreatedAt:  log.CreatedAt,
			}))

			exercisesMap[exKey].AddSet(set)
		}

		var exercises []entity.SessionExercise
		for _, ex := range exercisesMap {
			exercises = append(exercises, *ex)
		}

		session := *entity.NewTrainingSession(entity.WithTrainingSessionRestoreSpec(entity.TrainingSessionRestoreSpecification{
			ID:        sessionID,
			UserID:    raw.UserID,
			Date:      raw.Date,
			Exercises: exercises,
			Notes:     raw.Notes,
			CreatedAt: raw.CreatedAt,
		}))

		sessions = append(sessions, session)
	}

	return sessions, nil
}
