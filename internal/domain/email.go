package domain

import (
	"net/mail"
	"strings"
)

type Email struct{ value string }

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	address, err := mail.ParseAddress(value)
	if err != nil || address.Address != value || !strings.Contains(strings.Split(value, "@")[1], ".") {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: value}, nil
}
func (e Email) String() string { return e.value }
func (e Email) Validate() error {
	if _, err := NewEmail(e.value); err != nil {
		return ErrInvalidEmail
	}
	return nil
}
