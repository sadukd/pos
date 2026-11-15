package domain

import (
	"math"
	"strings"
)

type Criteria struct {
	id       CriteriaID
	name     string
	maxScore int
	weight   float64
}

func NewCriteria(id CriteriaID, name string, maxScore int, weight float64) (Criteria, error) {
	name = strings.TrimSpace(name)
	if id <= 0 {
		return Criteria{}, ErrInvalidID
	}
	if name == "" {
		return Criteria{}, ErrInvalidCriteriaName
	}
	if maxScore <= 0 {
		return Criteria{}, ErrInvalidMaxScore
	}
	if math.IsNaN(weight) || math.IsInf(weight, 0) || weight < 0 || weight > 1 {
		return Criteria{}, ErrInvalidWeight
	}
	return Criteria{id: id, name: name, maxScore: maxScore, weight: weight}, nil
}
func (c Criteria) ID() CriteriaID  { return c.id }
func (c Criteria) Name() string    { return c.name }
func (c Criteria) MaxScore() int   { return c.maxScore }
func (c Criteria) Weight() float64 { return c.weight }
func (c Criteria) Validate() error {
	_, err := NewCriteria(c.id, c.name, c.maxScore, c.weight)
	return err
}
func (c Criteria) Clone() Criteria { return c }

func normalizeCriteriaName(name string) string { return strings.ToLower(strings.TrimSpace(name)) }
