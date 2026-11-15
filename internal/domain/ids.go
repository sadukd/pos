package domain

type StudentID int
type CriteriaID int

type StudentIDGenerator interface{ NextStudentID() (StudentID, error) }
type CriteriaIDGenerator interface{ NextCriteriaID() (CriteriaID, error) }

type SequentialIDs struct {
	nextStudent  int
	nextCriteria int
}

func NewSequentialIDs() *SequentialIDs { return &SequentialIDs{} }
func (g *SequentialIDs) NextStudentID() (StudentID, error) {
	if g == nil {
		return 0, ErrInvalidID
	}
	g.nextStudent++
	return StudentID(g.nextStudent), nil
}
func (g *SequentialIDs) NextCriteriaID() (CriteriaID, error) {
	if g == nil {
		return 0, ErrInvalidID
	}
	g.nextCriteria++
	return CriteriaID(g.nextCriteria), nil
}
