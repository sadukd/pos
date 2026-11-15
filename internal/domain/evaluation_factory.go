package domain

type EvaluationFactory struct{}

func NewEvaluationFactory() EvaluationFactory { return EvaluationFactory{} }

func (EvaluationFactory) Create(rubric Rubric, criteriaID CriteriaID, score Score, comment string) (Evaluation, error) {
	criteria, err := rubric.FindCriteria(criteriaID)
	if err != nil {
		return Evaluation{}, err
	}
	return NewEvaluation(criteria, score, comment)
}
