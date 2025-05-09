package fiber

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/kimxuanhong/go-server/core"
)

type fiberContext struct {
	ctx *fiber.Ctx
}

func (f *fiberContext) Context() context.Context {
	return f.ctx.Context()
}

func (f *fiberContext) Param(name string) string {
	return f.ctx.Params(name)
}

func (f *fiberContext) Query(name string) string {
	return f.ctx.Query(name)
}

func (f *fiberContext) Header(name string) string {
	return f.ctx.Get(name)
}

func (f *fiberContext) Bind(obj interface{}) error {
	return f.ctx.BodyParser(obj)
}

func (f *fiberContext) JSON(code int, obj interface{}) {
	if err := f.ctx.Status(code).JSON(obj); err != nil {
		// Handle the error, for example by logging it
		_ = f.ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send JSON response"})
	}
}

func (f *fiberContext) String(code int, msg string) error {
	return f.ctx.Status(code).SendString(msg)
}

func (f *fiberContext) Status(code int) core.Context {
	f.ctx.Status(code)
	return f
}

func (f *fiberContext) SetHeader(key, value string) {
	f.ctx.Set(key, value)
}

// Method returns the HTTP method of the request.
func (f *fiberContext) Method() string {
	return f.ctx.Method()
}

// Path returns the request path.
func (f *fiberContext) Path() string {
	return f.ctx.Path()
}

// Next calls the next middleware in the chain.
func (f *fiberContext) Next() error {
	return f.ctx.Next()
}

func (f *fiberContext) Raw() interface{} {
	return f.ctx
}
