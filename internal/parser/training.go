package parser

import (
	"errors"
	"fmt"
	"gymnote/internal/helper"
	"strings"
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

func (p *parser) ParseExercises(s string) ([]Exercise, error) {
	lines := strings.Split(s, "\n")
	if len(lines) == 0 {
		return nil, errors.New("no data to process")
	}

	var exercises []Exercise

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		exs, err := p.parseExercise(line)
		if err != nil {
			return nil, fmt.Errorf("parse exercise error: %w", err)
		}

		exercises = append(exercises, exs)
	}

	return exercises, nil
}

func (p *parser) parseExercise(line string) (Exercise, error) {
	exs := Exercise{}
	parts := strings.SplitN(line, "-", 2)
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
	if len(fields) != 2 {
		return set, errors.New("invalid set format")
	}

	var err error
	set.Weight, err = helper.ParseFloat32(fields[0])
	if err != nil {
		return set, err
	}

	set.Reps, err = helper.ParseUint8(fields[1])
	if err != nil {
		return set, err
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
