package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kimxuanhong/go-server/core"
	"log"
	"net/http"
	"time"
)

// Server implements core.Server for Gin.
type Server struct {
	*core.DynamicRouter
	*core.ProviderRouter
	engine     *gin.Engine
	rootGroup  *gin.RouterGroup
	config     *core.Config
	httpServer *http.Server
}

func NewServer(cfg *core.Config) core.Server {
	gin.SetMode(cfg.Mode)
	engine := gin.New()
	rootGroup := engine.Group(cfg.RootPath)
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

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

	//add api from @Api tag
	s.LoadRouter()
	s.Routes(s.DynamicRouter.Routes)

	//add api from provider route
	s.Routes(s.ProviderRouter.Routes)

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
	handlers := make([]gin.HandlerFunc, 0, len(middleware)+1)
	for _, m := range middleware {
		handlers = append(handlers, transfer(m))
	}
	handlers = append(handlers, transfer(handler))
	s.rootGroup.Handle(method, relativePath, handlers...)
}

func (s *Server) Routes(routes []core.RouteConfig) {
	for _, r := range routes {
		s.Add(r.Method, r.Path, r.Handler, r.Middleware...)
	}
}

func (s *Server) Static(relativePath, root string) {
	s.engine.Static(relativePath, root)
}

func (s *Server) RootPath(relativePath string) {
	if relativePath != "" {
		s.config.RootPath = relativePath
	}
}

func (s *Server) HealthCheck() {
	s.engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	s.engine.GET("/liveness", func(c *gin.Context) {
		c.JSON(core.StatusOK, gin.H{"status": "alive"})
	})

	s.engine.GET("/readiness", func(c *gin.Context) {
		// Bạn có thể kiểm tra kết nối DB, Redis, etc. tại đây
		c.JSON(core.StatusOK, gin.H{"status": "ready"})
	})

	s.engine.POST("/terminate", func(c *gin.Context) {
		go func() {
			time.Sleep(1 * time.Second)
			_ = s.Shutdown(context.Background())
		}()
		c.JSON(core.StatusOK, gin.H{"status": "terminating"})
	})
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
