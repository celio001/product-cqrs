package middleware

import (
	"time"

	"github.com/celio001/product-command/pkg/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		tl := time.Now()

		err := c.Next()

		routePattern := c.Route().Path
		if routePattern == "" {
			routePattern = logger.SanitizeForLog(c.Path())
		}

		logger.Info("HTTP request processed",
			zap.String("http.request.method", c.Method()),
			zap.String("url.path", logger.SanitizeForLog(c.Path())),
			zap.String("url.route", routePattern),
			zap.String("client.ip", c.IP()),
			zap.Int("http.response.status)code", c.Response().StatusCode()),
			zap.Int64("event.duration", time.Since(tl).Microseconds()),
			zap.Int("http.response.body.bytes", len(c.Response().Body())),
		)
		return err
	}
}
