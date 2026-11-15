package domain

import (
	"errors"
	"math"
	"testing"
)

func mustCriteria(t *testing.T, id CriteriaID, name string, max int, weight float64) Criteria {
	t.Helper()
	c, err := NewCriteria(id, name, max, weight)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestValueObjectsAndStudentInvariants(t *testing.T) {
	name, err := NewFullName(" Ada ", "Lovelace ", " Byron ")
	if err != nil || name.String() != "Lovelace Ada Byron" {
		t.Fatalf("name=%q err=%v", name.String(), err)
	}
	if _, err := NewFullName(" ", "Lovelace", ""); !errors.Is(err, ErrInvalidFullName) {
		t.Fatalf("expected invalid full name, got %v", err)
	}
	email, err := NewEmail("ada@example.org")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewEmail("not-an-email"); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected invalid email, got %v", err)
	}
	s, err := NewStudent(1, name, " CS-101 ", email)
	if err != nil {
		t.Fatal(err)
	}
	if s.Group() != "CS-101" {
		t.Fatalf("group not normalized: %q", s.Group())
	}
	if err := s.ChangeGroup(" "); !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("expected invalid group, got %v", err)
	}
	if err := s.ChangeEmail(Email{}); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("expected invalid email, got %v", err)
	}
	if err := (Student{}).Validate(); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("zero student should be invalid: %v", err)
	}
}

func TestCriteriaValidationAndRubricUniqueness(t *testing.T) {
	if _, err := NewCriteria(1, "x", 0, .5); !errors.Is(err, ErrInvalidMaxScore) {
		t.Fatalf("max error: %v", err)
	}
	for _, weight := range []float64{-0.1, 1.1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NewCriteria(1, "x", 10, weight); !errors.Is(err, ErrInvalidWeight) {
			t.Errorf("weight %v error: %v", weight, err)
		}
	}
	c := mustCriteria(t, 1, "  Participation ", 100, 0)
	r := NewRubric()
	if err := r.AddCriteria(c); err != nil {
		t.Fatal(err)
	}
	if err := r.AddCriteria(mustCriteria(t, 2, "participation", 10, 1)); !errors.Is(err, ErrDuplicateCriteriaName) {
		t.Fatalf("expected duplicate name, got %v", err)
	}
	if err := r.AddCriteria(mustCriteria(t, 1, "Other", 10, 1)); !errors.Is(err, ErrDuplicateCriteriaID) {
		t.Fatalf("expected duplicate ID, got %v", err)
	}
	if _, err := r.FindCriteria(99); !errors.Is(err, ErrCriteriaNotFound) {
		t.Fatalf("missing criteria: %v", err)
	}
	if err := r.RemoveCriteria(99); !errors.Is(err, ErrCriteriaNotFound) {
		t.Fatalf("remove missing criteria: %v", err)
	}
}

func TestEvaluationsReplaceAndEnforceRubricMembership(t *testing.T) {
	c := mustCriteria(t, 7, "Exam", 20, 0.5)
	r := NewRubric()
	if err := r.AddCriteria(c); err != nil {
		t.Fatal(err)
	}
	se, _ := NewStudentEvaluations(3)
	score, _ := NewScore(0)
	e, err := NewEvaluation(c, score, " first ")
	if err != nil {
		t.Fatal(err)
	}
	if err := se.SetEvaluation(e, r); err != nil {
		t.Fatal(err)
	}
	if got, err := se.Evaluation(c.ID()); err != nil || got.Score().Value() != 0 || got.Comment() != "first" {
		t.Fatalf("got=%+v err=%v", got, err)
	}
	score, _ = NewScore(12)
	e, _ = NewEvaluation(c, score, "updated")
	if err := se.SetEvaluation(e, r); err != nil {
		t.Fatal(err)
	}
	got, _ := se.Evaluation(c.ID())
	if got.Score().Value() != 12 || got.Comment() != "updated" {
		t.Fatalf("upsert failed: %+v", got)
	}
	other := mustCriteria(t, 8, "Other", 20, 0.5)
	e, _ = NewEvaluation(other, score, "")
	if err := se.SetEvaluation(e, r); !errors.Is(err, ErrCriteriaNotFound) {
		t.Fatalf("non-member criteria accepted: %v", err)
	}
	tooHigh, _ := NewScore(21)
	if _, err := NewEvaluation(c, tooHigh, ""); !errors.Is(err, ErrInvalidScore) {
		t.Fatalf("score above max accepted: %v", err)
	}
	if err := se.RemoveEvaluation(c.ID()); err != nil {
		t.Fatal(err)
	}
	if err := se.RemoveEvaluation(c.ID()); !errors.Is(err, ErrEvaluationNotFound) {
		t.Fatalf("missing eval: %v", err)
	}
}

type fixedStrategy int

func (s fixedStrategy) Calculate(Rubric, StudentEvaluations) (int, error) { return int(s), nil }

func TestWeightedScoreStrategyAndStrategyInjection(t *testing.T) {
	r := NewRubric()
	c1 := mustCriteria(t, 1, "A", 20, .25)
	c2 := mustCriteria(t, 2, "B", 10, .5)
	if err := r.AddCriteria(c1); err != nil {
		t.Fatal(err)
	}
	if err := r.AddCriteria(c2); err != nil {
		t.Fatal(err)
	}
	se, _ := NewStudentEvaluations(1)
	zero, _ := NewScore(0)
	e, _ := NewEvaluation(c1, zero, "")
	if err := se.SetEvaluation(e, r); err != nil {
		t.Fatal(err)
	}

	if len(se.Evaluations()) != 1 {
		t.Fatalf("missing score was materialized: %d", len(se.Evaluations()))
	}
	calc, err := NewScoreCalculator(WeightedScoreStrategy{})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := calc.Calculate(r, se); err != nil || got != 0 {
		t.Fatalf("got %d, err %v", got, err)
	}
	se, _ = NewStudentEvaluations(1)
	score, _ := NewScore(6)
	e, _ = NewEvaluation(c1, score, "")
	_ = se.SetEvaluation(e, r)
	score, _ = NewScore(10)
	e, _ = NewEvaluation(c2, score, "")
	_ = se.SetEvaluation(e, r)
	if got, _ := calc.Calculate(r, se); got != 7 {
		t.Fatalf("rounding result = %d, want 7", got)
	}
	if got, err := NewScoreCalculator(fixedStrategy(4)); err != nil {
		t.Fatal(err)
	} else if result, _ := got.Calculate(r, se); result != 4 {
		t.Fatalf("custom strategy result = %d", result)
	}

	r2 := NewRubric()
	c3 := mustCriteria(t, 3, "C", 100, 1)
	_ = r2.AddCriteria(c3)
	se2, _ := NewStudentEvaluations(2)
	max, _ := NewScore(100)
	e, _ = NewEvaluation(c3, max, "")
	_ = se2.SetEvaluation(e, r2)
	if got, _ := calc.Calculate(r2, se2); got != 10 {
		t.Fatalf("capped result = %d, want 10", got)
	}
}

func TestSequentialIDsUseIndependentSequences(t *testing.T) {
	g := NewSequentialIDs()
	s1, _ := g.NextStudentID()
	c1, _ := g.NextCriteriaID()
	s2, _ := g.NextStudentID()
	c2, _ := g.NextCriteriaID()
	if s1 != 1 || c1 != 1 || s2 != 2 || c2 != 2 {
		t.Fatalf("unexpected sequences: %v %v %v %v", s1, c1, s2, c2)
	}
}
