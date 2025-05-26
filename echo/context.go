package echo

import (
	"context"
	"github.com/kimxuanhong/go-server/core"
	"github.com/labstack/echo/v4"
)

type echoContext struct {
	ctx echo.Context
}

func (e *echoContext) Context() context.Context {
	return e.ctx.Request().Context()
}

func (e *echoContext) Param(name string) string {
	return e.ctx.Param(name)
}

func (e *echoContext) Query(name string) string {
	return e.ctx.QueryParam(name)
}

func (e *echoContext) Header(name string) string {
	return e.ctx.Request().Header.Get(name)
}

func (e *echoContext) Bind(obj interface{}) error {
	return e.ctx.Bind(obj)
}

func (e *echoContext) JSON(code int, obj interface{}) {
	_ = e.ctx.JSON(code, obj)
}

func (e *echoContext) Abort() {
	e.ctx.Set("abort", true)
}

func (e *echoContext) AbortWithStatusJSON(code int, obj interface{}) {
	e.ctx.Response().WriteHeader(code)
	_ = e.ctx.JSON(code, obj)
	e.ctx.Set("abort", true)
}

func (e *echoContext) String(code int, msg string) error {
	return e.ctx.String(code, msg)
}

func (e *echoContext) Status(code int) core.Context {
	e.ctx.Response().WriteHeader(code)
	return e
}

func (e *echoContext) SetHeader(key, value string) {
	e.ctx.Response().Header().Set(key, value)
}

func (e *echoContext) Method() string {
	return e.ctx.Request().Method
}

func (e *echoContext) Path() string {
	return e.ctx.Path()
}

func (e *echoContext) Next() {
	// Echo không có `Next()` như Gin; nếu bạn cần chain middleware, hãy dùng echo.MiddlewareFunc
}

func (e *echoContext) Raw() interface{} {
	return e.ctx
}

func (e *echoContext) Set(key string, value interface{}) {
	e.ctx.Set(key, value)
	ctx := context.WithValue(e.ctx.Request().Context(), key, value)
	e.ctx.SetRequest(e.ctx.Request().WithContext(ctx))
}

func (e *echoContext) Get(key string) interface{} {
	if val := e.ctx.Get(key); val != nil {
		return val
	}
	return e.ctx.Request().Context().Value(key)
}

func (e *echoContext) GetString(key string) string {
	val := e.Get(key)
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func (e *echoContext) GetInt(key string) int {
	val := e.Get(key)
	if i, ok := val.(int); ok {
		return i
	}
	return 0
}
