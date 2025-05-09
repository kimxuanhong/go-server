package fiber

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/kimxuanhong/go-server/core"
	"log"
)

// Config defines server configuration.
type Config struct {
	Host string
	Port string
}

func (c *Config) GetAddr() string {
	return c.Host + ":" + c.Port
}

// Server implements core.Server for Fiber.
type Server struct {
	app    *fiber.App
	config *Config
}

func NewServer(cfg *Config) core.Server {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	return &Server{
		app:    app,
		config: cfg,
	}
}

func (s *Server) Start() error {
	addr := s.config.GetAddr()
	log.Printf("Server is running at %s", addr)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.app.Shutdown()
}

func (s *Server) RegisterMiddleware(middleware ...core.Handler) {
	for _, m := range middleware {
		s.app.Use(transfer(m))
	}
}

func (s *Server) RegisterRoutes(register func(rg core.RouterGroup)) {
	group := s.app.Group("/") // Trả về *fiber.Group
	register(&fiberRouterGroup{group: group})
}

func (s *Server) RegisterPrivateRoutes(register func(rg core.RouterGroup), middleware ...core.Handler) {
	private := s.app.Group("/private") // Trả về *fiber.Group
	for _, m := range middleware {
		private.Use(transfer(m))
	}
	register(&fiberRouterGroup{group: private})
}

func (s *Server) RegisterRoute(method, path string, handler core.Handler) {
	switch method {
	case "GET":
		s.app.Get(path, transfer(handler))
	case "POST":
		s.app.Post(path, transfer(handler))
	case "PUT":
		s.app.Put(path, transfer(handler))
	case "PATCH":
		s.app.Patch(path, transfer(handler))
	case "DELETE":
		s.app.Delete(path, transfer(handler))
	default:
		log.Printf("Unsupported method: %s", method)
	}
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		group := s.app.Group(r.Path)

		// Convert []core.Handler -> []fiber.Handler -> []interface{}
		var middlewareInterfaces []interface{}
		for _, m := range r.Middleware {
			middlewareInterfaces = append(middlewareInterfaces, transfer(m))
		}

		// Đăng ký middleware
		group.Use(middlewareInterfaces...)

		// Đăng ký route handler
		group.Add(r.Method, "/", transfer(r.Handler))
	}
}

type fiberRouterGroup struct {
	group fiber.Router
}

func (g *fiberRouterGroup) Register(method, path string, handler core.Handler) {
	g.group.Add(method, path, transfer(handler))
}

func (g *fiberRouterGroup) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		g.group.Use(transfer(m))
	}
}

func transfer(h core.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return h(&fiberContext{ctx: c})
	}
}
