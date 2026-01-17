package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/minhtranin/ct/internal/db"
	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/logger"
)

func main() {
	defer logger.Logger.Sync()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "CT",
	})

	// Serve static files from the "./static" directory
	app.Static("/static", "./static")

	// Middleware for logging HTTP requests
	app.Use(fiberLog.New(
		fiberLog.Config{
			Format:     "[${time}] ${status} ${latency} - ${method} ${path}\n",
			TimeFormat: "02-Jan-2006 15:04:05",
			TimeZone:   "Local",
		},
	))

	app.Get("/", func(c *fiber.Ctx) error {
		return handler.Home(c)
	})

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	uri := os.Getenv("MONGODB_URI")
	client, err := db.ConnectToMongoDB(uri)
	if err != nil {
		logger.Error("Main", "Error connecting to MongoDB: ")
	}

	defer client.Disconnect(context.Background())

	log.Fatal(app.Listen(":3000"))
}
