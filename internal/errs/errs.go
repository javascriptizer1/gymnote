package errs

import (
	"fmt"
)

var (
	ErrExerciseAlreadyExists  = fmt.Errorf("exercise already exists")
	ErrInvalidEventData       = fmt.Errorf("invalid event data: missing UserID or Text")
	ErrTrainingStarted        = fmt.Errorf("training is already started")
	ErrSessionNotFound        = fmt.Errorf("session not found")
	ErrExerciseNotFound       = fmt.Errorf("exercise not found")
	ErrSetNotFound            = fmt.Errorf("set not found")
	ErrFailedToInsertData     = fmt.Errorf("failed to insert training data")
	ErrAccountNotFound        = fmt.Errorf("account not found")
	ErrLinkCodeNotFound       = fmt.Errorf("link code not found or expired")
	ErrAccountAlreadyLinked   = fmt.Errorf("account is already linked")
	ErrCannotLinkSamePlatform = fmt.Errorf("cannot link accounts on the same platform")
)
