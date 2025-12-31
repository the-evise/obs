package config

import (
	"fmt"
	"os"
	"time"
)

/*
   Responsibility:
       - Defines Config
       - Applies Defaults
       - Validates inputs
   Owns:
       - Service Identity
       - Sampling ratios
       - Ports and endpoints
       - Environment flags
   Does NOT:
       - Initialize anything
       - Touch global state
   Used by:
       - runtime.go (only)
*/

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Interval time.Duration
}

func Load() (*Config, error) {
	interval := time.Minute
	if intervalStr := os.Getenv("CHECK_INTERVAL"); intervalStr != "" {
		parsed, err := time.ParseDuration(intervalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid CHECK_INTERVAL: %w", err)
		}
		interval = parsed
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	name := os.Getenv("APP_NAME")
	if name == "" {
		name = "checker"
	}

	cfg := &Config{
		App: AppConfig{
			Name:     name,
			Env:      env,
			Interval: interval,
		},
	}

	return cfg, nil
}
