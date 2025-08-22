package routes

import (
	"ecommerce/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
	cartController *controllers.CartController,
	orderController *controllers.OrderController,
) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",
			"message":      "Server is running",
			"architecture": "MVC Pattern",
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

		// Cart routes
		carts := v1.Group("/carts")
		{
			carts.POST("/", cartController.UpsertCart)
			carts.GET("/:cid/items", cartController.ListItems)
			carts.POST("/:cid/items", cartController.AddItem)
			carts.PATCH("/:cid/items/:iid", cartController.UpdateItem)
			carts.DELETE("/:cid/items/:iid", cartController.RemoveItem)
		}

		// Order routes
		orders := v1.Group("/orders")
		{
			orders.POST("/checkout", orderController.Checkout)
			orders.GET("/", orderController.ListOrders)
		}
	}

	// Legacy routes (for backward compatibility)
	setupLegacyRoutes(r, productController, cartController, orderController)
}

// setupLegacyRoutes maintains backward compatibility with existing API
func setupLegacyRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
	cartController *controllers.CartController,
	orderController *controllers.OrderController,
) {
	// Product routes
	r.GET("/products", productController.ListProducts)
	r.GET("/products/:id", productController.GetProduct)
	r.PUT("/products/:id/stock", productController.UpdateStock)

	// Cart routes
	r.POST("/carts", cartController.UpsertCart)
	r.GET("/carts/:cid/items", cartController.ListItems)
	r.POST("/carts/:cid/items", cartController.AddItem)
	r.PATCH("/carts/:cid/items/:iid", cartController.UpdateItem)
	r.DELETE("/carts/:cid/items/:iid", cartController.RemoveItem)

	// Order routes
	r.POST("/checkout", orderController.Checkout)
	r.GET("/orders", orderController.ListOrders)
}
