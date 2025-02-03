package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gymnote/internal/entity"
	"gymnote/internal/errs"
	"gymnote/internal/parser"
	"gymnote/internal/repository"
)

type Parser interface {
	ParseExercises(s string) ([]parser.Exercise, time.Time, error)
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

func (s *service) ParseTraining(ctx context.Context, e entity.Event) (*entity.TrainingSession, error) {
	if e.UserID == "" || e.Text == "" {
		return nil, errs.ErrInvalidEventData
	}

	parsedExercises, date, err := s.parser.ParseExercises(e.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse exercises: %w", err)
	}

	var exercises []entity.SessionExercise

	for exsIDX, parsedExercise := range parsedExercises {
		sets := make([]entity.Set, 0, len(parsedExercise.Sets))

		exercise, err := s.db.GetExerciseByName(ctx, parsedExercise.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get exercise ID for '%s': %w", parsedExercise.Name, err)
		}

		for setIDX, set := range parsedExercise.Sets {
			sets = append(sets, *entity.NewSet(entity.WithSetInitSpec(
				entity.SetInitSpecification{
					ExerciseID: exercise.ID(),
					UserID:     e.UserID,
					Number:     uint8(setIDX + 1),
					Weight:     set.Weight,
					Reps:       set.Reps,
					Difficulty: set.Difficulty,
					Notes:      set.Notes,
				})),
			)
		}

		exercises = append(exercises, *entity.NewSessionExercise(&exercise, sets, entity.WithSessionExerciseInitSpec(
			entity.SessionExerciseInitSpecification{
				Number: uint8(exsIDX + 1),
			},
		)))
	}

	session := entity.NewTrainingSession(entity.WithTrainingSessionInitSpec(
		entity.TrainingSessionInitSpecification{
			UserID:    e.UserID,
			Date:      date,
			Notes:     "",
			Exercises: exercises,
		},
	))

	if err := s.db.InsertTrainingSession(ctx, *session); err != nil {
		return nil, fmt.Errorf("failed to insert training session: %w", err)
	}

	if err := s.db.InsertTrainingLogs(ctx, *session); err != nil {
		return nil, fmt.Errorf("failed to insert training logs: %w", err)
	}

	return session, nil
}

func (s *service) GetExerciseProgression(ctx context.Context, userID string, exerciseID uuid.UUID) ([]entity.ExerciseProgression, error) {
	now := time.Now()
	fromDate := now.AddDate(-1, 0, 0)
	toDate := now

	return s.db.GetExerciseProgression(ctx, userID, exerciseID, fromDate, toDate)
}

func (s *service) CreateExercise(ctx context.Context, name, muscleGroup, equipment string) error {
	_, err := s.db.GetExerciseByName(ctx, name)
	if err == nil {
		return errs.ErrExerciseAlreadyExists
	}
	if !errors.Is(err, errs.ErrExerciseNotFound) {
		return fmt.Errorf("failed to check existing exercise: %w", err)
	}

	exercise := entity.NewExercise(entity.WithExerciseInitSpec(entity.ExerciseInitSpecification{
		Name:        name,
		MuscleGroup: muscleGroup,
		Equipment:   equipment,
	}))

	return s.db.InsertExercise(ctx, *exercise)
}

func (s *service) GetTrainingSessions(ctx context.Context, userID string, fromDate, toDate *time.Time) ([]entity.TrainingSession, error) {
	now := time.Now()

	if fromDate == nil {
		defaultFrom := now.AddDate(0, 0, -14)
		fromDate = &defaultFrom
	}

	if toDate == nil {
		toDate = &now
	}

	sessions, err := s.db.GetTrainingSessions(ctx, userID, *fromDate, *toDate)
	return sessions, err
}

func (s *service) GetExercisesByMuscleGroup(ctx context.Context, muscleGroup string) ([]entity.Exercise, error) {
	return s.db.GetExercisesByMuscleGroup(ctx, muscleGroup)
}

func (s *service) StartTraining(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	session, err := s.cache.GetSession(ctx, userID)
	if err == nil && session != nil {
		return nil, errs.ErrTrainingStarted
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

func (s *service) AddTrainingExercise(ctx context.Context, userID string, exerciseID uuid.UUID) error {
	session, err := s.getSession(ctx, userID)
	if err != nil {
		return err
	}

	return s.addExerciseToSession(ctx, session, exerciseID)
}

func (s *service) AddOrUpdateSet(ctx context.Context, userID string, weight float32, reps uint8, notes string) error {
	session, err := s.getSession(ctx, userID)
	if err != nil {
		return err
	}

	activeExercise := session.ActiveExercise()
	if activeExercise == nil {
		return errs.ErrExerciseNotFound
	}

	lastSet := activeExercise.LastSet()
	if lastSet == nil {
		return errs.ErrSetNotFound
	}

	if lastSet.Weight() == 0 || lastSet.Reps() == 0 {
		lastSet.SetWeight(weight)
		lastSet.SetReps(reps)
		lastSet.SetNotes(notes)
		return s.cache.SaveSession(ctx, session)
	}

	newSet := entity.NewSet(entity.WithSetInitSpec(
		entity.SetInitSpecification{
			UserID:     lastSet.UserID(),
			ExerciseID: lastSet.ExerciseID(),
			Number:     lastSet.Number() + 1,
			Weight:     weight,
			Reps:       reps,
			Notes:      notes,
		},
	))

	activeExercise.AddSet(newSet)

	return s.cache.SaveSession(ctx, session)
}

func (s *service) EndSession(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	session, err := s.getSession(ctx, userID)
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

func (s *service) ClearSession(ctx context.Context, userID string) error {
	session, err := s.cache.GetSession(ctx, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return errs.ErrSessionNotFound
	}

	return s.cache.DeleteSession(ctx, userID)
}

func (s *service) GetCurrentSession(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	return s.cache.GetSession(ctx, userID)
}

func (s *service) getSession(ctx context.Context, userID string) (*entity.TrainingSession, error) {
	session, err := s.cache.GetSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errs.ErrSessionNotFound
	}

	return session, nil
}

func (s *service) addExerciseToSession(ctx context.Context, session *entity.TrainingSession, exerciseID uuid.UUID) error {
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
