package entity

type UserState string

const (
	StateNone                        UserState = "none"
	StateAwaitingSetInput            UserState = "awaiting_set_input"
	StateAwaitingExerciseInput       UserState = "awaiting_create_exercise_input"
	StateAwaitingTrainingInput       UserState = "awaiting_upload_training_input"
	StateAwaitingGetTrainingsInput   UserState = "awaiting_get_trainings_input"
	StateAwaitingExerciseProgression UserState = "awaiting_exercise_progression"
)
