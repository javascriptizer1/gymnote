package formatter

import (
	"fmt"
	"strings"

	"gymnote/internal/entity"
)

type formatter struct{}

func New() *formatter {
	return &formatter{}
}

func (f *formatter) FormatTrainingLogs(sessions []entity.TrainingSession) string {
	var sb strings.Builder

	for _, session := range sessions {
		sb.WriteString(fmt.Sprintf("%s\n", session.Date().Format("2006-01-02")))

		for _, ex := range session.Exercises() {
			setStrings := []string{}

			for _, set := range ex.Sets() {
				setStr := fmt.Sprintf("%.1f,%d", set.Weight(), set.Reps())
				if set.Notes() != "" {
					setStr += fmt.Sprintf(" (%s)", set.Notes())
				}
				setStrings = append(setStrings, setStr)
			}

			sb.WriteString(fmt.Sprintf("%d. %s - %s\n", ex.Number(), ex.Exercise.Name(), strings.Join(setStrings, "; ")))
		}
		sb.WriteString("\n\n")
	}

	return sb.String()
}
