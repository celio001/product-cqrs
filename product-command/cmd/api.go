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
	"github.com/celio001/product-command/pkg/postgres"
	"github.com/spf13/cobra"
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

	pg, err := postgres.ConectPostgres(cmd.Context(), config.GetString("POSTGRES_DB_DSN"))
	if err != nil {
		logger.Fatal("error conection postrgres database")
	}

	if err := pg.Ping(cmd.Context()); err != nil {
		pg.Close()
		logger.Fatal("error ping postrgres database")
	}

	tx := database.New(pg)

	productTopic := kafka.NewKafkaProducer(config.GetStrings("KAFKA_BROKERS"), config.GetString("KAFKA_PRODUCT_TOPIC"))
	brandTopic := kafka.NewKafkaProducer(config.GetStrings("KAFKA_BROKERS"), config.GetString("KAFKA_BRAND_TOPIC"))

	producer := producer.NewProducerCommand(productTopic, brandTopic)

	brandsRepo := brands_repository.NewBrandsRepository(pg, tx)
	brandsSvc := brands_service.NewBrandSvc(brandsRepo, producer)

	categoriesRepo := categories_repository.NewCategoriesRepo(pg, tx)
	categoriesSvc := categories_service.NewCategoriesSvc(categoriesRepo, producer)

	inventoryRepo := inventory_repository.NewInventoryRepo(pg, tx)
	fiscalRepo := fiscal_repository.NewFiscalRepo(pg, tx)
	productRepo := product_repository.NewProductRepo(pg, tx)

	

	productSvc := product_service.NewProductSvc(productRepo, fiscalRepo, inventoryRepo, categoriesRepo, brandsRepo, producer)

	f := fiber.CreateApp(brandsSvc, categoriesSvc, productSvc)

	lifecycle.New(cmd.Context(), "api", f.Start, f.Stop)

	return nil
}
