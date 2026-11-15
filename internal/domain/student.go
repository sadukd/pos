package domain

import "strings"

type Student struct {
	id    StudentID
	name  FullName
	group string
	email Email
}

func NewStudent(id StudentID, name FullName, group string, email Email) (Student, error) {
	if id <= 0 {
		return Student{}, ErrInvalidID
	}
	if err := name.Validate(); err != nil {
		return Student{}, err
	}
	group = strings.TrimSpace(group)
	if group == "" {
		return Student{}, ErrInvalidGroup
	}
	if err := email.Validate(); err != nil {
		return Student{}, err
	}
	return Student{id: id, name: name, group: group, email: email}, nil
}
func (s Student) ID() StudentID      { return s.id }
func (s Student) FullName() FullName { return s.name }
func (s Student) Group() string      { return s.group }
func (s Student) Email() Email       { return s.email }
func (s *Student) ChangeFullName(name FullName) error {
	if s == nil {
		return ErrStudentNotFound
	}
	if err := name.Validate(); err != nil {
		return err
	}
	s.name = name
	return nil
}
func (s *Student) ChangeGroup(group string) error {
	if s == nil {
		return ErrStudentNotFound
	}
	group = strings.TrimSpace(group)
	if group == "" {
		return ErrInvalidGroup
	}
	s.group = group
	return nil
}
func (s *Student) ChangeEmail(email Email) error {
	if s == nil {
		return ErrStudentNotFound
	}
	if err := email.Validate(); err != nil {
		return err
	}
	s.email = email
	return nil
}
func (s Student) Validate() error {
	if s.id <= 0 {
		return ErrInvalidID
	}
	if err := s.name.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(s.group) == "" {
		return ErrInvalidGroup
	}
	return s.email.Validate()
}
func (s Student) Clone() Student { return s }
