package fiber

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/kimxuanhong/go-server/core"
	"log"
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

func (f *fiberContext) Abort() {
	// Fiber không có Abort, nên ta dừng luồng bằng cách không gọi next()
}

func (f *fiberContext) AbortWithStatusJSON(code int, obj interface{}) {
	f.JSON(code, obj)
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
func (f *fiberContext) Next() {
	if err := f.ctx.Next(); err != nil {
		log.Printf("Next() error: %v", err)
	}
}

func (f *fiberContext) Raw() interface{} {
	return f.ctx
}

func (f *fiberContext) Set(key string, value interface{}) {
	f.ctx.Locals(key, value)
}

func (f *fiberContext) Get(key string) interface{} {
	return f.ctx.Locals(key)
}

func (f *fiberContext) GetString(key string) string {
	val := f.ctx.Locals(key)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

func (f *fiberContext) GetInt(key string) int {
	val := f.ctx.Locals(key)
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
