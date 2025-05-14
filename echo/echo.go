package echo

import (
	"context"
	"github.com/kimxuanhong/go-server/core"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"time"
)

// Server implements core.Server for Echo.
type Server struct {
	*core.DynamicRouter
	*core.ProviderRouter
	engine     *echo.Echo
	rootGroup  *echo.Group
	config     *core.Config
	httpServer *http.Server
}

func NewServer(cfg *core.Config) core.Server {
	engine := echo.New()
	engine.HideBanner = true
	engine.Debug = cfg.Mode == "debug"
	engine.Use(middleware.Logger())
	engine.Use(middleware.Recover())
	engine.Pre(middleware.RemoveTrailingSlash()) // Remove trailing /

	rootGroup := engine.Group(cfg.RootPath)

	return &Server{
		DynamicRouter:  &core.DynamicRouter{},
		ProviderRouter: &core.ProviderRouter{},
		engine:         engine,
		rootGroup:      rootGroup,
		config:         cfg,
	}
}

func (s *Server) Start() error {
	addr := s.config.GetAddr()
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	// Add APIs from @Api and provider
	s.LoadRouter()
	s.Routes(s.DynamicRouter.Routes)
	s.Routes(s.ProviderRouter.Routes)

	// Debug: Print registered routes
	for _, r := range s.engine.Routes() {
		log.Printf("Route: %s %s -> %s", r.Method, r.Path, r.Name)
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
		s.rootGroup.Use(transferMiddleware(m))
	}
}

func (s *Server) AddGroup(relativePath string, register func(rg core.RouterGroup), middleware ...core.Handler) {
	group := s.rootGroup.Group(relativePath)
	for _, m := range middleware {
		group.Use(transferMiddleware(m))
	}
	register(&RouterGroup{group: group})
}

func (s *Server) Add(method, relativePath string, handler core.Handler, middleware ...core.Handler) {
	handlers := make([]echo.MiddlewareFunc, 0, len(middleware))
	for _, m := range middleware {
		handlers = append(handlers, transferMiddleware(m))
	}
	s.rootGroup.Add(method, relativePath, transfer(handler))
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		s.Add(r.Method, r.Path, r.Handler, r.Middleware...)
	}
}

func (s *Server) Static(relativePath, root string) {
	s.engine.Static(relativePath, root)
}

func (s *Server) HealthCheck() {
	s.engine.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
	})

	s.engine.GET("/liveness", func(c echo.Context) error {
		return c.JSON(core.StatusOK, map[string]string{"status": "alive"})
	})

	s.engine.GET("/readiness", func(c echo.Context) error {
		// Kiểm tra DB, Redis, etc.
		return c.JSON(core.StatusOK, map[string]string{"status": "ready"})
	})

	s.engine.POST("/terminate", func(c echo.Context) error {
		go func() {
			time.Sleep(1 * time.Second)
			_ = s.Shutdown(context.Background())
		}()
		return c.JSON(core.StatusOK, map[string]string{"status": "terminating"})
	})
}

type RouterGroup struct {
	group *echo.Group
}

func (g *RouterGroup) Add(method, path string, handler core.Handler, middleware ...core.Handler) {
	for _, m := range middleware {
		g.group.Use(transferMiddleware(m))
	}
	g.group.Add(method, path, transfer(handler))
}

func transfer(h core.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		h(&echoContext{ctx: c})
		return nil
	}
}

func transferMiddleware(h core.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// wrap echo.Context thành core.Context
			ctx := &echoContext{ctx: c}

			// Call your handler
			h(ctx)

			// Kiểm tra nếu middleware không gọi abort, thì tiếp tục chuỗi xử lý
			if c.Get("abort") != true {
				return next(c)
			}
			return nil
		}
	}
}
