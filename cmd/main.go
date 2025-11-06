package main

import (
	"fmt"
	"net"
	"os"

	"bats.com/local-server/api/handlers"
	"bats.com/local-server/io/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"time"
)

func loadUserDB(userDB *db.UserDB) {
	if err := userDB.LoadUsersFromJSON("sample/users/users.json"); err != nil {
		fmt.Printf("Failed to load users: %v\n", err)
		return
	}
}

func loadFertilizerDB(fertilizerDB *db.FertilizerDB) {
	if err := fertilizerDB.LoadFromJSON("sample/fertilizer/fertilizers_500.json"); err != nil {
		fmt.Printf("Failed to load fertilizers: %v\n", err)
		return
	}
}

func main() {
	args := os.Args
	userDB := db.GetUserDB()
	fertilizerDB := db.GetFertilizerDB()

	// If "refresh" is passed, load JSON data into SQLite
	if len(args) > 1 && args[1] == "refresh" {
		// Case 1: refresh everything
		if len(args) == 2 {
			fmt.Println("Refreshing ALL databases from JSON...")
			loadUserDB(userDB)
			loadFertilizerDB(fertilizerDB)
			fmt.Println("✅ All databases refreshed.")
			return
		}

		// Case 2: refresh specific type
		if len(args) >= 3 {
			switch args[2] {
			case "users":
				fmt.Println("Refreshing USERS database from JSON...")
				loadUserDB(userDB)
				fmt.Println("✅ Users refreshed.")
				return

			case "fertilizers":
				fmt.Println("Refreshing FERTILIZERS database from JSON...")
				loadFertilizerDB(fertilizerDB)
				fmt.Println("✅ Fertilizers refreshed.")
				return

			default:
				fmt.Printf("Unknown refresh target: %s\n", args[2])
				fmt.Println("Valid options: refresh, refresh users, refresh fertilizers")
				return
			}
		}
	}

	// Otherwise, start Fiber app
	fiberApp := newFiberApp() // your Fiber app setup
	setUpRoutes(fiberApp)     // register routes

	listener, port := getListener(9090) // custom listener logic
	fmt.Printf("Server listening on port http://[::]:%d\n", port)

	if err := fiberApp.Listener(listener); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
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
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	handlers.SetUpAuthRoutes(v1)
	handlers.SetUpQuotesRoutes(v1)
	handlers.SetUpUserRoutes(v1)
	handlers.SetUpfertilizerRoutes(v1)
}
