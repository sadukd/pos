package application

import (
	"pos/internal/application/ports"
	"pos/internal/domain"
)

type EvaluationService struct {
	students    ports.StudentRepository
	rubric      ports.RubricRepository
	evaluations ports.StudentEvaluationsRepository
	factory     domain.EvaluationFactory
	calc        *domain.ScoreCalculator
}

func NewEvaluationService(s ports.StudentRepository, r ports.RubricRepository, e ports.StudentEvaluationsRepository, f domain.EvaluationFactory, c *domain.ScoreCalculator) *EvaluationService {
	return &EvaluationService{s, r, e, f, c}
}
func (s *EvaluationService) Set(req EvaluationRequest) (EvaluationResponse, error) {
	studentID := domain.StudentID(req.StudentID)
	if _, err := s.students.Find(studentID); err != nil {
		return EvaluationResponse{}, err
	}
	r, err := s.rubric.Get()
	if err != nil {
		return EvaluationResponse{}, err
	}
	evs, err := s.evaluations.Find(studentID)
	if err != nil {
		return EvaluationResponse{}, err
	}
	score, err := domain.NewScore(req.Score)
	if err != nil {
		return EvaluationResponse{}, err
	}
	ev, err := s.factory.Create(r, domain.CriteriaID(req.CriteriaID), score, req.Comment)
	if err != nil {
		return EvaluationResponse{}, err
	}
	if err = evs.SetEvaluation(ev, r); err != nil {
		return EvaluationResponse{}, err
	}
	if err = s.evaluations.Save(evs); err != nil {
		return EvaluationResponse{}, err
	}
	return EvaluationResponse{req.StudentID, req.CriteriaID, req.Score, ev.Comment()}, nil
}
func (s *EvaluationService) Remove(studentID, criteriaID int) error {
	v, e := s.evaluations.Find(domain.StudentID(studentID))
	if e != nil {
		return e
	}
	if e = v.RemoveEvaluation(domain.CriteriaID(criteriaID)); e != nil {
		return e
	}
	return s.evaluations.Save(v)
}
func (s *EvaluationService) ForStudent(id int) (StudentResult, error) {
	student, e := s.students.Find(domain.StudentID(id))
	if e != nil {
		return StudentResult{}, e
	}
	rubric, e := s.rubric.Get()
	if e != nil {
		return StudentResult{}, e
	}
	evs, e := s.evaluations.Find(domain.StudentID(id))
	if e != nil {
		return StudentResult{}, e
	}
	final, e := s.calc.Calculate(rubric, evs)
	if e != nil {
		return StudentResult{}, e
	}
	out := StudentResult{Student: studentResponse(student), FinalScore: final}
	for _, c := range rubric.Criteria() {
		ev, err := evs.Evaluation(c.ID())
		if err == domain.ErrEvaluationNotFound {
			out.Evaluations = append(out.Evaluations, EvaluationResponse{id, int(c.ID()), -1, ""})
			continue
		}
		if err != nil {
			return StudentResult{}, err
		}
		out.Evaluations = append(out.Evaluations, EvaluationResponse{id, int(c.ID()), ev.Score().Value(), ev.Comment()})
	}
	return out, nil
}
func (s *EvaluationService) AllResults() ([]StudentResult, error) {
	students, e := s.students.List()
	if e != nil {
		return nil, e
	}
	out := make([]StudentResult, 0, len(students))
	for _, v := range students {
		r, e := s.ForStudent(int(v.ID()))
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	return out, nil
}
