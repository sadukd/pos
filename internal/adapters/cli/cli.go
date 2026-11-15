package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"pos/internal/application"
	"strconv"
	"strings"
	"unicode/utf8"
)

const pageSize = 10
const menuEdge = "≡≡≡≡≡≡"

type menuOption struct {
	label  string
	action func() error
}
type CLI struct {
	app *application.POSFacade
	in  *bufio.Scanner
	out io.Writer
}

func New(app *application.POSFacade, r io.Reader, w io.Writer) *CLI {
	return &CLI{app, bufio.NewScanner(r), w}
}

func (c *CLI) Run() error {
	err := c.runMenu("ПОС", []menuOption{{"Students", c.studentsMenu}, {"Criteria", c.criteriaMenu}, {"Evaluations", c.evaluationsMenu}, {"Results", c.resultsMenu}, {"Export CSV", c.export}}, "Exit")
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}
func (c *CLI) runMenu(title string, opts []menuOption, quit string) error {
	page := 0
	for {
		start := page * pageSize
		end := min(start+pageSize, len(opts))
		if err := c.title(title); err != nil {
			return err
		}
		for i, o := range opts[start:end] {
			if _, err := fmt.Fprintf(c.out, "%s. %s\n", menuKey(i), o.label); err != nil {
				return err
			}
		}
		if err := c.end(title); err != nil {
			return err
		}
		prompt := "q. " + quit + "\n> "
		if len(opts) > pageSize {
			prompt = "< previous | > next | q " + strings.ToLower(quit) + "\n> "
		}
		s, e := c.read(prompt)
		if e != nil {
			return e
		}
		switch s {
		case "q":
			return nil
		case "<":
			if page > 0 {
				page--
			} else {
				fmt.Fprintln(c.out, "No previous page.")
			}
		case ">":
			if end < len(opts) {
				page++
			} else {
				fmt.Fprintln(c.out, "No next page.")
			}
		default:
			i, ok := menuIndex(s, end-start)
			if !ok {
				fmt.Fprintln(c.out, "Invalid choice.")
				continue
			}
			if e := opts[start+i].action(); e != nil {
				return e
			}
		}
	}
}
func (c *CLI) studentsMenu() error {
	return c.runMenu("Students", []menuOption{{"List students", c.listStudents}, {"Add student", c.addStudent}, {"Edit student", c.editStudent}, {"Delete student", c.deleteStudent}}, "Back")
}
func (c *CLI) criteriaMenu() error {
	return c.runMenu("Criteria", []menuOption{{"List criteria", c.listCriteria}, {"Add criterion", c.addCriteria}, {"Delete criterion", c.deleteCriteria}}, "Back")
}
func (c *CLI) evaluationsMenu() error {
	return c.runMenu("Evaluations", []menuOption{{"Set or replace evaluation", c.setEvaluation}, {"Delete evaluation", c.deleteEvaluation}}, "Back")
}
func (c *CLI) resultsMenu() error {
	return c.runMenu("Results", []menuOption{{"Student result", c.studentResult}, {"All results", c.allResults}}, "Back")
}
func (c *CLI) listStudents() error {
	v, e := c.app.Students.List()
	if e != nil {
		return c.showError(e)
	}
	rows := make([]string, len(v))
	for i, s := range v {
		rows[i] = fmt.Sprintf("#%d %s | %s | %s", s.ID, s.FullName, s.Group, s.Email)
	}
	return c.showPage("Students", rows)
}
func (c *CLI) listCriteria() error {
	v, e := c.app.Rubric.List()
	if e != nil {
		return c.showError(e)
	}
	rows := make([]string, len(v))
	for i, x := range v {
		rows[i] = fmt.Sprintf("#%d %s (max %d, weight %g)", x.ID, x.Name, x.MaxScore, x.Weight)
	}
	return c.showPage("Criteria", rows)
}
func (c *CLI) addStudent() error {
	r, e := c.studentFields()
	if e != nil {
		return e
	}
	v, e := c.app.Students.Create(r)
	if e != nil {
		return c.showError(e)
	}
	return c.message(fmt.Sprintf("Added student #%d: %s", v.ID, v.FullName))
}
func (c *CLI) editStudent() error {
	id, e := c.readInt("Student ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	old, e := c.app.Students.Get(id)
	if e != nil {
		return c.showError(e)
	}
	fmt.Fprintf(c.out, "Selected: #%d %s | %s | %s\n", old.ID, old.FullName, old.Group, old.Email)
	r, e := c.studentFields()
	if e != nil {
		return e
	}
	v, e := c.app.Students.Update(id, r)
	if e != nil {
		return c.showError(e)
	}
	return c.message("Updated: " + v.FullName)
}
func (c *CLI) deleteStudent() error {
	id, e := c.readInt("Student ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	v, e := c.app.Students.Get(id)
	if e != nil {
		return c.showError(e)
	}
	fmt.Fprintf(c.out, "Selected: #%d %s | %s | %s\n", v.ID, v.FullName, v.Group, v.Email)
	yes, e := c.confirm("Delete this student? (y/n): ")
	if e != nil {
		return e
	}
	if !yes {
		return c.message("Deletion cancelled.")
	}
	return c.result("Student", c.app.Students.Delete(id))
}
func (c *CLI) addCriteria() error {
	name, e := c.readNonempty("Name: ")
	if e != nil {
		return e
	}
	max, e := c.readInt("Maximum score: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	weight, e := c.readFloat("Weight (0..1): ", 0, 1)
	if e != nil {
		return e
	}
	v, e := c.app.Rubric.Add(application.CriteriaRequest{Name: name, MaxScore: max, Weight: weight})
	if e != nil {
		return c.showError(e)
	}
	return c.message(fmt.Sprintf("Added criterion #%d %s.", v.ID, v.Name))
}
func (c *CLI) deleteCriteria() error {
	id, e := c.readInt("Criteria ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	all, e := c.app.Rubric.List()
	if e != nil {
		return c.showError(e)
	}
	var label string
	for _, v := range all {
		if v.ID == id {
			label = fmt.Sprintf("#%d %s (max %d, weight %g)", v.ID, v.Name, v.MaxScore, v.Weight)
		}
	}
	if label == "" {
		return c.showError(fmt.Errorf("criterion #%d not found", id))
	}
	fmt.Fprintf(c.out, "Selected: %s\n", label)
	yes, e := c.confirm("Delete this criterion? (y/n): ")
	if e != nil {
		return e
	}
	if !yes {
		return c.message("Deletion cancelled.")
	}
	return c.result("Criterion", c.app.Rubric.Remove(id))
}
func (c *CLI) setEvaluation() error {
	sid, e := c.readInt("Student ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	cid, e := c.readInt("Criteria ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	score, e := c.readInt("Score: ", 0, math.MaxInt)
	if e != nil {
		return e
	}
	comment, e := c.read("Comment (optional): ")
	if e != nil {
		return e
	}
	_, e = c.app.Evaluations.Set(application.EvaluationRequest{StudentID: sid, CriteriaID: cid, Score: score, Comment: comment})
	return c.result("Evaluation", e)
}
func (c *CLI) deleteEvaluation() error {
	sid, e := c.readInt("Student ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	cid, e := c.readInt("Criteria ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	r, e := c.app.Evaluations.ForStudent(sid)
	if e != nil {
		return c.showError(e)
	}
	var selected string
	for _, v := range r.Evaluations {
		if v.CriteriaID == cid {
			score := fmt.Sprint(v.Score)
			if v.Score < 0 {
				score = "N/A"
			}
			selected = fmt.Sprintf("student #%d, criterion #%d, score %s, comment %q", sid, cid, score, v.Comment)
		}
	}
	if selected == "" || strings.Contains(selected, "score N/A") {
		return c.showError(errors.New("evaluation not found"))
	}
	fmt.Fprintf(c.out, "Selected: %s\n", selected)
	yes, e := c.confirm("Delete this evaluation? (y/n): ")
	if e != nil {
		return e
	}
	if !yes {
		return c.message("Deletion cancelled.")
	}
	return c.result("Evaluation", c.app.Evaluations.Remove(sid, cid))
}
func (c *CLI) studentResult() error {
	id, e := c.readInt("Student ID: ", 1, math.MaxInt)
	if e != nil {
		return e
	}
	r, e := c.app.Evaluations.ForStudent(id)
	if e != nil {
		return c.showError(e)
	}
	return c.showPage("Result: "+r.Student.FullName, resultRows([]application.StudentResult{r}))
}
func (c *CLI) allResults() error {
	v, e := c.app.Evaluations.AllResults()
	if e != nil {
		return c.showError(e)
	}
	return c.showPage("Results", resultRows(v))
}
func resultRows(rows []application.StudentResult) []string {
	out := []string{}
	for _, r := range rows {
		out = append(out, fmt.Sprintf("#%d %s — final %d", r.Student.ID, r.Student.FullName, r.FinalScore))
		for _, v := range r.Evaluations {
			s := fmt.Sprint(v.Score)
			if v.Score < 0 {
				s = "N/A"
			}
			out = append(out, fmt.Sprintf("  criterion #%d: %s %s", v.CriteriaID, s, v.Comment))
		}
	}
	return out
}
func (c *CLI) export() error {
	if e := c.app.Export(); e != nil {
		return c.showError(e)
	}
	return c.message("Exported to ./results.csv.")
}
func (c *CLI) studentFields() (application.StudentRequest, error) {
	first, e := c.readNonempty("First name: ")
	if e != nil {
		return application.StudentRequest{}, e
	}
	last, e := c.readNonempty("Last name: ")
	if e != nil {
		return application.StudentRequest{}, e
	}
	middle, e := c.read("Middle name (optional): ")
	if e != nil {
		return application.StudentRequest{}, e
	}
	group, e := c.readNonempty("Group: ")
	if e != nil {
		return application.StudentRequest{}, e
	}
	email, e := c.readNonempty("Email: ")
	return application.StudentRequest{FirstName: first, LastName: last, MiddleName: middle, Group: group, Email: email}, e
}
func (c *CLI) showPage(title string, rows []string) error {
	page := 0
	for {
		start := page * pageSize
		end := min(start+pageSize, len(rows))
		if e := c.title(title); e != nil {
			return e
		}
		if len(rows) == 0 {
			fmt.Fprintln(c.out, "(empty)")
		} else {
			for i, v := range rows[start:end] {
				fmt.Fprintf(c.out, "%s. %s\n", menuKey(i), v)
			}
		}
		if e := c.end(title); e != nil {
			return e
		}
		s, e := c.read("< previous | > next | q back\n> ")
		if e != nil {
			return e
		}
		switch s {
		case "q":
			return nil
		case "<":
			if page > 0 {
				page--
			} else {
				fmt.Fprintln(c.out, "No previous page.")
			}
		case ">":
			if end < len(rows) {
				page++
			} else {
				fmt.Fprintln(c.out, "No next page.")
			}
		default:
			fmt.Fprintln(c.out, "Invalid command.")
		}
	}
}
func (c *CLI) title(s string) error {
	_, e := fmt.Fprintf(c.out, "\n%s %s %s\n", menuEdge, s, menuEdge)
	return e
}
func (c *CLI) end(s string) error {
	_, e := fmt.Fprintln(c.out, strings.Repeat("≡", 2*utf8.RuneCountInString(menuEdge)+utf8.RuneCountInString(s)+2))
	return e
}
func (c *CLI) read(prompt string) (string, error) {
	if _, e := fmt.Fprint(c.out, prompt); e != nil {
		return "", e
	}
	if !c.in.Scan() {
		if e := c.in.Err(); e != nil {
			return "", e
		}
		return "", io.EOF
	}
	return strings.TrimSpace(c.in.Text()), nil
}
func (c *CLI) readNonempty(prompt string) (string, error) {
	for {
		s, e := c.read(prompt)
		if e != nil {
			return "", e
		}
		if s != "" {
			return s, nil
		}
		fmt.Fprintln(c.out, "Value cannot be empty.")
	}
}
func (c *CLI) readInt(prompt string, minV, maxV int) (int, error) {
	for {
		s, e := c.read(prompt)
		if e != nil {
			return 0, e
		}
		n, pe := strconv.Atoi(s)
		if pe == nil && n >= minV && n <= maxV {
			return n, nil
		}
		fmt.Fprintf(c.out, "Enter a whole number from %d to %d.\n", minV, maxV)
	}
}
func (c *CLI) readFloat(prompt string, minV, maxV float64) (float64, error) {
	for {
		s, e := c.read(prompt)
		if e != nil {
			return 0, e
		}
		n, pe := strconv.ParseFloat(s, 64)
		if pe == nil && !math.IsNaN(n) && !math.IsInf(n, 0) && n >= minV && n <= maxV {
			return n, nil
		}
		fmt.Fprintf(c.out, "Enter a finite number from %g to %g.\n", minV, maxV)
	}
}
func (c *CLI) confirm(prompt string) (bool, error) {
	for {
		s, e := c.read(prompt)
		if e != nil {
			return false, e
		}
		switch strings.ToLower(s) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		default:
			fmt.Fprintln(c.out, "Enter y or n.")
		}
	}
}
func (c *CLI) message(s string) error  { _, e := fmt.Fprintln(c.out, s); return e }
func (c *CLI) showError(e error) error { fmt.Fprintf(c.out, "Not completed: %v\n", e); return nil }
func (c *CLI) result(what string, e error) error {
	if e != nil {
		return c.showError(e)
	}
	return c.message(what + " operation completed.")
}
func menuKey(i int) string {
	if i == 9 {
		return "0"
	}
	return strconv.Itoa(i + 1)
}
func menuIndex(s string, count int) (int, bool) {
	i := -1
	if s == "0" {
		i = 9
	} else if n, e := strconv.Atoi(s); e == nil {
		i = n - 1
	}
	return i, i >= 0 && i < count && i < pageSize
}
