package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/weknora/weknora/internal/server"
)

// @title WeKnora API
// @version 1.0
// @description WeKnora — A collaborative knowledge management platform.
// @termsOfService http://swagger.io/terms/

// @contact.name WeKnora Support
// @contact.url https://github.com/weknora/weknora/issues

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables from .env file if present.
	// In production, environment variables should be set directly.
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] No .env file found, using system environment variables")
	}

	// Retrieve the application port from environment, default to 3000.
	// Changed default from 8080 to 3000 to avoid conflicts with other local services.
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	// Initialize and start the HTTP server.
	srv, err := server.New()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize server: %v", err)
	}

	log.Printf("[INFO] WeKnora server starting on port %s", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatalf("[FATAL] Server encountered a fatal error: %v", err)
	}
}
