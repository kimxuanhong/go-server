package core

import (
	"context"
)

// Handler defines a generic HTTP handler.
type Handler func(Context)

// RouterGroup defines a generic route group.
type RouterGroup interface {
	Add(method, path string, handler Handler, middleware ...Handler)
}

// RouteConfig defines a route configuration.
type RouteConfig struct {
	Path       string
	Method     string
	Handler    Handler
	Middleware []Handler
}

// Server defines generic server operations.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
	Use(middleware ...Handler)
	AddGroup(relativePath string, register func(rg RouterGroup), middleware ...Handler)
	Add(method, relativePath string, handler Handler, middleware ...Handler)
	RegisterHandlersWithTags(...interface{})
	RegisterHandlers(handlers ...interface{})
	Routes(routes []RouteConfig)
	Static(relativePath, root string)
	HealthCheck()
}
