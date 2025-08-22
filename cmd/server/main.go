package main

import (
	"log"
	"os"

	"ecommerce/internal/controllers"
	"ecommerce/internal/db"
	"ecommerce/internal/repositories"
	"ecommerce/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories (Data Access Layer)
	productRepo := repositories.NewProductRepository(database)
	cartRepo := repositories.NewCartRepository(database)
	orderRepo := repositories.NewOrderRepository(database)

	// Initialize controllers (MVC Controllers)
	productController := controllers.NewProductController(productRepo)
	cartController := controllers.NewCartController(cartRepo)
	orderController := controllers.NewOrderController(orderRepo)

	// Setup router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "X-Admin-Secret"},
	}))

	// Setup routes using MVC pattern
	routes.SetupRoutes(r, productController, cartController, orderController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📁 Using MVC architecture pattern")
	log.Printf("🔗 API endpoints available at:")
	log.Printf("   - Health check: http://localhost:%s/health", port)
	log.Printf("   - Products: http://localhost:%s/products", port)
	log.Printf("   - New API v1: http://localhost:%s/api/v1/", port)

	r.Run(":" + port)
}
