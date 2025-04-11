package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"gymnote/internal/entity"
)

func (m *mongodb) InsertTrainingLogs(ctx context.Context, req entity.TrainingSession) error {
	docsToWrite := make([]any, 0, req.SetCount())
	for _, exs := range req.Exercises() {
		for _, set := range exs.Sets() {
			docsToWrite = append(docsToWrite, NewSetRow(WithSetRestoreSpec(SetRowRestoreSpecification{
				ID:             set.ID().String(),
				UserID:         set.UserID(),
				SessionID:      req.ID().String(),
				SessionDate:    req.Date(),
				ExerciseName:   exs.Name(),
				ExerciseNumber: exs.Number(),
				SetNumber:      set.Number(),
				MuscleGroup:    exs.MuscleGroup(),
				ExerciseID:     set.ExerciseID().String(),
				Number:         set.Number(),
				Weight:         set.Weight(),
				Reps:           set.Reps(),
				Difficulty:     set.Difficulty(),
				Notes:          set.Notes(),
				CreatedAt:      set.CreatedAt(),
			})))
		}
	}

	if _, err := m.logColl.InsertMany(ctx, docsToWrite); err != nil {
		return fmt.Errorf("insert exercise error: %w", err)
	}

	return nil
}

func (m *mongodb) GetExerciseProgression(ctx context.Context, userID string, exerciseID uuid.UUID, fromDate, toDate time.Time) ([]entity.ExerciseProgression, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "exercise_id", Value: exerciseID.String()},
			{Key: "session_date", Value: bson.D{
				{Key: "$gte", Value: fromDate},
				{Key: "$lte", Value: toDate},
			}},
		}}},

		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "exercise_name", Value: "$exercise_name"},
				{Key: "session_date", Value: "$session_date"},
			}},
			{Key: "max_weight", Value: bson.D{{Key: "$max", Value: "$weight"}}},
			{Key: "max_reps", Value: bson.D{{Key: "$max", Value: "$reps"}}},
		}}},

		{{Key: "$sort", Value: bson.D{{Key: "_id.session_date", Value: 1}}}},
	}

	cursor, err := m.logColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate exercise progression: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("close cursor err: %v", err)
		}
	}()

	var result []entity.ExerciseProgression

	for cursor.Next(ctx) {
		var progress struct {
			ID struct {
				ExerciseName string    `bson:"exercise_name"`
				SessionDate  time.Time `bson:"session_date"`
			} `bson:"_id"`
			MaxWeight float32 `bson:"max_weight"`
			MaxReps   uint8   `bson:"max_reps"`
		}

		if err := cursor.Decode(&progress); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}

		result = append(result, entity.ExerciseProgression{
			ExerciseName: progress.ID.ExerciseName,
			SessionDate:  progress.ID.SessionDate,
			Weight:       progress.MaxWeight,
			Reps:         progress.MaxReps,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}

func (m *mongodb) GetLastSetsForExercise(ctx context.Context, userID string, exerciseID uuid.UUID, limitDays int64) ([]entity.ExerciseProgression, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userID},
			{Key: "exercise_id", Value: exerciseID.String()},
		}}},

		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$session_date"},
		}}},

		{{Key: "$sort", Value: bson.D{
			{Key: "_id", Value: -1},
		}}},

		{{Key: "$limit", Value: limitDays}},

		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "dates", Value: bson.D{
				{Key: "$push", Value: "$_id"},
			}},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: m.logColl.Name()},

			{Key: "let", Value: bson.D{
				{Key: "dates", Value: "$dates"},
				{Key: "user_id", Value: userID},
				{Key: "exercise_id", Value: exerciseID.String()},
			}},

			{Key: "pipeline", Value: mongo.Pipeline{

				{{Key: "$match", Value: bson.D{

					{Key: "$expr", Value: bson.D{

						{Key: "$and", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$user_id", "$$user_id"}}},
							bson.D{{Key: "$eq", Value: bson.A{"$exercise_id", "$$exercise_id"}}},
							bson.D{{Key: "$in", Value: bson.A{"$session_date", "$$dates"}}},
						}},
					}},
				}}},

				{{Key: "$sort", Value: bson.D{
					{Key: "session_date", Value: 1},
					{Key: "set_number", Value: 1},
				}}},
			}},

			{Key: "as", Value: "logs"},
		}}},

		{{Key: "$unwind", Value: "$logs"}},

		{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: "$logs"},
		}}},

		{{Key: "$project", Value: bson.D{
			{Key: "session_date", Value: 1},
			{Key: "weight", Value: 1},
			{Key: "reps", Value: 1},
		}}},
	}

	cursor, err := m.logColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to get last sets for exercise: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("close cursor err: %v", err)
		}
	}()

	var result []entity.ExerciseProgression

	for cursor.Next(ctx) {
		var log struct {
			SessionDate time.Time `bson:"session_date"`
			Weight      float64   `bson:"weight"`
			Reps        uint8     `bson:"reps"`
		}

		if err := cursor.Decode(&log); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}

		result = append(result, entity.ExerciseProgression{
			SessionDate: log.SessionDate,
			Weight:      float32(log.Weight),
			Reps:        log.Reps,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return result, nil
}
