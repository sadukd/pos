package repository

import (
	"pos/internal/domain"
	"sort"
)

type Students struct {
	data map[domain.StudentID]domain.Student
}

func NewStudents() *Students { return &Students{map[domain.StudentID]domain.Student{}} }
func (r *Students) Save(v domain.Student) error {
	if e := v.Validate(); e != nil {
		return e
	}
	if r.data == nil {
		r.data = map[domain.StudentID]domain.Student{}
	}
	r.data[v.ID()] = v.Clone()
	return nil
}
func (r *Students) Find(id domain.StudentID) (domain.Student, error) {
	v, ok := r.data[id]
	if !ok {
		return domain.Student{}, domain.ErrStudentNotFound
	}
	return v.Clone(), nil
}
func (r *Students) List() ([]domain.Student, error) {
	out := []domain.Student{}
	for _, v := range r.data {
		out = append(out, v.Clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out, nil
}
func (r *Students) Delete(id domain.StudentID) error {
	if _, ok := r.data[id]; !ok {
		return domain.ErrStudentNotFound
	}
	delete(r.data, id)
	return nil
}

type Rubrics struct{ value domain.Rubric }

func NewRubrics() *Rubrics                     { return &Rubrics{domain.NewRubric()} }
func (r *Rubrics) Save(v domain.Rubric) error  { r.value = v.Clone(); return nil }
func (r *Rubrics) Get() (domain.Rubric, error) { return r.value.Clone(), nil }

type Evaluations struct {
	data map[domain.StudentID]domain.StudentEvaluations
}

func NewEvaluations() *Evaluations {
	return &Evaluations{map[domain.StudentID]domain.StudentEvaluations{}}
}
func (r *Evaluations) Save(v domain.StudentEvaluations) error {
	if r.data == nil {
		r.data = map[domain.StudentID]domain.StudentEvaluations{}
	}
	r.data[v.StudentID()] = v.Clone()
	return nil
}
func (r *Evaluations) Find(id domain.StudentID) (domain.StudentEvaluations, error) {
	v, ok := r.data[id]
	if !ok {
		return domain.StudentEvaluations{}, domain.ErrStudentNotFound
	}
	return v.Clone(), nil
}
func (r *Evaluations) Delete(id domain.StudentID) error {
	if _, ok := r.data[id]; !ok {
		return domain.ErrStudentNotFound
	}
	delete(r.data, id)
	return nil
}
func (r *Evaluations) List() ([]domain.StudentEvaluations, error) {
	out := []domain.StudentEvaluations{}
	for _, v := range r.data {
		out = append(out, v.Clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StudentID() < out[j].StudentID() })
	return out, nil
}
