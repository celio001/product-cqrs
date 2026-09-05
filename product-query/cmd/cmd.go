package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/celio001/product-cqrs/product-query/config"
	product_repository "github.com/celio001/product-cqrs/product-query/internal/modules/product/repository"
	product_service "github.com/celio001/product-cqrs/product-query/internal/modules/product/service"
	query_router "github.com/celio001/product-cqrs/product-query/internal/router"
	"github.com/celio001/product-cqrs/product-query/pkg/lifecycle"
	"github.com/celio001/product-cqrs/product-query/pkg/logger"
	mongodb "github.com/celio001/product-cqrs/product-query/pkg/mongo"
	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var apiCommand = &cobra.Command{
	Use:  "api",
	RunE: apiExecute,
}

func init() {
	rootCmd.AddCommand(apiCommand)
}

func apiExecute(cmd *cobra.Command, args []string) error {
	cfgs := config.LoadEnvs()

	logger.Init(cfgs.ServiceName, cfgs.ServiceVersion, cfgs.Env)
	defer logger.Sync()

	mongoClient, err := mongodb.ConnectMongoDB(cmd.Context(), cfgs.MongoDB.DSN)
	if err != nil {
		logger.Fatal("error connecting mongo database",
			zap.String("error.type", "ConnectDatabase"),
			zap.String("error", err.Error()),
		)
	}

	if err := mongoClient.Ping(cmd.Context(), nil); err != nil {
		mongoClient.Disconnect(cmd.Context())
		logger.Fatal("error ping mongo",
			zap.String("error.type", "ConnectDatabase"),
			zap.String("error", err.Error()),
		)
	}

	productRepo := product_repository.NewProductRepository(mongoClient)
	productSvc := product_service.NewProductService(productRepo)

	r := query_router.NewSetupRouters(chi.NewRouter(), productSvc)

	srv := &http.Server{
		Addr:              ":" + cfgs.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return lifecycle.New(cmd.Context(), cfgs.ServiceName,
		func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("http server stopped unexpectedly",
						zap.String("error.message", err.Error()),
						zap.String("error.code", "HTTP_SERVER_ERROR"),
					)
				}
			}()
			return nil
		},
		func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := srv.Shutdown(shutdownCtx); err != nil {
				return err
			}

			return mongoClient.Disconnect(ctx)
		},
	)
}
