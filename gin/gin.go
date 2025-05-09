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
		if h, ok := m.(func(*gin.Context)); ok {
			s.engine.Use(h)
		}
	}
}

func (s *Server) RegisterRoutes(register func(rg core.RouterGroup)) {
	register(&ginRouterGroup{rg: s.engine.Group("/")})
}

func (s *Server) RegisterPrivateRoutes(register func(rg core.RouterGroup), middleware ...core.Handler) {
	private := s.engine.Group("/private")
	for _, m := range middleware {
		if h, ok := m.(func(*gin.Context)); ok {
			private.Use(h)
		}
	}
	register(&ginRouterGroup{rg: private})
}

func (s *Server) RegisterRoute(method, path string, handler core.Handler) {
	if h, ok := handler.(func(*gin.Context)); ok {
		switch method {
		case "GET":
			s.engine.GET(path, h)
		case "POST":
			s.engine.POST(path, h)
		case "PUT":
			s.engine.PUT(path, h)
		case "PATCH":
			s.engine.PATCH(path, h)
		case "DELETE":
			s.engine.DELETE(path, h)
		default:
			log.Printf("Unsupported method: %s", method)
		}
	}
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		group := s.engine.Group(r.Path)
		for _, m := range r.Middleware {
			if h, ok := m.(func(*gin.Context)); ok {
				group.Use(h)
			}
		}
		s.RegisterRoute(r.Method, r.Path, r.Handler)
	}
}

type ginRouterGroup struct {
	rg *gin.RouterGroup
}

func (g *ginRouterGroup) Register(method, path string, handler core.Handler) {
	if h, ok := handler.(func(*gin.Context)); ok {
		g.rg.Handle(method, path, h)
	}
}

func (g *ginRouterGroup) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		if h, ok := m.(func(*gin.Context)); ok {
			g.rg.Use(h)
		}
	}
}
