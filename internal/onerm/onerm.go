package onerm

import "math"

type Formula string

const (
	FormulaEpley    Formula = "epley"
	FormulaBrzycki  Formula = "brzycki"
	FormulaLander   Formula = "lander"
	FormulaLombardi Formula = "lombardi"
	FormulaMayhew   Formula = "mayhew"
	FormulaOConner  Formula = "oconner"
	FormulaWathan   Formula = "wathan"
)

type Result struct {
	Formula Formula
	Value   float64
}

type Summary struct {
	Results []Result
	Average float64
}

func Calculate(weight float64, reps int) Summary {
	if weight <= 0 || reps <= 0 {
		return Summary{}
	}

	w := weight
	r := float64(reps)

	results := make([]Result, 0, 7)

	// Эпли: 1RM = w * (1 + r/30)
	epley := w * (1 + r/30)
	results = append(results, Result{Formula: FormulaEpley, Value: epley})

	// Бжицки: 1RM = w * 36 / (37 - r)
	if d := 37 - r; d > 0 {
		brzycki := w * 36 / d
		results = append(results, Result{Formula: FormulaBrzycki, Value: brzycki})
	}

	// Лэндера: 1RM = w * 100 / (101.3 - 2.67123*r)
	if d := 101.3 - 2.67123*r; d > 0 {
		lander := w * 100 / d
		results = append(results, Result{Formula: FormulaLander, Value: lander})
	}

	// Ломбарди: 1RM = w * r^0.10
	lombardi := w * math.Pow(r, 0.10)
	results = append(results, Result{Formula: FormulaLombardi, Value: lombardi})

	// Мэйхью: 1RM = (100 * w) / (52.2 + 41.9*e^{-0.055*r})
	mayhew := (100 * w) / (52.2 + 41.9*math.Exp(-0.055*r))
	results = append(results, Result{Formula: FormulaMayhew, Value: mayhew})

	// О'Коннор: 1RM = w * (1 + 0.025*r)
	oconner := w * (1 + 0.025*r)
	results = append(results, Result{Formula: FormulaOConner, Value: oconner})

	// Ватан: 1RM = (100 * w) / (48.8 + 53.8*e^{-0.075*r})
	wathan := (100 * w) / (48.8 + 53.8*math.Exp(-0.075*r))
	results = append(results, Result{Formula: FormulaWathan, Value: wathan})

	var sum float64
	for _, r := range results {
		sum += r.Value
	}

	avg := 0.0
	if len(results) > 0 {
		avg = sum / float64(len(results))
	}

	return Summary{
		Results: results,
		Average: avg,
	}
}
