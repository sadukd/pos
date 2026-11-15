package domain

type Score struct{ value int }

func NewScore(value int) (Score, error) {
	if value < 0 {
		return Score{}, ErrInvalidScore
	}
	return Score{value: value}, nil
}
func (s Score) Value() int { return s.value }
func (s Score) Validate(maxScore int) error {
	if maxScore <= 0 || s.value < 0 || s.value > maxScore {
		return ErrInvalidScore
	}
	return nil
}
