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

func (f *fiberContext) JSON(code int, obj interface{}) error {
	return f.ctx.Status(code).JSON(obj)
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

func (f *fiberContext) Raw() interface{} {
	return f.ctx
}
