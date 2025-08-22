package handlers

import (
	"log"
	"os"
	"strconv"

	"ecommerce/internal/services"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		log.Printf("Error getting products: %v", err)
		c.JSON(500, gin.H{
			"error":   "Failed to get products",
			"details": err.Error(),
			"code":    "PRODUCTS_FETCH_ERROR",
		})
		return
	}
	c.JSON(200, gin.H{"items": products})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", idStr)
		c.JSON(400, gin.H{
			"error":   "Invalid product ID",
			"details": "Product ID must be a valid number",
			"code":    "INVALID_PRODUCT_ID",
		})
		return
	}

	product, err := h.service.GetProduct(id)
	if err != nil {
		log.Printf("Error getting product ID %d: %v", id, err)
		if err.Error() == "sql: no rows in result set" {
			c.JSON(404, gin.H{
				"error":   "Product not found",
				"details": "No product exists with the given ID",
				"code":    "PRODUCT_NOT_FOUND",
				"id":      id,
			})
		} else {
			c.JSON(500, gin.H{
				"error":   "Failed to get product",
				"details": err.Error(),
				"code":    "PRODUCT_FETCH_ERROR",
				"id":      id,
			})
		}
		return
	}
	c.JSON(200, product)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	secret := c.GetHeader("X-Admin-Secret")
	if secret != os.Getenv("ADMIN_SECRET") {
		log.Printf("Unauthorized stock update attempt from IP: %s", c.ClientIP())
		c.JSON(401, gin.H{
			"error":   "Unauthorized",
			"details": "Valid admin secret required",
			"code":    "UNAUTHORIZED_ACCESS",
		})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid product ID for stock update: %s", idStr)
		c.JSON(400, gin.H{
			"error":   "Invalid product ID",
			"details": "Product ID must be a valid number",
			"code":    "INVALID_PRODUCT_ID",
		})
		return
	}

	type StockRequest struct {
		Stock int `json:"stock"`
	}
	var req StockRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Invalid JSON in stock update request: %v", err)
		c.JSON(400, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
			"code":    "INVALID_JSON",
		})
		return
	}

	if req.Stock < 0 {
		c.JSON(400, gin.H{
			"error":   "Invalid stock value",
			"details": "Stock cannot be negative",
			"code":    "INVALID_STOCK_VALUE",
		})
		return
	}

	if err := h.service.UpdateStock(id, req.Stock); err != nil {
		log.Printf("Error updating stock for product ID %d: %v", id, err)
		c.JSON(500, gin.H{
			"error":   "Failed to update stock",
			"details": err.Error(),
			"code":    "STOCK_UPDATE_ERROR",
			"id":      id,
		})
		return
	}

	log.Printf("Stock updated successfully for product ID %d to %d", id, req.Stock)
	c.JSON(200, gin.H{
		"ok":      true,
		"message": "Stock updated successfully",
		"id":      id,
		"stock":   req.Stock,
	})
}
