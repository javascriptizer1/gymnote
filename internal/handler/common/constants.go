package common

const (
	StartCommand                  = "start"
	HelpCommand                   = "help"
	StartTrainingCommand          = "start_training"
	ClearTrainingCommand          = "clear_training"
	CreateExerciseCommand         = "create_exercise"
	UploadTrainingCommand         = "upload_training"
	GetTrainingsCommand           = "get_trainings"
	GetExerciseProgressionCommand = "get_exercise_progression"
	GetExerciseHistoryCommand     = "get_exercise_history"
	OneRMCommand                  = "one_rm"
)

const (
	MusclePrefix                      = "muscle:"
	ExercisePrefix                    = "exercise:"
	ConfirmationFinishTrainingPrefix  = "confirmation_finish_training:"
	AcceptFinishTrainingPrefix        = "accept_finish_training:"
	RejectFinishTrainingPrefix        = "reject_finish_training:"
	StartNewExercisePrefix            = "start_new_exercise:"
	StartGetExerciseProgressionPrefix = "start_progression:"
	StartGetExerciseHistoryPrefix     = "start_exercise_history:"
	BackToMuscleGroups                = "back_to_muscle_groups"
)

const (
	NextDirection = "next"
	PrevDirection = "prev"
)

var (
	MuscleGroupSmilePrefix = "💪🏋️🦵🤷 "
	MuscleGroupsWithSmiles = []string{"💪 Спина", "🏋️ Грудь", "🦵 Ноги", "💪 Руки", "🤷 Плечи"}
	AvailableMuscleGroups  = []string{"Спина", "Грудь", "Ноги", "Руки", "Плечи"}
)

const (
	DaysForSetStatistics   = int64(3)
	DaysForExerciseHistory = int64(20)
)
