package entity

type UserState string

const (
	StateNone                  UserState = "none"
	StateAwaitingSetInput      UserState = "awaiting_set_input"
	StateAwaitingExerciseInput UserState = "awaiting_create_exercise_input"
)
