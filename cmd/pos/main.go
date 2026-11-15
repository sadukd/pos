package main

import (
	"go.uber.org/dig"
	"log"
	"os"
	"pos/internal/adapters/cli"
	exporter "pos/internal/adapters/export"
	"pos/internal/adapters/repository"
	"pos/internal/application"
	"pos/internal/application/ports"
	"pos/internal/domain"
)

func main() {
	c := dig.New()
	provide(c, func() ports.StudentRepository { return repository.NewStudents() }, func() ports.RubricRepository { return repository.NewRubrics() }, func() ports.StudentEvaluationsRepository { return repository.NewEvaluations() }, func() ports.ResultExporter { return exporter.NewCSVExporter() }, domain.NewSequentialIDs, func(g *domain.SequentialIDs) domain.StudentIDGenerator { return g }, func(g *domain.SequentialIDs) domain.CriteriaIDGenerator { return g }, domain.NewEvaluationFactory, func() domain.ScoreStrategy { return domain.WeightedScoreStrategy{} }, domain.NewScoreCalculator, application.NewStudentService, application.NewRubricService, application.NewEvaluationService, application.NewPOSFacade, func(f *application.POSFacade) *cli.CLI { return cli.New(f, os.Stdin, os.Stdout) })
	if err := c.Invoke(func(v *cli.CLI) {
		if e := v.Run(); e != nil {
			log.Print(e)
		}
	}); err != nil {
		log.Fatal(err)
	}
}
func provide(c *dig.Container, values ...interface{}) {
	for _, v := range values {
		if err := c.Provide(v); err != nil {
			log.Fatal(err)
		}
	}
}
