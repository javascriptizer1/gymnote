package errs

import (
	"fmt"
)

var (
	ErrExerciseAlreadyExists = fmt.Errorf("exercise already exists")
	ErrInvalidEventData      = fmt.Errorf("invalid event data: missing UserID or Text")
	ErrTrainingStarted       = fmt.Errorf("training is already started")
	ErrSessionNotFound       = fmt.Errorf("session not found")
	ErrExerciseNotFound      = fmt.Errorf("exercise not found")
	ErrSetNotFound           = fmt.Errorf("set not found")
	ErrFailedToInsertData    = fmt.Errorf("failed to insert training data")
)
