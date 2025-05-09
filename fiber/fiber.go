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
		if h, ok := m.(func(*fiber.Ctx) error); ok {
			s.app.Use(h)
		} else {
			log.Printf("Invalid middleware type for Fiber")
		}
	}
}

func (s *Server) RegisterRoutes(register func(rg core.RouterGroup)) {
	group := s.app.Group("/") // Trả về *fiber.Group
	register(&fiberRouterGroup{group: group})
}

func (s *Server) RegisterPrivateRoutes(register func(rg core.RouterGroup), middleware ...core.Handler) {
	private := s.app.Group("/private") // Trả về *fiber.Group
	for _, m := range middleware {
		if h, ok := m.(func(*fiber.Ctx) error); ok {
			private.Use(h)
		} else {
			log.Printf("Invalid middleware type for Fiber")
		}
	}
	register(&fiberRouterGroup{group: private})
}

func (s *Server) RegisterRoute(method, path string, handler core.Handler) {
	if h, ok := handler.(func(*fiber.Ctx) error); ok {
		s.app.Add(method, path, h)
	} else {
		log.Printf("Invalid handler type for Fiber")
	}
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		group := s.app.Group(r.Path)
		fiberMiddleware := make([]interface{}, 0, len(r.Middleware))
		for _, m := range r.Middleware {
			if h, ok := m.(func(*fiber.Ctx) error); ok {
				fiberMiddleware = append(fiberMiddleware, h)
			} else {
				log.Printf("Invalid middleware type for Fiber in route %s", r.Path)
			}
		}
		group.Use(fiberMiddleware...)
		// Register the route handler
		s.RegisterRoute(r.Method, r.Path, r.Handler)
	}
}

type fiberRouterGroup struct {
	group fiber.Router // Sửa: dùng *fiber.Group thay vì *fiber.App
}

func (g *fiberRouterGroup) Register(method, path string, handler core.Handler) {
	if h, ok := handler.(func(*fiber.Ctx) error); ok {
		g.group.Add(method, path, h)
	} else {
		log.Printf("Invalid handler type for Fiber")
	}
}

func (g *fiberRouterGroup) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		if h, ok := m.(func(*fiber.Ctx) error); ok {
			g.group.Use(h)
		} else {
			log.Printf("Invalid middleware type for Fiber")
		}
	}
}
