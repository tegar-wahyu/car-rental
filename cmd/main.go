package main

import (
	"car-rental/pkg/database"
	"car-rental/pkg/routes"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	database.Connect()
	database.Migrate()
	database.SeedData()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS to allow all methods
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Configure trusted proxies
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies == "" {
		// If not specified, only trust localhost/loopback addresses
		r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	} else if strings.ToLower(trustedProxies) == "nil" || strings.ToLower(trustedProxies) == "none" {
		// Disable trusted proxies completely
		r.SetTrustedProxies(nil)
	} else {
		// Use the provided list of trusted proxies
		r.SetTrustedProxies(strings.Split(trustedProxies, ","))
	}

	// Setup routes
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
