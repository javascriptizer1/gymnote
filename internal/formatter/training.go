package formatter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"gymnote/internal/entity"
)

type formatter struct{}

func New() *formatter {
	return &formatter{}
}

func (f *formatter) FormatTrainingLogs(sessions []entity.TrainingSession) string {
	var sb strings.Builder

	for _, session := range sessions {
		sb.WriteString(fmt.Sprintf("%s\n", session.Date().Format(time.DateOnly)))

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

func (f *formatter) FormatLastSets(sets []entity.ExerciseProgression) string {
	var sb strings.Builder

	if len(sets) == 0 {
		return ""
	}

	grouped := make(map[string][]entity.ExerciseProgression)
	dates := []string{}

	for _, set := range sets {
		dateStr := set.SessionDate.Format(time.DateOnly)
		if _, exists := grouped[dateStr]; !exists {
			dates = append(dates, dateStr)
		}
		grouped[dateStr] = append(grouped[dateStr], set)
	}

	sort.Strings(dates)

	for _, date := range dates {
		sb.WriteString(fmt.Sprintf("%s\n", date))

		setStrings := []string{}
		for _, set := range grouped[date] {
			setStrings = append(setStrings, fmt.Sprintf("%.1f кг x %d", set.Weight, set.Reps))
		}

		sb.WriteString(strings.Join(setStrings, "; ") + "\n\n")
	}

	return sb.String()
}
