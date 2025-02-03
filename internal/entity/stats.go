package entity

import "time"

type ExerciseProgression struct {
	ExerciseName string
	SessionDate  time.Time
	Weight       float32
	Reps         uint8
}
