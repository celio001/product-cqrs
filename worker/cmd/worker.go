package cmd

import (
	"context"
	"errors"

	"github.com/celio001/product-cqrs/worker/config"
	"github.com/celio001/product-cqrs/worker/internal/modules/consumer"
	"github.com/celio001/product-cqrs/worker/internal/modules/producer-dlq"
	product_respository "github.com/celio001/product-cqrs/worker/internal/modules/product/respository"
	product_service "github.com/celio001/product-cqrs/worker/internal/modules/product/service"
	"github.com/celio001/product-cqrs/worker/pkg/kafka"
	"github.com/celio001/product-cqrs/worker/pkg/lifecycle"
	"github.com/celio001/product-cqrs/worker/pkg/logger"
	"github.com/celio001/product-cqrs/worker/pkg/mongodb"
	opentelemetry "github.com/celio001/product-cqrs/worker/pkg/openTelemetry"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	workerCommand = &cobra.Command{
		Use:  "worker",
		RunE: worker,
	}
)

func init() {
	rootCmd.AddCommand(workerCommand)
}

func worker(cmd *cobra.Command, args []string) error {
	cfgs := config.LoadEnvs()

	logger.Init(cfgs.ServiceName, cfgs.ServiceVersion, cfgs.Env)
	defer logger.Sync()

	tp, err := opentelemetry.InitTracerProvider(cmd.Context(), cfgs.ServiceName, cfgs.ServiceVersion, cfgs.JaegerConfig.URL)
	if err != nil {
		logger.Fatal("error initializing tracer provider", zap.String("error", err.Error()))
	}
	defer func() {
		if err := tp.Shutdown(cmd.Context()); err != nil {
			logger.Fatal("error shutting down tracer provider", zap.String("error", err.Error()))
		}
	}()

	tracer := otel.Tracer("product-worker")

	mongo, err := mongodb.ConnectMongoDB(cmd.Context(), cfgs.MongoDB.DSN)
	if err != nil {
		logger.Fatal("error connecting mongo database",
			zap.String("error.type", "ConnectDatabase"),
			zap.String("error", err.Error()),
		)
	}

	err = mongo.Ping(cmd.Context(), nil)
	if err != nil {
		logger.Fatal("error ping mongo",
			zap.String("error.type", "ConnectDatabse"),
			zap.String("error", err.Error()))
	}

	productTopic := kafka.NewKafkaConsumer(cfgs.KafkaCfg.KafkaBrokers, cfgs.KafkaCfg.ProductTopic)
	defer productTopic.Close()

	productDlqTopic := kafka.NewKafkaProducer(cfgs.KafkaCfg.KafkaBrokers, cfgs.KafkaCfg.ProdctDlqTopic)
	defer productDlqTopic.Close()

	productConsumer := consumer.NewConsumerTopics(productTopic)
	productDlqProducer := producer.NewProducerDlq(productDlqTopic)

	productRepo := product_respository.NewProductRepository(mongo)
	productSvc := product_service.NewProductService(productRepo, productConsumer, productDlqProducer, tracer)

	return lifecycle.New(cmd.Context(), cfgs.ServiceName,
		func(ctx context.Context) error {
			productSvc.CreateProductSvc(ctx)
			return nil
		},
		func(ctx context.Context) error {
			return errors.Join(
				productDlqTopic.Close(),
				productTopic.Close(),
				mongo.Disconnect(ctx),
			)
		},
	)
}
