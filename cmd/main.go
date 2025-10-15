package main

import (
	"fmt"
	"net"
	"os"

	"bats.com/local-server/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"time"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println(args)
	}
	fiberApp := newFiberApp()
	setUpRoutes(fiberApp)
	listener, port := getListener(9090)
	fmt.Printf("Server listening on port http://[::]:%d\n", port)
	err := fiberApp.Listener(listener)
	if err != nil {
		log.Error(err)
		return
	}
}

func getListener(preferredPort int) (net.Listener, int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", preferredPort))
	if err == nil {
		return listener, preferredPort
	}

	fmt.Printf("Port %d is in use, picking a free port...\n", preferredPort)
	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to find a free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return listener, port
}

func newFiberApp() *fiber.App {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        10,              // Maximum number of requests
		Expiration: 1 * time.Second, // Time window for rate limit
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Limit by client IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too Many Requests")
		},
	}))
	return app
}

func setUpRoutes(app *fiber.App) {
	v1 := app.Group("/api/v1")
	handlers.SetUpQuotesRoutes(v1)
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
}
