package main

import (
	"github.com/kimxuanhong/go-server/core"
	"github.com/kimxuanhong/go-server/gin"
	"log"
)

func main() {
	// Tạo cấu hình cho server
	cfg := &gin.Config{
		Host: "localhost",
		Port: "8081",
		Mode: "debug",
	}
	// Khởi tạo server với cấu hình
	server := gin.NewServer(cfg)

	// Đăng ký middleware toàn cục
	server.RegisterMiddleware(func(c core.Context) {
		log.Printf("Middleware: %s %s", c.Method(), c.Path())
		c.Next()
	})

	// Đăng ký các routes
	server.RegisterRoutes(func(rg core.RouterGroup) {
		// Đăng ký một route đơn giản
		rg.Register("GET", "/hello", func(c core.Context) {
			c.JSON(200, map[string]string{
				"message": "Hello, world!",
			})
		})
	})

	// Đăng ký các routes riêng tư
	server.RegisterPrivateRoutes(func(rg core.RouterGroup) {
		// Đăng ký một route cho nhóm private
		rg.Register("GET", "/private", func(c core.Context) {
			c.JSON(200, map[string]string{
				"message": "Private route accessed",
			})
		})
	}, func(c core.Context) {
		// Middleware cho private routes
		log.Println("Private route accessed")
		c.Next()
	})

	server.RegisterRoute("GET", "/ping", func(c core.Context) {
		c.JSON(200, map[string]string{
			"message": "Hello, world!",
		})
	})

	// Bắt đầu chạy server
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
