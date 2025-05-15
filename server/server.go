package server

import (
	"github.com/kimxuanhong/go-server/core"
	"github.com/kimxuanhong/go-server/echo"
	"github.com/kimxuanhong/go-server/fiber"
	"github.com/kimxuanhong/go-server/gin"
	"log"
)

func NewServer(configs ...*core.Config) core.Server {
	cfg := core.GetConfig(configs...)
	switch cfg.Engine {
	case "gin":
		return gin.NewServer(cfg)
	case "fiber":
		return fiber.NewServer(cfg)
	case "echo":
		return echo.NewServer(cfg)
	default:
		log.Fatalf("Can not init server with engine = %s", cfg.Engine)
	}
	return nil
}
