package fiber

import (
	"context"
	"fmt"

	"github.com/celio001/product-command/config"
	v1 "github.com/celio001/product-command/internal/fiber/v1"
	"github.com/celio001/product-command/internal/middleware"
	f "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func (h HttpServer) Start(ctx context.Context) error {

	h.app.Use(middleware.LoggerMiddleware())
	api := h.app.Group("/api")

	h.app.Get(healthcheck.LivenessEndpoint, healthcheck.New())

	h.app.Get("/health", healthcheck.New(healthcheck.Config{
		ResponseFormat: healthcheck.FormatJSON,
	}))

	v1Router := api.Group(v1.HandlerPath)
	v1.RegisterRouter(v1Router, h.brandsSvc, h.categoriesSvc)

	addr := fmt.Sprint(":", config.GetString("HTTP_PORT"))
	
	return h.app.Listen(addr, f.ListenConfig{
		DisableStartupMessage: true,
	})
}

func (h HttpServer) Stop(ctx context.Context) error {
	return h.app.Shutdown()
}
