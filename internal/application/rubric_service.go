package application

import (
	"pos/internal/application/ports"
	"pos/internal/domain"
)

type RubricService struct {
	repo        ports.RubricRepository
	evaluations ports.StudentEvaluationsRepository
	ids         domain.CriteriaIDGenerator
}

func NewRubricService(r ports.RubricRepository, ev ports.StudentEvaluationsRepository, ids domain.CriteriaIDGenerator) *RubricService {
	return &RubricService{r, ev, ids}
}
func (s *RubricService) Add(req CriteriaRequest) (CriteriaResponse, error) {
	id, e := s.ids.NextCriteriaID()
	if e != nil {
		return CriteriaResponse{}, e
	}
	c, e := domain.NewCriteria(id, req.Name, req.MaxScore, req.Weight)
	if e != nil {
		return CriteriaResponse{}, e
	}
	r, e := s.repo.Get()
	if e != nil {
		return CriteriaResponse{}, e
	}
	if e = r.AddCriteria(c); e != nil {
		return CriteriaResponse{}, e
	}
	if e = s.repo.Save(r); e != nil {
		return CriteriaResponse{}, e
	}
	return CriteriaResponse{int(id), c.Name(), c.MaxScore(), c.Weight()}, nil
}
func (s *RubricService) List() ([]CriteriaResponse, error) {
	r, e := s.repo.Get()
	if e != nil {
		return nil, e
	}
	out := []CriteriaResponse{}
	for _, c := range r.Criteria() {
		out = append(out, CriteriaResponse{int(c.ID()), c.Name(), c.MaxScore(), c.Weight()})
	}
	return out, nil
}
func (s *RubricService) Remove(id int) error {
	r, e := s.repo.Get()
	if e != nil {
		return e
	}
	if _, e = r.FindCriteria(domain.CriteriaID(id)); e != nil {
		return e
	}
	all, e := s.evaluations.List()
	if e != nil {
		return e
	}
	for _, ev := range all {
		if _, e = ev.Evaluation(domain.CriteriaID(id)); e == nil {
			return domain.ErrCriteriaInUse
		}
	}
	if e = r.RemoveCriteria(domain.CriteriaID(id)); e != nil {
		return e
	}
	return s.repo.Save(r)
}
