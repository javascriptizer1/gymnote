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
	cache  repository.Cache
	parser Parser
}

func New(db repository.DB, cache repository.Cache, parser Parser) *service {
	return &service{
		db:     db,
		cache:  cache,
		parser: parser,
	}
}

func (s *service) ParseTraining(ctx context.Context, e entity.Event) error {
	if e.UserID == "" || e.Text == "" {
		return errors.New("invalid event data: missing UserID or Text")
	}

	parsedExercises, err := s.parser.ParseExercises(e.Text)
	if err != nil {
		return fmt.Errorf("failed to parse exercises: %w", err)
	}

	session := entity.TrainingSession{
		// ID:            uuid.New(),
		// UserID:        e.UserID,
		// Date:          time.Now(),
		// Notes:         "",
		// ExerciseCount: uint8(len(parsedExercises)),
		// Exercises:     make([]entity.SessionExercise, 0, len(parsedExercises)),
	}

	for _, parsedExercise := range parsedExercises {
		sets := make([]entity.Set, 0, len(parsedExercise.Sets))
		totalVolume := float32(0)

		_, err := s.db.GetExerciseByName(ctx, parsedExercise.Name)
		if err != nil {
			return fmt.Errorf("failed to get exercise ID for '%s': %w", parsedExercise.Name, err)
		}

		for setIDX, set := range parsedExercise.Sets {
			totalVolume += set.Weight * float32(set.Reps)
			sets = append(sets, *entity.NewSet(entity.WithSetInitSpec(
				entity.SetInitSpecification{
					UserID: e.UserID,
					// ExerciseID: exercise.ID,
					Number:     uint8(setIDX + 1),
					Weight:     set.Weight,
					Reps:       set.Reps,
					Difficulty: set.Difficulty,
					Notes:      set.Notes,
				})),
			)
		}

		// session.Exercises = append(session.Exercises, entity.SessionExercise{
		// 	ID:             uuid.New(),
		// 	Exercise:       exercise,
		// 	ExerciseNumber: uint8(exsIDX + 1),
		// 	Sets:           sets,
		// })
	}

	if err := s.db.InsertTrainingSession(ctx, session); err != nil {
		return fmt.Errorf("failed to insert training session: %w", err)
	}

	if err := s.db.InsertTrainingLogs(ctx, session); err != nil {
		return fmt.Errorf("failed to insert training logs: %w", err)
	}

	return nil
}

func (s *service) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error) {
	return s.db.GetExercisesByMuscleGroup(ctx, muscleGroup)
}

func (s *service) StartTraining(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	session, err := s.cache.GetSession(ctx, userID)
	if err == nil && session != nil {
		return nil, errors.New("training is already started")
	}

	session = entity.NewTrainingSession(entity.WithTrainingSessionInitSpec(entity.TrainingSessionInitSpecification{
		UserID:    userID,
		Date:      time.Now(),
		Exercises: []entity.SessionExercise{},
		Notes:     "",
	}))

	err = s.cache.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *service) AddExercise(ctx context.Context, userID string, exerciseID uuid.UUID) error {
	session, err := s.cache.GetSession(ctx, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}

	exercise, err := s.db.GetExerciseByID(ctx, exerciseID)
	if err != nil {
		return err
	}

	sets := []entity.Set{*entity.NewSet(entity.WithSetInitSpec(entity.SetInitSpecification{
		UserID:     session.UserID(),
		ExerciseID: exerciseID,
		Number:     1,
	}))}

	session.AddExercise(entity.NewSessionExercise(
		&exercise,
		sets,
		entity.WithSessionExerciseInitSpec(entity.SessionExerciseInitSpecification{
			Number: session.ExerciseCount() + 1,
		})),
	)

	return s.cache.SaveSession(ctx, session)
}

func (s *service) AddSet(ctx context.Context, userID string, set *entity.Set) error {
	session, err := s.cache.GetSession(ctx, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}

	// exercise := session.ActiveExercise()

	// exercise.AddSet()

	return s.cache.SaveSession(ctx, session)
}

func (s *service) UpdateActiveSet(ctx context.Context, userID string, weight float32, reps uint8) error {
	session, err := s.cache.GetSession(ctx, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}

	activeExercise := session.ActiveExercise()
	if activeExercise == nil {
		return errors.New("exercise not found")
	}

	activeSet := activeExercise.ActiveSet()
	if activeSet == nil {
		return errors.New("set not found")
	}

	activeSet.SetWeight(weight)
	activeSet.SetReps(reps)

	return s.cache.SaveSession(ctx, session)
}

func (s *service) EndSession(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	session, err := s.cache.GetSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := s.db.InsertTrainingSession(ctx, *session); err != nil {
		return nil, fmt.Errorf("failed to insert training session: %w", err)
	}

	if err := s.db.InsertTrainingLogs(ctx, *session); err != nil {
		return nil, fmt.Errorf("failed to insert training logs: %w", err)
	}

	err = s.cache.DeleteSession(ctx, userID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *service) GetCurrentSession(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	return s.cache.GetSession(ctx, userID)
}
