package domain

import "strings"

type FullName struct{ firstName, lastName, middleName string }

func NewFullName(firstName, lastName, middleName string) (FullName, error) {
	firstName, lastName, middleName = strings.TrimSpace(firstName), strings.TrimSpace(lastName), strings.TrimSpace(middleName)
	if firstName == "" || lastName == "" {
		return FullName{}, ErrInvalidFullName
	}
	return FullName{firstName: firstName, lastName: lastName, middleName: middleName}, nil
}
func (n FullName) FirstName() string  { return n.firstName }
func (n FullName) LastName() string   { return n.lastName }
func (n FullName) MiddleName() string { return n.middleName }
func (n FullName) String() string {
	parts := []string{n.lastName, n.firstName}
	if n.middleName != "" {
		parts = append(parts, n.middleName)
	}
	return strings.Join(parts, " ")
}
func (n FullName) Validate() error {
	if strings.TrimSpace(n.firstName) == "" || strings.TrimSpace(n.lastName) == "" {
		return ErrInvalidFullName
	}
	return nil
}
