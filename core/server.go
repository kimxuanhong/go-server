package core

import (
	"context"
)

// Handler defines a generic HTTP handler.
type Handler func(Context) error

type Config interface {
	GetAddr() string
}

// RouterGroup defines a generic route group.
type RouterGroup interface {
	Register(method, path string, handler Handler)
	Use(middleware ...Handler)
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
	RegisterMiddleware(middleware ...Handler)
	RegisterRoutes(register func(rg RouterGroup))
	RegisterPrivateRoutes(register func(rg RouterGroup), middleware ...Handler)
	RegisterRoute(method, path string, handler Handler)
	Routes(routes []RouteConfig)
}
