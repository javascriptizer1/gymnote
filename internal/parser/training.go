package parser

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gymnote/internal/helper"
)

const (
	DifficultyEasy   = "легко"
	DifficultyMedium = "средне"
	DifficultyHard   = "тяжело"
)

type Exercise struct {
	Name string
	Sets []Set
}

type Set struct {
	Weight     float32
	Reps       uint8
	Difficulty string
	Notes      string
}

type parser struct{}

func New() *parser {
	return &parser{}
}

// Пример текста
// 2024-02-15
// 1. Жим в Хаммере - 40,12 (легко); 40,12; 40,12
// 2. Жим лежа - 50,10 (легко); 50,10; 95,1
// 3. Жим на наклонной скамье - 50,10 (легко);
// 4. Жим на наклонной скамье -  50,10 (легко); 80,1 (с помощью); 60,5 (хорошо)
// 5. Жим гантелей лежа - 25,10 (нормально); 25,10 (нормально)
// 6. Разводки гантелей лежа - 15,12 (средне)
// 7. Разгибание в блоке на трицепс - 42,12 (легко); 50,12 (легко); 50,12 (на коленях, средне)

func (p *parser) ParseExercises(s string) ([]Exercise, time.Time, error) {
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return nil, time.Time{}, errors.New("no data to process")
	}

	var date time.Time
	firstLine := strings.TrimSpace(lines[0])
	parsedTime, ok := p.isValidDate(firstLine)
	if ok {
		date = parsedTime
		lines = lines[1:]
	} else {
		date = time.Now()
	}

	var exercises []Exercise

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		exs, err := p.parseExercise(line)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("parse exercise error: %w", err)
		}

		exercises = append(exercises, exs)
	}

	return exercises, date, nil
}

func (p *parser) parseExercise(line string) (Exercise, error) {
	exs := Exercise{}
	parts := strings.SplitN(line, " - ", 2)
	if len(parts) != 2 {
		return exs, errors.New("invalid exercise format")
	}

	exerciseName := strings.TrimSpace(strings.Split(parts[0], ".")[1])
	setsData := strings.Split(parts[1], ";")

	var sets []Set
	for _, setData := range setsData {
		if setData == "" {
			continue
		}
		set, err := p.parseSet(strings.TrimSpace(setData))
		if err != nil {
			return exs, err
		}
		sets = append(sets, set)
	}

	exs.Name = exerciseName
	exs.Sets = sets

	return exs, nil
}

func (p *parser) parseSet(setData string) (Set, error) {
	set := Set{}

	if strings.Contains(setData, "(") {
		start := strings.Index(setData, "(")
		end := strings.Index(setData, ")")
		if start == -1 || end == -1 || end <= start {
			return set, errors.New("invalid set notes format")
		}
		set.Notes = strings.TrimSpace(setData[start+1 : end])
		setData = strings.TrimSpace(setData[:start])
	}

	fields := strings.Split(setData, ",")
	if len(fields) == 1 {
		reps, err := helper.ParseUint8(strings.TrimSpace(fields[0]))
		if err != nil {
			return set, fmt.Errorf("invalid reps format: %w", err)
		}
		set.Weight = 1
		set.Reps = reps
	} else if len(fields) == 2 {
		weight, err := helper.ParseFloat32(strings.TrimSpace(fields[0]))
		if err != nil {
			return set, fmt.Errorf("invalid weight format: %w", err)
		}

		reps, err := helper.ParseUint8(strings.TrimSpace(fields[1]))
		if err != nil {
			return set, fmt.Errorf("invalid reps format: %w", err)
		}

		set.Weight = weight
		set.Reps = reps
	} else {
		return set, errors.New("invalid set format")
	}

	if strings.Contains(set.Notes, DifficultyEasy) {
		set.Difficulty = DifficultyEasy
	} else if strings.Contains(set.Notes, DifficultyMedium) {
		set.Difficulty = DifficultyMedium
	} else if strings.Contains(set.Notes, DifficultyHard) {
		set.Difficulty = DifficultyHard
	} else {
		set.Difficulty = "-"
	}

	return set, nil
}

func (p *parser) isValidDate(dateStr string) (time.Time, bool) {
	value, err := time.Parse("2006-01-02", dateStr)
	return value, err == nil
}
