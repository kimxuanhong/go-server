package main

import (
	"github.com/kimxuanhong/go-server/core"
	"github.com/kimxuanhong/go-server/example/gin/internal/api"
	"github.com/kimxuanhong/go-server/gin"
	"log"
	"net/http"
)

func main() {
	// Tạo cấu hình cho server
	cfg := &core.Config{
		Host:     "localhost",
		Port:     "8081",
		Mode:     "debug",
		RootPath: "/api/v1/",
	}
	// Khởi tạo server với cấu hình
	server := gin.NewServer(cfg)

	// Đăng ký middleware toàn cục
	server.Use(func(c core.Context) {
		log.Printf("Middleware: %s %s", c.Method(), c.Path())
		c.Next()
	})

	// Đăng ký các routes
	server.AddGroup("/index", func(rg core.RouterGroup) {
		// Đăng ký một route đơn giản
		rg.Add("GET", "/hello", func(c core.Context) {
			c.JSON(200, map[string]string{
				"message": "Hello, world!",
			})
		}, func(c core.Context) {
			log.Println("Test /hello")
			c.Next()
		})
	}, func(c core.Context) {
		log.Println("Test /index")
		c.Next()
	})

	server.Add("GET", "/ping", func(c core.Context) {
		c.JSON(200, map[string]string{
			"message": "pong",
		})
	}, func(c core.Context) {
		log.Println("Test /ping")
		c.Next()
	})

	funcHandler := []core.RouteConfig{
		{
			Path:   "/users/index",
			Method: http.MethodGet,
			Handler: func(c core.Context) {
				c.JSON(200, map[string]string{
					"message": "/users/index",
				})
			},
			Middleware: []core.Handler{func(c core.Context) {
				log.Println("Test /users/index")
				c.Next()
			}},
		},
	}
	server.Routes(funcHandler)

	server.SetHandlers(&api.MyApiHandler{})
	server.HealthCheck()

	// Bắt đầu chạy server
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
