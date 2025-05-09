package main

import (
	"github.com/kimxuanhong/go-server/core"
	"github.com/kimxuanhong/go-server/fiber"
	"log"
)

func main() {
	// Tạo cấu hình cho server
	cfg := &fiber.Config{
		Host: "localhost",
		Port: "8080",
	}
	// Khởi tạo server với cấu hình
	server := fiber.NewServer(cfg)

	// Đăng ký middleware toàn cục
	server.RegisterMiddleware(func(c core.Context) error {
		log.Printf("Middleware: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	// Đăng ký các routes
	server.RegisterRoutes(func(rg core.RouterGroup) {
		// Đăng ký một route đơn giản
		rg.Register("GET", "/hello", func(c core.Context) error {
			c.JSON(200, map[string]string{
				"message": "Hello, world!",
			})
			return nil
		})
	})

	// Đăng ký các routes riêng tư
	server.RegisterPrivateRoutes(func(rg core.RouterGroup) {
		// Đăng ký một route cho nhóm private
		rg.Register("GET", "/private", func(c core.Context) error {
			c.JSON(200, map[string]string{
				"message": "Private route accessed",
			})
			return nil
		})
	}, func(c core.Context) error {
		// Middleware cho private routes
		log.Println("Private route accessed")
		return c.Next()
	})

	server.RegisterRoute("GET", "/ping", func(c core.Context) error {
		c.JSON(200, map[string]string{
			"message": "Hello, world!",
		})
		return nil
	})

	// Bắt đầu chạy server
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
