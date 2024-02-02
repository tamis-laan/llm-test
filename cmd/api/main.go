package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"llm/internal/clients"
	_ "llm/internal/clients"
	"llm/internal/routes"
	"llm/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/golang-migrate/migrate/v4"
)

// Embed documentation folder into binary
var (
	//go:embed docs
	docs embed.FS
)

// Server
func main() {

	// Load configuration
	config := utils.GetConfig()

	// Get the database connection
	pool := clients.GetPool()

	// Close PostgreSQL connection
	defer pool.Close()

	// Migration flag
	migrateFlag := flag.Bool("migrate", false, "Auto migrate before starting server")

	// Parse arguments
	flag.Parse()

	// Get the migration client
	migration := clients.GetMigration()

	// Auto migrate to latest version
	if *migrateFlag {

		// Log to console
		log.Println("Migrate database to latest version")

		// Start migration
		err := migration.Up()

		// Check for migration errors
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	}

	// Get the migration version
	version, dirty, _ := migration.Version()

	// Close migration connection
	migration.Close()

	// Exit on dirty database
	if dirty {
		log.Fatal("Database is dirty!")
	}

	// Log migration version
	log.Printf("Database version %d", version)

	// Create app
	app := fiber.New(fiber.Config{
		AppName: "GIO API",
	})

	// Add logging
	app.Use(logger.New())

	// Server 500 error handler
	app.Use(recover.New())

	// Compress communications
	app.Use(compress.New())

	// Set the LLM routes
	routes.SetRouteLLM(app.Group("/llm"))

	// Health route
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check if we can reach the database
		err := pool.Ping(context.Background())

		// Send the error on failure
		if err != nil {
			log.Println("Cannot connect to database")
			return c.Status(500).SendString("Cannot connect to database")
		}

		// Everything is oké
		return c.SendStatus(200)
	})

	// Test route
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("LLM API - go to /docs")
	})

	// Panic route
	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("PANIC!!")
	})

	// Error route
	app.Get("/error", func(c *fiber.Ctx) error {
		// Return 503 error
		return fiber.ErrServiceUnavailable
	})

	// Get the pod hostname
	hostname, _ := os.Hostname()

	// Metrics route with dashboard
	app.Get("/metrics", monitor.New(monitor.Config{Title: hostname}))

	// Host docs
	app.Use("/docs", filesystem.New(filesystem.Config{
		Root:       http.FS(docs),
		PathPrefix: "docs",
		Browse:     true,
	}))

	// Server connection string
	con := fmt.Sprintf(":%d", config.Server.Port)

	// Start web server
	app.Listen(con)

}
