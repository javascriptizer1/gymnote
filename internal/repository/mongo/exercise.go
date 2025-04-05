package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"gymnote/internal/entity"
	"gymnote/internal/errs"
)

func (m *mongodb) InsertExercise(ctx context.Context, req entity.Exercise) error {
	row := NewExerciseRow(WithExerciseRowRestoreSpec(ExerciseRowRestoreSpecification{
		ID:          req.ID().String(),
		CreatedAt:   req.CreatedAt(),
		Name:        req.Name(),
		MuscleGroup: req.MuscleGroup(),
		Equipment:   req.Equipment(),
	}))

	_, err := m.exerciseColl.InsertOne(ctx, row)
	if err != nil {
		return fmt.Errorf("failed to insert exercise: %w", err)
	}

	return nil
}

func (m *mongodb) GetExerciseByName(ctx context.Context, name string) (entity.Exercise, error) {
	var row ExerciseRow
	filter := bson.M{"name": name}

	err := m.exerciseColl.FindOne(ctx, filter).Decode(&row)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Exercise{}, errs.ErrExerciseNotFound
		}
		return entity.Exercise{}, fmt.Errorf("failed to get exercise by name: %w", err)
	}

	return *row.ToEntity(), nil
}

func (m *mongodb) GetExerciseByID(ctx context.Context, id uuid.UUID) (entity.Exercise, error) {
	var row ExerciseRow
	filter := bson.M{"id": id.String()}

	err := m.exerciseColl.FindOne(ctx, filter).Decode(&row)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Exercise{}, errs.ErrExerciseNotFound
		}
		return entity.Exercise{}, fmt.Errorf("failed to get exercise by id: %w", err)
	}

	return *row.ToEntity(), nil
}

func (m *mongodb) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error) {
	filter := bson.M{"muscle_group": muscleGroup}
	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := m.exerciseColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("close cursor err: %v", err)
		}
	}()

	var exercises []entity.Exercise
	for cursor.Next(ctx) {
		var row ExerciseRow
		if err := cursor.Decode(&row); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}
		exercises = append(exercises, *row.ToEntity())
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return exercises, nil
}
