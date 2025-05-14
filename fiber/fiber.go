package fiber

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/kimxuanhong/go-server/core"
	"log"
	"time"
)

// Server implements core.Server for Fiber.
type Server struct {
	*core.DynamicRouter
	*core.ProviderRouter
	app       *fiber.App
	rootGroup fiber.Router
	config    *core.Config
}

func NewServer(cfg *core.Config) core.Server {
	app := fiber.New()
	rootGroup := app.Group(cfg.RootPath)
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	return &Server{
		DynamicRouter:  &core.DynamicRouter{},
		ProviderRouter: &core.ProviderRouter{},
		app:            app,
		rootGroup:      rootGroup,
		config:         cfg,
	}
}

func (s *Server) Start() error {
	addr := s.config.GetAddr()

	//add api from @Api tag
	s.LoadRouter()
	s.Routes(s.DynamicRouter.Routes)

	//add api from provider route
	s.Routes(s.ProviderRouter.Routes)

	// Debug: Print registered routes
	for _, route := range s.app.GetRoutes(true) {
		log.Printf("Route: %s %s -> %s", route.Method, route.Path, route.Name)
	}

	log.Printf("Server is running at %s", addr)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.app.Shutdown()
}

func (s *Server) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		s.rootGroup.Use(transfer(m))
	}
}

func (s *Server) AddGroup(relativePath string, register func(rg core.RouterGroup), middleware ...core.Handler) {
	group := s.rootGroup.Group(relativePath)
	for _, m := range middleware {
		group.Use(transfer(m))
	}
	register(&RouterGroup{group: group})
}

func (s *Server) Add(method, relativePath string, handler core.Handler, middleware ...core.Handler) {
	for _, m := range middleware {
		s.rootGroup.Use(relativePath, transfer(m))
	}
	s.app.Add(method, relativePath, transfer(handler))
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		s.Add(r.Method, r.Path, r.Handler, r.Middleware...)
	}
}

func (s *Server) Static(relativePath, root string) {
	s.app.Static(relativePath, root)
}

func (s *Server) HealthCheck() {
	s.app.Add("GET", "/ping", func(c *fiber.Ctx) error {
		return c.Status(core.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})

	s.app.Add("GET", "/liveness", func(c *fiber.Ctx) error {
		return c.Status(core.StatusOK).JSON(fiber.Map{
			"status": "alive",
		})
	})

	s.app.Add("GET", "/readiness", func(c *fiber.Ctx) error {
		// Thêm logic kiểm tra database, cache, v.v. nếu cần
		return c.Status(core.StatusOK).JSON(fiber.Map{
			"status": "ready",
		})
	})

	s.app.Post("/terminate", func(c *fiber.Ctx) error {
		go func() {
			time.Sleep(1 * time.Second)
			_ = s.Shutdown(context.Background())
		}()
		return c.Status(core.StatusOK).JSON(fiber.Map{
			"status": "terminating",
		})
	})
}

type RouterGroup struct {
	group fiber.Router
}

func (g *RouterGroup) Add(method, path string, handler core.Handler, middleware ...core.Handler) {
	for _, m := range middleware {
		g.group.Use(path, transfer(m))
	}
	g.group.Add(method, path, transfer(handler))
}

func transfer(h core.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h(&fiberContext{ctx: c})
		return nil
	}
}
