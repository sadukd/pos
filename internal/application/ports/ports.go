package ports

import "pos/internal/domain"

type StudentRepository interface {
	Save(domain.Student) error
	Find(domain.StudentID) (domain.Student, error)
	List() ([]domain.Student, error)
	Delete(domain.StudentID) error
}
type RubricRepository interface {
	Save(domain.Rubric) error
	Get() (domain.Rubric, error)
}
type StudentEvaluationsRepository interface {
	Save(domain.StudentEvaluations) error
	Find(domain.StudentID) (domain.StudentEvaluations, error)
	Delete(domain.StudentID) error
	List() ([]domain.StudentEvaluations, error)
}
type ResultExporter interface {
	Export([]domain.Student, domain.Rubric, []domain.StudentEvaluations, map[domain.StudentID]int) error
}
