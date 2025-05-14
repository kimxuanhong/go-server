package core

import "context"

type Context interface {
	// Context Request-scoped context (for timeouts, cancellation)
	Context() context.Context

	// Param Input
	Param(name string) string
	Query(name string) string
	Header(name string) string
	Bind(obj interface{}) error

	// JSON Output
	JSON(code int, obj interface{})
	Abort()
	AbortWithStatusJSON(code int, obj interface{})
	String(code int, msg string) error
	Status(code int) Context
	SetHeader(key, value string)

	Method() string
	Path() string
	Next()

	// Raw access if needed
	Raw() interface{}

	Set(key string, value interface{})
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
}
