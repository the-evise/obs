package main

import (
	"context"
	"log"
	"time"

	"github.com/yourorg/obs/internal/config"
	internalLog "github.com/yourorg/obs/internal/log"
	"github.com/yourorg/obs/internal/runtime"
	"github.com/yourorg/obs/internal/service"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := internalLog.NewStdLogger()

	rt := runtime.NewRuntime(ctx, logger, cfg, 10*time.Second)

	checker := service.NewChecker(logger, cfg.App.Interval, cfg.App.Name)
	rt.RegisterServices(checker)

	if err := rt.Run(); err != nil && err != context.Canceled {
		log.Fatalf("runtime stopped: %v", err)
	}
}
