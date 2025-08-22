package main

import (
	"log"
	"os"

	"ecommerce/internal/db"
	"ecommerce/internal/handlers"
	"ecommerce/internal/services"

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

	// Initialize services
	productService := services.NewProductService(database)
	cartService := services.NewCartService(database)
	orderService := services.NewOrderService(database)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "X-Admin-Secret"},
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// Product routes
	r.GET("/products", productHandler.ListProducts)
	r.GET("/products/:id", productHandler.GetProduct)
	r.PUT("/products/:id/stock", productHandler.UpdateStock)

	// Cart routes
	r.POST("/carts", cartHandler.UpsertCart)
	r.GET("/carts/:cid/items", cartHandler.ListItems)
	r.POST("/carts/:cid/items", cartHandler.AddItem)
	r.PATCH("/carts/:cid/items/:iid", cartHandler.UpdateItem)
	r.DELETE("/carts/:cid/items/:iid", cartHandler.RemoveItem)

	// Order routes
	r.POST("/checkout", orderHandler.Checkout)
	r.GET("/orders", orderHandler.ListOrders)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
