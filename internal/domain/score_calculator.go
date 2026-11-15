package domain

import (
	"errors"
	"math"
)

type ScoreStrategy interface {
	Calculate(rubric Rubric, evaluations StudentEvaluations) (int, error)
}
type WeightedScoreStrategy struct{}

func (WeightedScoreStrategy) Calculate(rubric Rubric, evaluations StudentEvaluations) (int, error) {
	if err := evaluations.Validate(rubric); err != nil {
		return 0, err
	}
	total := 0.0
	for _, c := range rubric.Criteria() {
		e, err := evaluations.Evaluation(c.id)
		if errors.Is(err, ErrEvaluationNotFound) {
			continue
		}
		if err != nil {
			return 0, err
		}
		total += float64(e.score.value) * c.weight
	}
	return int(math.Round(math.Min(total, 10))), nil
}

type ScoreCalculator struct{ strategy ScoreStrategy }

func NewScoreCalculator(strategy ScoreStrategy) (*ScoreCalculator, error) {
	if strategy == nil {
		return nil, ErrInvalidStrategy
	}
	return &ScoreCalculator{strategy: strategy}, nil
}
func (c *ScoreCalculator) Calculate(rubric Rubric, evaluations StudentEvaluations) (int, error) {
	if c == nil || c.strategy == nil {
		return 0, ErrInvalidStrategy
	}
	return c.strategy.Calculate(rubric, evaluations)
}
