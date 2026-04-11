package vk

import (
	"gymnote/internal/handler/common"
)

const (
	// commands
	startCommand                  = common.StartCommand
	helpCommand                   = common.HelpCommand
	startTrainingCommand          = common.StartTrainingCommand
	clearTrainingCommand          = common.ClearTrainingCommand
	createExerciseCommand         = common.CreateExerciseCommand
	uploadTrainingCommand         = common.UploadTrainingCommand
	getTrainingsCommand           = common.GetTrainingsCommand
	getExerciseProgressionCommand = common.GetExerciseProgressionCommand
	getExerciseHistoryCommand     = common.GetExerciseHistoryCommand
	oneRMCommand                  = common.OneRMCommand
	// callbacks
	musclePrefix                      = common.MusclePrefix
	exercisePrefix                    = common.ExercisePrefix
	confirmationFinishTrainingPrefix  = common.ConfirmationFinishTrainingPrefix
	acceptFinishTrainingPrefix        = common.AcceptFinishTrainingPrefix
	rejectFinishTrainingPrefix        = common.RejectFinishTrainingPrefix
	startNewExercisePrefix            = common.StartNewExercisePrefix
	startGetExerciseProgressionPrefix = common.StartGetExerciseProgressionPrefix
	startGetExerciseHistoryPrefix     = common.StartGetExerciseHistoryPrefix
	backToMuscleGroups                = common.BackToMuscleGroups

	nextDirection = common.NextDirection
	prevDirection = common.PrevDirection
)

const (
	startText                                 = common.StartTextVK
	helpText                                  = common.HelpTextVK
	clearTrainingDoneText                     = common.ClearTrainingDoneText
	donateAuthorText                          = common.DonateAuthorText
	startTrainingText                         = common.StartTrainingTextVK
	muscleGroupDoneText                       = common.MuscleGroupDoneTextVK
	muscleGroupSelectText                     = common.MuscleGroupSelectText
	startProgressionMuscleGroupSelectText     = common.StartProgressionMuscleGroupSelectText
	startExerciseHistoryMuscleGroupSelectText = common.StartExerciseHistoryMuscleGroupSelectText
	exerciseText                              = common.ExerciseText
	lastSetsText                              = common.LastSetsText
	setText                                   = common.SetText
	exerciseCreatedText                       = common.ExerciseCreatedText
	startNewExerciseText                      = common.StartNewExerciseText
	finishTrainingConfirmationText            = common.FinishTrainingConfirmationText
	finishTrainingText                        = common.FinishTrainingText
	finishText                                = common.FinishText
	startOneRMText                            = common.StartOneRMText
	notFoundTrainingsText                     = common.NotFoundTrainingsText
	startCreateExerciseText                   = common.StartCreateExerciseText
	startGetTrainingsText                     = common.StartGetTrainingsText
	startUploadTrainingText                   = common.StartUploadTrainingText
	paginationNextText                        = common.PaginationNextText
	paginationPrevText                        = common.PaginationPrevText
	loadingProgressionText                    = common.LoadingProgressionText
	backToMuscleGroupsText                    = common.BackToMuscleGroupsText
	backToExercisesText                       = common.BackToExercisesText

	adminOnlyText                     = common.AdminOnlyText
	answerYes                         = common.AnswerYes
	answerNo                          = common.AnswerNo
	emptyExerciseNameText             = common.EmptyExerciseNameText
	unknownMuscleGroupText            = common.UnknownMuscleGroupText
	exerciseWithNameAlreadyExistsText = common.ExerciseWithNameAlreadyExistsText
	unknownCommandText                = common.UnknownCommandText

	// error messages
	errStartTraining     = common.ErrStartTraining
	errNoTraining        = common.ErrNoTraining
	errExerciseLoad      = common.ErrExerciseLoad
	errClearTraining     = common.ErrClearTraining
	errUploadTraining    = common.ErrUploadTraining
	errGetTrainings      = common.ErrGetTrainings
	errNoExercises       = common.ErrNoExercises
	errAddExercise       = common.ErrAddExercise
	errProgression       = common.ErrProgression
	errInvalidFormat     = common.ErrInvalidFormat
	errParseData         = common.ErrParseData
	errGeneral           = common.ErrGeneral
	errInvalidExerciseID = common.ErrInvalidExerciseID
	errCreateExercise    = common.ErrCreateExercise
	errInternal          = common.ErrInternal
)

var (
	muscleGroupSmilePrefix = common.MuscleGroupSmilePrefix
	muscleGroupsWithSmiles = common.MuscleGroupsWithSmiles
	availableMuscleGroups  = common.AvailableMuscleGroups
)
