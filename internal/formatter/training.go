package formatter

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"gymnote/internal/entity"
)

type formatter struct{}

func New() *formatter {
	return &formatter{}
}

func (f *formatter) FormatTrainingLogs(sessions []entity.TrainingSession, showDuration bool) string {
	var sb strings.Builder

	for _, session := range sessions {
		if showDuration && session.Duration() > 0 {
			fmt.Fprintf(&sb, "%s (%s)\n", session.Date().Format(time.DateOnly), FormatDuration(session.Duration()))
		} else {
			fmt.Fprintf(&sb, "%s\n", session.Date().Format(time.DateOnly))
		}

		for _, ex := range session.Exercises() {
			setStrings := []string{}

			for _, set := range ex.Sets() {
				setStr := fmt.Sprintf("%s,%d", formatWeight(set.Weight()), set.Reps())
				if set.Notes() != "" {
					setStr += fmt.Sprintf(" (%s)", set.Notes())
				}
				setStrings = append(setStrings, setStr)
			}

			fmt.Fprintf(&sb, "%d. %s - %s\n", ex.Number(), ex.Exercise.Name(), strings.Join(setStrings, "; "))
		}

		sb.WriteString("\n\n")
	}

	return sb.String()
}

func FormatWeightFloat(v float64) string {
	if math.Mod(v, 1) == 0 {
		return fmt.Sprintf("%.0f", v)
	}

	return fmt.Sprintf("%.1f", v)
}

func formatWeight(weight float32) string {
	return FormatWeightFloat(float64(weight))
}

func FormatDuration(d time.Duration) string {
	if d <= 0 {
		return "0 мин."
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d ч. %d мин.", hours, minutes)
	}

	return fmt.Sprintf("%d мин.", minutes)
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
			setStrings = append(setStrings, fmt.Sprintf("%s x %d", formatWeight(set.Weight), set.Reps))
		}

		sb.WriteString(strings.Join(setStrings, "; ") + "\n\n")
	}

	return sb.String()
}
