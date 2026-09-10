package cmd

import (
	"github.com/celio001/product-command/config"
	"github.com/celio001/product-command/internal/database"
	"github.com/celio001/product-command/internal/fiber"
	brands_repository "github.com/celio001/product-command/internal/modules/brands/repository"
	brands_service "github.com/celio001/product-command/internal/modules/brands/service"
	categories_repository "github.com/celio001/product-command/internal/modules/categories/repository"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
	fiscal_repository "github.com/celio001/product-command/internal/modules/fiscal/repository"
	inventory_repository "github.com/celio001/product-command/internal/modules/inventory/repository"
	"github.com/celio001/product-command/internal/modules/producer"
	product_repository "github.com/celio001/product-command/internal/modules/product/repository"
	product_service "github.com/celio001/product-command/internal/modules/product/service"
	"github.com/celio001/product-command/pkg/kafka"
	"github.com/celio001/product-command/pkg/lifecycle"
	"github.com/celio001/product-command/pkg/logger"
	opentelemetry "github.com/celio001/product-command/pkg/openTelemetry"
	"github.com/celio001/product-command/pkg/postgres"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	httpCommand = &cobra.Command{
		Use:  "api",
		RunE: httpExecute,
	}
)

func init() {
	rootCmd.AddCommand(httpCommand)
}

func httpExecute(cmd *cobra.Command, args []string) error {

	logger.Init(config.GetString("SERVICE_NAME"), config.GetString("SERVICE_VERSION"), config.GetString("ENV"))

	tp, err := opentelemetry.InitTracerProvider(cmd.Context(), config.GetString("SERVICE_NAME"), config.GetString("SERVICE_VERSION"), config.GetString("JAEGER_URL"))
	if err != nil {
		logger.Fatal("error initializing tracer provider", zap.String("error", err.Error()))
	}
	defer func() {
		if err := tp.Shutdown(cmd.Context()); err != nil {
			logger.Fatal("error shutting down tracer provider", zap.String("error", err.Error()))
		}
	}()

	tracer := otel.Tracer(config.GetString("SERVICE_NAME"))

	pg, err := postgres.ConectPostgres(cmd.Context(), config.GetString("POSTGRES_DB_DSN"))
	if err != nil {
		logger.Fatal("error conection postrgres database", zap.String("error", err.Error()))
	}

	if err := pg.Ping(cmd.Context()); err != nil {
		pg.Close()
		logger.Fatal("error ping postrgres database", zap.String("error", err.Error()))
	}

	tx := database.New(pg)

	productTopic := kafka.NewKafkaProducer(config.GetStrings("KAFKA_BROKERS"), config.GetString("KAFKA_PRODUCT_TOPIC"))
	productTopicDeleted := kafka.NewKafkaProducer(config.GetStrings("KAFKA_BROKERS"), config.GetString("KAFKA_PRODUCT_DELETED_TOPIC"))
	brandTopic := kafka.NewKafkaProducer(config.GetStrings("KAFKA_BROKERS"), config.GetString("KAFKA_BRAND_TOPIC"))
	categoryTopic := kafka.NewKafkaProducer(config.GetStrings("KAFKA_BROKERS"), config.GetString("KAFKA_CATEGORY_TOPIC"))

	producer := producer.NewProducerCommand(productTopic, productTopicDeleted, brandTopic, categoryTopic, tracer)

	brandsRepo := brands_repository.NewBrandsRepository(pg, tx)
	brandsSvc := brands_service.NewBrandSvc(brandsRepo, producer, tracer)

	categoriesRepo := categories_repository.NewCategoriesRepo(pg, tx)
	categoriesSvc := categories_service.NewCategoriesSvc(categoriesRepo, producer)

	inventoryRepo := inventory_repository.NewInventoryRepo(pg, tx)
	fiscalRepo := fiscal_repository.NewFiscalRepo(pg, tx)
	productRepo := product_repository.NewProductRepo(pg, tx)

	productSvc := product_service.NewProductSvc(productRepo, fiscalRepo, inventoryRepo, categoriesRepo, brandsRepo, producer, tracer)

	f := fiber.CreateApp(brandsSvc, categoriesSvc, productSvc)

	lifecycle.New(cmd.Context(), "api", f.Start, f.Stop)

	return nil
}
