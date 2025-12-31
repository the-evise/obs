package service

import (
	"context"
	"time"

	"github.com/yourorg/obs/internal/health"
	"github.com/yourorg/obs/internal/log"
	"github.com/yourorg/obs/internal/runtime"
)

type Checker struct {
	logger   log.Logger
	interval time.Duration
	health   health.Report
}

func NewChecker(logger log.Logger, interval time.Duration, name string) *Checker {
	if interval <= 0 {
		interval = time.Minute
	}

	if name == "" {
		name = "checker"
	}

	return &Checker{
		logger:   logger,
		interval: interval,
		health: health.Report{
			Name:   name,
			Status: health.StatusHealthy,
		},
	}
}

func (c *Checker) Name() string {
	return c.health.Name
}

func (c *Checker) Health() health.Report {
	return c.health
}

func (c *Checker) Run(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("checker stopping")
			return ctx.Err()

		case <-ticker.C:
			if err := c.checkOnce(); err != nil {
				c.health.Status = health.StatusDegraded
				c.health.Details = err.Error()

				c.logger.Error("check failed", log.F("error", err))
				return runtime.Fatal(err)
			}

			c.health.Status = health.StatusHealthy
			c.health.Details = ""
		}
	}
}

func (c *Checker) checkOnce() error {
	// TODO: actual check logic
	return nil
}
