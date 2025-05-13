package fiber

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/kimxuanhong/go-server/core"
	"log"
	"path"
)

// Server implements core.Server for Fiber.
type Server struct {
	*core.DynamicRouter
	*core.ProviderRouter
	app    *fiber.App
	config *core.Config
}

func NewServer(cfg *core.Config) core.Server {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	return &Server{
		DynamicRouter:  &core.DynamicRouter{},
		ProviderRouter: &core.ProviderRouter{},
		app:            app,
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

	log.Printf("Server is running at %s", addr)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.app.Shutdown()
}

func (s *Server) Use(middleware ...core.Handler) {
	for _, m := range middleware {
		s.app.Use(transfer(m))
	}
}

func (s *Server) AddGroup(relativePath string, register func(rg core.RouterGroup), middleware ...core.Handler) {
	group := s.app.Group(path.Join(s.config.RootPath, relativePath))
	for _, m := range middleware {
		group.Use(transfer(m))
	}
	register(&RouterGroup{group: group})
}

func (s *Server) Add(method, relativePath string, handler core.Handler, middleware ...core.Handler) {
	for _, m := range middleware {
		s.app.Use(path.Join(s.config.RootPath, relativePath), transfer(m))
	}
	s.app.Add(method, path.Join(s.config.RootPath, relativePath), transfer(handler))
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		s.Add(r.Method, r.Path, r.Handler, r.Middleware...)
	}
}

func (s *Server) Static(relativePath, root string) {
	s.app.Static(relativePath, root)
}

func (s *Server) RootPath(relativePath string) {
	if relativePath != "" {
		s.config.RootPath = relativePath
	}
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
