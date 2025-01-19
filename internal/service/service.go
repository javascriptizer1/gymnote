package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gymnote/internal/entity"
	"gymnote/internal/parser"
	"gymnote/internal/repository"
)

type Parser interface {
	ParseExercises(s string) ([]parser.Exercise, error)
}

type service struct {
	db     repository.DB
	parser Parser
}

func New(db repository.DB, parser Parser) *service {
	return &service{
		db:     db,
		parser: parser,
	}
}

func (s *service) Process(ctx context.Context, e entity.Event) error {
	if e.UserID == "" || e.Text == "" {
		return errors.New("invalid event data: missing UserID or Text")
	}

	parsedExercises, err := s.parser.ParseExercises(e.Text)
	if err != nil {
		return fmt.Errorf("failed to parse exercises: %w", err)
	}

	session := entity.TrainingSession{
		ID:            uuid.New(),
		UserID:        e.UserID,
		Date:          time.Now(),
		Notes:         "",
		ExerciseCount: uint8(len(parsedExercises)),
		Exercises:     make([]entity.SessionExercise, 0, len(parsedExercises)),
	}

	for exsIDX, parsedExercise := range parsedExercises {
		sets := make([]entity.Set, 0, len(parsedExercise.Sets))
		totalVolume := float32(0)

		exercise, err := s.db.GetExerciseByName(ctx, parsedExercise.Name)
		if err != nil {
			return fmt.Errorf("failed to get exercise ID for '%s': %w", parsedExercise.Name, err)
		}

		for setIDX, set := range parsedExercise.Sets {
			totalVolume += set.Weight * float32(set.Reps)
			sets = append(sets, entity.Set{
				ID:         uuid.New(),
				UserID:     e.UserID,
				ExerciseID: exercise.ID,
				SetNumber:  uint8(setIDX + 1),
				Weight:     set.Weight,
				Reps:       set.Reps,
				Difficulty: set.Difficulty,
				Notes:      set.Notes,
			})
		}

		session.TotalVolume += totalVolume
		session.Exercises = append(session.Exercises, entity.SessionExercise{
			ID:             uuid.New(),
			Exercise:       exercise,
			ExerciseNumber: uint8(exsIDX + 1),
			Sets:           sets,
		})
	}

	if err := s.db.InsertTrainingSession(ctx, session); err != nil {
		return fmt.Errorf("failed to insert training session: %w", err)
	}

	if err := s.db.InsertTrainingLogs(ctx, session); err != nil {
		return fmt.Errorf("failed to insert training logs: %w", err)
	}

	return nil
}
