package response

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type Builder struct {
	status  int
	message string
	data    any
	meta    any
	err     any
}

func New() *Builder {
	return &Builder{
		status: http.StatusOK,
	}
}

func (b *Builder) Status(code int) *Builder {
	b.status = code
	return b
}

func (b *Builder) Message(msg string) *Builder {
	b.message = msg
	return b
}

func (b *Builder) Data(d any) *Builder {
	b.data = d
	return b
}

func (b *Builder) Meta(m any) *Builder {
	b.meta = m
	return b
}

func (b *Builder) Error(err any) *Builder {
	b.err = err
	return b
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Error   any    `json:"error,omitempty"`
}

type Builder struct {
	message string
	data    any
	meta    any
	err     any
}

func (b *Builder) Send(status any, c fiber.Ctx) fiber.Ctx {
	resp := Response{
		Status:  b.status,
		Message: b.message,
		Data:    b.data,
		Meta:    b.meta,
		Error:   b.err,
	}
	return c.Status(b.status).JSONP(resp)
}
