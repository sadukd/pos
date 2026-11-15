package domain

type StudentEvaluations struct {
	studentID   StudentID
	evaluations map[CriteriaID]Evaluation
}

func NewStudentEvaluations(studentID StudentID) (StudentEvaluations, error) {
	if studentID <= 0 {
		return StudentEvaluations{}, ErrInvalidID
	}
	return StudentEvaluations{studentID: studentID, evaluations: make(map[CriteriaID]Evaluation)}, nil
}
func (s StudentEvaluations) StudentID() StudentID { return s.studentID }
func (s StudentEvaluations) Evaluations() []Evaluation {
	items := make([]Evaluation, 0, len(s.evaluations))
	for _, e := range s.evaluations {
		items = append(items, e)
	}
	return items
}
func (s StudentEvaluations) Evaluation(id CriteriaID) (Evaluation, error) {
	e, ok := s.evaluations[id]
	if !ok {
		return Evaluation{}, ErrEvaluationNotFound
	}
	return e, nil
}
func (s *StudentEvaluations) SetEvaluation(e Evaluation, rubric Rubric) error {
	if s == nil || s.studentID <= 0 {
		return ErrInvalidID
	}
	c, err := rubric.FindCriteria(e.criteriaID)
	if err != nil {
		return err
	}
	if err := e.Validate(c); err != nil {
		return err
	}
	if s.evaluations == nil {
		s.evaluations = make(map[CriteriaID]Evaluation)
	}
	s.evaluations[e.criteriaID] = e
	return nil
}
func (s *StudentEvaluations) RemoveEvaluation(id CriteriaID) error {
	if s == nil {
		return ErrEvaluationNotFound
	}
	if _, ok := s.evaluations[id]; !ok {
		return ErrEvaluationNotFound
	}
	delete(s.evaluations, id)
	return nil
}
func (s StudentEvaluations) Clone() StudentEvaluations {
	clone, _ := NewStudentEvaluations(s.studentID)
	for id, e := range s.evaluations {
		clone.evaluations[id] = e
	}
	return clone
}
func (s StudentEvaluations) Validate(rubric Rubric) error {
	if s.studentID <= 0 {
		return ErrInvalidID
	}
	for id, e := range s.evaluations {
		c, err := rubric.FindCriteria(id)
		if err != nil {
			return err
		}
		if err := e.Validate(c); err != nil {
			return err
		}
	}
	return nil
}
