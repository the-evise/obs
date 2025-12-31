package example

import (
	"context"

	"github.com/yourorg/obs/internal/health"
	"github.com/yourorg/obs/internal/runtime"
)

type ExampleService struct {
	health health.Report
}

func New() *ExampleService {
	return &ExampleService{
		health: health.Report{
			Name:   "example",
			Status: health.StatusHealthy,
		},
	}
}

func (s *ExampleService) Name() string {
	return "example"
}

func (s *ExampleService) Health() health.Report {
	return s.health
}

func (s *ExampleService) Run(ctx context.Context) error {
	if err := s.doWork(ctx); err != nil {
		s.health.Status = health.StatusUnhealthy
		s.health.Details = err.Error()

		return runtime.Fatal(err)
	}
	return nil
}
