package application

type StudentRequest struct{ FirstName, LastName, MiddleName, Group, Email string }
type StudentResponse struct {
	ID                     int
	FullName, Group, Email string
}
type CriteriaRequest struct {
	Name     string
	MaxScore int
	Weight   float64
}
type CriteriaResponse struct {
	ID       int
	Name     string
	MaxScore int
	Weight   float64
}
type EvaluationRequest struct {
	StudentID  int
	CriteriaID int
	Score      int
	Comment    string
}
type EvaluationResponse struct {
	StudentID  int
	CriteriaID int
	Score      int
	Comment    string
}
type StudentResult struct {
	Student     StudentResponse
	Evaluations []EvaluationResponse
	FinalScore  int
}
