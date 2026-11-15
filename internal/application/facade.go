package application

import (
	"pos/internal/application/ports"
	"pos/internal/domain"
)

type POSFacade struct {
	Students    *StudentService
	Rubric      *RubricService
	Evaluations *EvaluationService
	exporter    ports.ResultExporter
	studentRepo ports.StudentRepository
	rubricRepo  ports.RubricRepository
	evalRepo    ports.StudentEvaluationsRepository
}

func NewPOSFacade(s *StudentService, r *RubricService, e *EvaluationService, x ports.ResultExporter, sr ports.StudentRepository, rr ports.RubricRepository, er ports.StudentEvaluationsRepository) *POSFacade {
	return &POSFacade{s, r, e, x, sr, rr, er}
}
func (f *POSFacade) Export() error {
	students, e := f.studentRepo.List()
	if e != nil {
		return e
	}
	rubric, e := f.rubricRepo.Get()
	if e != nil {
		return e
	}
	evaluations, e := f.evalRepo.List()
	if e != nil {
		return e
	}
	finals := map[domain.StudentID]int{}
	for _, s := range students {
		result, e := f.Evaluations.ForStudent(int(s.ID()))
		if e != nil {
			return e
		}
		finals[s.ID()] = result.FinalScore
	}
	return f.exporter.Export(students, rubric, evaluations, finals)
}
