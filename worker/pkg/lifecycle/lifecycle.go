package lifecycle

import (
	"context"
	"errors"
	"time"

	"os/signal"
	"syscall"

	"github.com/celio001/product-cqrs/worker/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func New(ctx context.Context, name string, onStart, onStop func(ctx context.Context) error) error {
	lifeCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, gCtx := errgroup.WithContext(lifeCtx)

	g.Go(func() error {
		logger.Info("worker started",
			zap.String("app.component", name),
			zap.String("event.action", "worker_started"),
			zap.String("event.outcome", "success"),
		)
		return onStart(gCtx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("initiating graceful shutdown",
			zap.String("app.component", name),
			zap.String("event.action", "graceful_shutdown_started"),
		)

		stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()

		err := onStop(stopCtx)

		return err
	})

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("shutdown completed with error",
			zap.String("error.message", err.Error()),
			zap.String("error.code", "SHUTDOWN_FAILED"),
			zap.String("app.component", name),
			zap.String("event.action", "shutdown_completed"),
			zap.String("event.outcome", "failure"),
		)
	} else {
		logger.Info("shutdown completed successfully",
			zap.String("app.component", name),
			zap.String("event.action", "shutdown_completed"),
			zap.String("event.outcome", "success"))
	}

	return err
}
