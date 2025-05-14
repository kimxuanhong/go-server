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

func (g *ginContext) JSON(code int, obj interface{}) {
	g.ctx.JSON(code, obj)
}

func (g *ginContext) Abort() {
	g.ctx.Abort()
}

func (g *ginContext) AbortWithStatusJSON(code int, obj interface{}) {
	g.ctx.AbortWithStatusJSON(code, obj)
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

// Method returns the HTTP method of the request.
func (g *ginContext) Method() string {
	return g.ctx.Request.Method
}

// Path returns the request path.
func (g *ginContext) Path() string {
	return g.ctx.FullPath()
}

// Next calls the next middleware in the chain.
func (g *ginContext) Next() {
	g.ctx.Next()
}

func (g *ginContext) Raw() interface{} {
	return g.ctx
}

func (g *ginContext) Set(key string, value interface{}) {
	g.ctx.Set(key, value)
}

func (g *ginContext) Get(key string) interface{} {
	val, _ := g.ctx.Get(key)
	return val
}

func (g *ginContext) GetString(key string) string {
	return g.ctx.GetString(key)
}

func (g *ginContext) GetInt(key string) int {
	return g.ctx.GetInt(key)
}
