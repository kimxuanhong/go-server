package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kimxuanhong/go-server/core"
	"log"
	"net/http"
	"os"
)

// Config defines server configuration.
type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

func NewConfig() *Config {
	return &Config{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnv("SERVER_PORT", "8080"),
		Mode: getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Config) GetAddr() string {
	return c.Host + ":" + c.Port
}

// Server implements core.Server for Gin.
type Server struct {
	engine     *gin.Engine
	config     *Config
	httpServer *http.Server
}

func NewServer(cfg *Config) core.Server {
	gin.SetMode(cfg.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	return &Server{
		engine: engine,
		config: cfg,
	}
}

func (s *Server) Start() error {
	addr := s.config.GetAddr()
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	log.Printf("Server is running at %s", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		s.engine.Use(transfer(m))
	}
}

func (s *Server) AddGroup(relativePath string, register func(rg core.RouterGroup), middleware ...core.Handler) {
	group := s.engine.Group(relativePath)
	for _, m := range middleware {
		group.Use(transfer(m))
	}
	register(&RouterGroup{group: group})
}

func (s *Server) Add(method, path string, handler core.Handler, middleware ...core.Handler) {
	handlers := make([]gin.HandlerFunc, 0, len(middleware)+1)
	for _, m := range middleware {
		handlers = append(handlers, transfer(m))
	}
	handlers = append(handlers, transfer(handler))
	s.engine.Handle(method, path, handlers...)
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		s.Add(r.Method, r.Path, r.Handler, r.Middleware...)
	}
}

type RouterGroup struct {
	group *gin.RouterGroup
}

func (g *RouterGroup) Add(method, path string, handler core.Handler, middleware ...core.Handler) {
	for _, m := range middleware {
		g.group.Use(transfer(m))
	}
	g.group.Handle(method, path, transfer(handler))
}

func transfer(h core.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&ginContext{ctx: c})
	}
}
