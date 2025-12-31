package runtime

import (
	"context"

	"github.com/yourorg/obs/internal/health"
)

type Service interface {
	Name() string
	Run(ctx context.Context) error
	Health() health.Report
}
