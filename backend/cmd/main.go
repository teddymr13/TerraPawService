package main

import (
	"log"
	"os"

	"github.com/TerraPaw/backend/db"
	"github.com/TerraPaw/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize database
	db.InitDB()

	// Patch Dummy Data (4000 records)
	db.PatchLargeData()
	// Ensure Food Data exists (if skipped by PatchLargeData)
	db.EnsureFoodData()
	// Ensure Categories exist
	db.SeedCategories()

	// Seed data if empty
	db.SeedData()

	// Create Gin router
	router := gin.Default()

	// 404 Handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status":  "error",
			"message": "Resource not found",
			"code":    404,
		})
	})

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Setup CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Register routes
	routes.SetupRoutes(router)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
