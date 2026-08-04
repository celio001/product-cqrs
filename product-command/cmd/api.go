package cmd

import (
	"github.com/celio001/product-command/config"
	brands_repository "github.com/celio001/product-command/internal/modules/brands/repository"
	brands_service "github.com/celio001/product-command/internal/modules/brands/service"
	"github.com/celio001/product-command/internal/fiber"
	categories_repository "github.com/celio001/product-command/internal/modules/categories/repository"
	categories_service "github.com/celio001/product-command/internal/modules/categories/service"
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

	brandsRepo := brands_repository.NewBrandsRepository(pg)
	brandsSvc := brands_service.NewBrandSvc(brandsRepo)

	categoriesRepo := categories_repository.NewCategoriesRepo(pg)
	categoriesSvc := categories_service.NewCategoriesSvc(categoriesRepo)

	f := fiber.CreateApp(brandsSvc, categoriesSvc)

	lifecycle.New(cmd.Context(), "api", f.Start, f.Stop)

	return nil
}
