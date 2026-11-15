package domain

import "strings"

type Evaluation struct {
	criteriaID CriteriaID
	score      Score
	comment    string
}

func NewEvaluation(criteria Criteria, score Score, comment string) (Evaluation, error) {
	if err := criteria.Validate(); err != nil {
		return Evaluation{}, err
	}
	if err := score.Validate(criteria.maxScore); err != nil {
		return Evaluation{}, err
	}
	return Evaluation{criteriaID: criteria.id, score: score, comment: strings.TrimSpace(comment)}, nil
}
func (e Evaluation) CriteriaID() CriteriaID { return e.criteriaID }
func (e Evaluation) Score() Score           { return e.score }
func (e Evaluation) Comment() string        { return e.comment }
func (e Evaluation) Validate(criteria Criteria) error {
	if err := criteria.Validate(); err != nil {
		return err
	}
	if e.criteriaID != criteria.id {
		return ErrCriterionNotInRubric
	}
	return e.score.Validate(criteria.maxScore)
}
