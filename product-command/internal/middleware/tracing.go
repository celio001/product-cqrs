package middleware

import (
	"github.com/celio001/product-command/config"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type fiberHeaderCarrier struct {
	ctx fiber.Ctx
}

func (c fiberHeaderCarrier) Get(key string) string {
	return c.ctx.Get(key)
}

func (fiberHeaderCarrier) Set(string, string) {}

func (c fiberHeaderCarrier) Keys() []string {
	headers := c.ctx.GetReqHeaders()
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	return keys
}

func OpenTelemetryMiddleware() fiber.Handler {
	tracer := otel.Tracer(config.GetString("SERVICE_NAME"))

	return func(c fiber.Ctx) error {
		parentCtx := otel.GetTextMapPropagator().Extract(c.Context(), fiberHeaderCarrier{ctx: c})
		ctx, span := tracer.Start(parentCtx, "HTTP "+c.Method(), trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		c.SetContext(ctx)
		err := c.Next()

		route := c.FullPath()
		if route == "" {
			route = c.Path()
		}
		span.SetName(c.Method() + " " + route)
		span.SetAttributes(
			attribute.String("http.request.method", c.Method()),
			attribute.String("url.path", c.Path()),
			attribute.String("url.route", route),
			attribute.Int("http.response.status_code", c.Response().StatusCode()),
		)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}

		return err
	}
}
