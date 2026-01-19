package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"github.com/minhtranin/ct/internal/db"
	"github.com/minhtranin/ct/internal/handler"
	"github.com/minhtranin/ct/internal/logger"
	"github.com/minhtranin/ct/internal/router"
)

func main() {
	defer logger.Logger.Sync()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Connect to MongoDB
	uri := os.Getenv("MONGODB_URI")
	client, err := db.ConnectToMongoDB(uri)
	if err != nil {
		logger.Error("Main", "Error connecting to MongoDB", zap.String("error", err.Error()))
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Initialize auth service
	handler.InitAuth(client)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "CT",
		// Increase body limit for larger requests
		BodyLimit: 4 * 1024 * 1024, // 4MB
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

	// Register routes
	router.Register(app)

	log.Fatal(app.Listen(":3000"))
}
