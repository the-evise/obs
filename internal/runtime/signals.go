package runtime

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/obs/internal/log"
)

func (r *Runtime) listenForSignals() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-ch
		r.Logger.Info("shutdown signal received", log.F("signal", sig.String()))
		r.Shutdown()
	}()
}
