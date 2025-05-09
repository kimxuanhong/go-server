package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kimxuanhong/go-server/core"
)

type ginContext struct {
	ctx *gin.Context
}

func (g *ginContext) Context() context.Context {
	return g.ctx.Request.Context()
}

func (g *ginContext) Param(name string) string {
	return g.ctx.Param(name)
}

func (g *ginContext) Query(name string) string {
	return g.ctx.Query(name)
}

func (g *ginContext) Header(name string) string {
	return g.ctx.GetHeader(name)
}

func (g *ginContext) Bind(obj interface{}) error {
	return g.ctx.ShouldBind(obj)
}

func (g *ginContext) JSON(code int, obj interface{}) error {
	g.ctx.JSON(code, obj)
	return nil
}

func (g *ginContext) String(code int, msg string) error {
	g.ctx.String(code, msg)
	return nil
}

func (g *ginContext) Status(code int) core.Context {
	g.ctx.Status(code)
	return g
}

func (g *ginContext) SetHeader(key, value string) {
	g.ctx.Header(key, value)
}

func (g *ginContext) Raw() interface{} {
	return g.ctx
}
