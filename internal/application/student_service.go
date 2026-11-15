package application

import (
	"fmt"
	"pos/internal/application/ports"
	"pos/internal/domain"
)

type StudentService struct {
	students    ports.StudentRepository
	evaluations ports.StudentEvaluationsRepository
	ids         domain.StudentIDGenerator
}

func NewStudentService(s ports.StudentRepository, e ports.StudentEvaluationsRepository, ids domain.StudentIDGenerator) *StudentService {
	return &StudentService{s, e, ids}
}
func studentResponse(s domain.Student) StudentResponse {
	return StudentResponse{int(s.ID()), s.FullName().String(), s.Group(), s.Email().String()}
}
func (s *StudentService) Create(r StudentRequest) (StudentResponse, error) {
	id, err := s.ids.NextStudentID()
	if err != nil {
		return StudentResponse{}, err
	}
	n, err := domain.NewFullName(r.FirstName, r.LastName, r.MiddleName)
	if err != nil {
		return StudentResponse{}, err
	}
	e, err := domain.NewEmail(r.Email)
	if err != nil {
		return StudentResponse{}, err
	}
	v, err := domain.NewStudent(id, n, r.Group, e)
	if err != nil {
		return StudentResponse{}, err
	}
	ev, _ := domain.NewStudentEvaluations(id)
	if err = s.students.Save(v); err != nil {
		return StudentResponse{}, err
	}
	if err = s.evaluations.Save(ev); err != nil {
		return StudentResponse{}, fmt.Errorf("create student evaluations: %w", err)
	}
	return studentResponse(v), nil
}
func (s *StudentService) List() ([]StudentResponse, error) {
	all, err := s.students.List()
	if err != nil {
		return nil, err
	}
	out := make([]StudentResponse, 0, len(all))
	for _, v := range all {
		out = append(out, studentResponse(v))
	}
	return out, nil
}
func (s *StudentService) Get(id int) (StudentResponse, error) {
	v, e := s.students.Find(domain.StudentID(id))
	if e != nil {
		return StudentResponse{}, e
	}
	return studentResponse(v), nil
}
func (s *StudentService) Update(id int, r StudentRequest) (StudentResponse, error) {
	v, e := s.students.Find(domain.StudentID(id))
	if e != nil {
		return StudentResponse{}, e
	}
	n, e := domain.NewFullName(r.FirstName, r.LastName, r.MiddleName)
	if e != nil {
		return StudentResponse{}, e
	}
	mail, e := domain.NewEmail(r.Email)
	if e != nil {
		return StudentResponse{}, e
	}
	if e = v.ChangeFullName(n); e != nil {
		return StudentResponse{}, e
	}
	if e = v.ChangeGroup(r.Group); e != nil {
		return StudentResponse{}, e
	}
	if e = v.ChangeEmail(mail); e != nil {
		return StudentResponse{}, e
	}
	if e = s.students.Save(v); e != nil {
		return StudentResponse{}, e
	}
	return studentResponse(v), nil
}
func (s *StudentService) Delete(id int) error {
	key := domain.StudentID(id)
	if _, e := s.students.Find(key); e != nil {
		return e
	}
	if e := s.students.Delete(key); e != nil {
		return e
	}
	return s.evaluations.Delete(key)
}
