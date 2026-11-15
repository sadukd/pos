package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"pos/internal/domain"
	"sort"
	"strconv"
)

type CSVExporter struct{ Path string }

func NewCSVExporter() *CSVExporter { return &CSVExporter{Path: "./results.csv"} }
func (x *CSVExporter) Export(students []domain.Student, r domain.Rubric, evals []domain.StudentEvaluations, finals map[domain.StudentID]int) error {
	path := x.Path
	if path == "" {
		path = "./results.csv"
	}
	file, e := os.Create(path)
	if e != nil {
		return e
	}
	defer file.Close()
	w := csv.NewWriter(file)
	criteria := r.Criteria()
	sort.Slice(criteria, func(i, j int) bool { return criteria[i].ID() < criteria[j].ID() })
	sort.Slice(students, func(i, j int) bool { return students[i].ID() < students[j].ID() })
	header := []string{"student_id", "full_name", "group", "email"}
	for _, c := range criteria {
		header = append(header, c.Name())
	}
	header = append(header, "final_score")
	if e = w.Write(header); e != nil {
		return e
	}
	byStudent := map[domain.StudentID]domain.StudentEvaluations{}
	for _, v := range evals {
		byStudent[v.StudentID()] = v
	}
	for _, s := range students {
		row := []string{strconv.Itoa(int(s.ID())), s.FullName().String(), s.Group(), s.Email().String()}
		evs := byStudent[s.ID()]
		for _, c := range criteria {
			ev, err := evs.Evaluation(c.ID())
			if err == domain.ErrEvaluationNotFound {
				row = append(row, "N/A")
			} else if err != nil {
				return err
			} else {
				row = append(row, strconv.Itoa(ev.Score().Value()))
			}
		}
		row = append(row, strconv.Itoa(finals[s.ID()]))
		if e = w.Write(row); e != nil {
			return e
		}
	}
	w.Flush()
	if e = w.Error(); e != nil {
		return fmt.Errorf("write CSV: %w", e)
	}
	return nil
}
