package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kimxuanhong/go-server/core"
	"log"
	"net/http"
)

// Config defines server configuration.
type Config struct {
	Mode string
	Host string
	Port string
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

func (s *Server) RegisterMiddleware(middleware ...core.Handler) {
	for _, m := range middleware {
		s.engine.Use(transfer(m))
	}
}

func (s *Server) RegisterRoutes(register func(rg core.RouterGroup)) {
	register(&ginRouterGroup{group: s.engine.Group("/")})
}

func (s *Server) RegisterPrivateRoutes(register func(rg core.RouterGroup), middleware ...core.Handler) {
	private := s.engine.Group("/private") // Trả về *gin.RouterGroup
	for _, m := range middleware {
		private.Use(transfer(m))
	}
	register(&ginRouterGroup{group: private})
}

func (s *Server) RegisterRoute(method, path string, handler core.Handler) {
	switch method {
	case "GET":
		s.engine.GET(path, transfer(handler))
	case "POST":
		s.engine.POST(path, transfer(handler))
	case "PUT":
		s.engine.PUT(path, transfer(handler))
	case "PATCH":
		s.engine.PATCH(path, transfer(handler))
	case "DELETE":
		s.engine.DELETE(path, transfer(handler))
	default:
		log.Printf("Unsupported method: %s", method)
	}
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		group := s.engine.Group(r.Path)
		var middlewareInterfaces []gin.HandlerFunc
		for _, m := range r.Middleware {
			middlewareInterfaces = append(middlewareInterfaces, transfer(m))
		}
		group.Use(middlewareInterfaces...)
		group.Handle(r.Method, "/", transfer(r.Handler))
	}
}

type ginRouterGroup struct {
	group *gin.RouterGroup
}

func (g *ginRouterGroup) Register(method, path string, handler core.Handler) {
	g.group.Handle(method, path, transfer(handler))
}

func (g *ginRouterGroup) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		g.group.Use(transfer(m))
	}
}

func transfer(h core.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(&ginContext{ctx: c}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
	}
}
