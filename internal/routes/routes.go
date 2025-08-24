package routes

import (
	"ecommerce/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",
			"message":      "Server is running",
			"architecture": "MVC Pattern",
			"apis":         []string{"products"},
		})
	})

	// API version 1 group
	v1 := r.Group("/api/v1")
	{
		// Product routes
		products := v1.Group("/products")
		{
			products.GET("/", productController.ListProducts)
			products.GET("/:id", productController.GetProduct)
			products.PUT("/:id/stock", productController.UpdateStock)
		}
	}

	// Legacy routes (for backward compatibility)
	setupLegacyRoutes(r, productController)
}

// setupLegacyRoutes maintains backward compatibility with existing API
func setupLegacyRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
) {
	// Product routes only
	r.GET("/products", productController.ListProducts)
	r.GET("/products/:id", productController.GetProduct)
	r.PUT("/products/:id/stock", productController.UpdateStock)
}
