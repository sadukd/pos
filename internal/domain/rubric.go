package domain

type Rubric struct{ criteria map[CriteriaID]Criteria }

func NewRubric() Rubric { return Rubric{criteria: make(map[CriteriaID]Criteria)} }
func (r Rubric) Criteria() []Criteria {
	items := make([]Criteria, 0, len(r.criteria))
	for _, c := range r.criteria {
		items = append(items, c)
	}
	return items
}
func (r Rubric) FindCriteria(id CriteriaID) (Criteria, error) {
	c, ok := r.criteria[id]
	if !ok {
		return Criteria{}, ErrCriteriaNotFound
	}
	return c, nil
}
func (r *Rubric) AddCriteria(c Criteria) error {
	if r == nil {
		return ErrInvalidCriteria
	}
	if err := c.Validate(); err != nil {
		return err
	}
	if r.criteria == nil {
		r.criteria = make(map[CriteriaID]Criteria)
	}
	if _, exists := r.criteria[c.id]; exists {
		return ErrDuplicateCriteriaID
	}
	for _, existing := range r.criteria {
		if normalizeCriteriaName(existing.name) == normalizeCriteriaName(c.name) {
			return ErrDuplicateCriteriaName
		}
	}
	r.criteria[c.id] = c
	return nil
}
func (r *Rubric) RemoveCriteria(id CriteriaID) error {
	if r == nil {
		return ErrCriteriaNotFound
	}
	if _, ok := r.criteria[id]; !ok {
		return ErrCriteriaNotFound
	}
	delete(r.criteria, id)
	return nil
}
func (r Rubric) Contains(c Criteria) bool {
	if err := c.Validate(); err != nil {
		return false
	}
	current, ok := r.criteria[c.id]
	return ok && current == c
}
func (r Rubric) Clone() Rubric {
	clone := NewRubric()
	for id, c := range r.criteria {
		clone.criteria[id] = c.Clone()
	}
	return clone
}
