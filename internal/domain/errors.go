package domain

import "errors"

var (
	ErrInvalidID             = errors.New("id must be positive")
	ErrInvalidFullName       = errors.New("invalid full name")
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidGroup          = errors.New("group must not be empty")
	ErrInvalidCriteria       = errors.New("invalid criteria")
	ErrInvalidWeight         = errors.New("weight must be finite and between 0 and 1")
	ErrInvalidMaxScore       = errors.New("maximum score must be positive")
	ErrInvalidCriteriaName   = errors.New("criteria name must not be empty")
	ErrInvalidScore          = errors.New("invalid score")
	ErrStudentNotFound       = errors.New("student not found")
	ErrCriteriaNotFound      = errors.New("criteria not found")
	ErrEvaluationNotFound    = errors.New("evaluation not found")
	ErrDuplicateCriteriaName = errors.New("criteria name already exists")
	ErrDuplicateCriteriaID   = errors.New("criteria id already exists")
	ErrCriteriaInUse         = errors.New("criteria has evaluations; remove them first")
	ErrCriterionNotInRubric  = errors.New("criteria does not belong to rubric")
	ErrInvalidStrategy       = errors.New("score strategy must not be nil")
)
