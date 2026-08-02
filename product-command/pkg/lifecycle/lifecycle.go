package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/celio001/product-command/config"
	"github.com/celio001/product-command/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func New(ctx context.Context, name string, onStart, onStop func(ctx context.Context) error) {
	lifeCtx, cancel := context.WithCancel(ctx)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		sig := <-c
		logger.Info("received signal, initing shutdown",
			zap.String("api.product-command.signal", sig.String()),
			zap.String("event.action", "shutdown_initiated"),
		)
		cancel()
	}()
	g, gCtx := errgroup.WithContext(lifeCtx)

	g.Go(func() error {
		logger.Info("server started",
			zap.String("app.product-command.port", config.GetString("HTTP_PORT")),
			zap.String("event.action", "server_started"),
			zap.String("event.outcome", "success"),
		)
		return onStart(gCtx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("initiating graceful shutdown",
			zap.String("app.product-command.component", name),
			zap.String("event.action", "graceful_shutdown_started"),
		)

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()

		err := onStop(stopCtx)

		return err
	})

	if err := g.Wait(); err != nil {
		logger.Error("shutdown completed with error",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "SHUTDOWN_FAILED"),
			zap.String("app.product-command.component", name),
			zap.String("event.action", "shutdown_completed"),
			zap.String("event.outcome", "failure"),
		)
	} else {
		logger.Info("shutdown completed succesfully",
			zap.String("app.product-command.component", name),
			zap.String("event.action", "shutdown_completed"),
			zap.String("event.outcome", "success"))
	}
}
