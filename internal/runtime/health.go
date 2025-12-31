package runtime

import (
	"github.com/yourorg/obs/internal/health"
)

func Health(r *Runtime) []health.Report {
	reports := make([]health.Report, 0, len(r.services))

	for _, svc := range r.services {
		reports = append(reports, svc.Health())
	}

	return reports
}
