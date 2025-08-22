package controllers

import (
	"log"
	"os"
	"strconv"

	"ecommerce/internal/repositories"
	"ecommerce/internal/views"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	repo repositories.ProductRepository
}

func NewProductController(repo repositories.ProductRepository) *ProductController {
	return &ProductController{repo: repo}
}

func (ctrl *ProductController) ListProducts(c *gin.Context) {
	products, err := ctrl.repo.GetAll()
	if err != nil {
		log.Printf("Error getting products: %v", err)
		views.InternalServerErrorResponse(c, "Failed to retrieve products")
		return
	}

	response := views.FormatProductList(products)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *ProductController) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", idStr)
		views.ValidationErrorResponse(c, "Product ID must be a valid number")
		return
	}

	product, err := ctrl.repo.GetByID(id)
	if err != nil {
		log.Printf("Error getting product ID %d: %v", id, err)
		if err.Error() == "sql: no rows in result set" {
			views.NotFoundResponse(c, "Product")
		} else {
			views.InternalServerErrorResponse(c, "Failed to retrieve product")
		}
		return
	}

	response := views.FormatProduct(product)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *ProductController) UpdateStock(c *gin.Context) {
	// ตรวจสอบ admin secret
	secret := c.GetHeader("X-Admin-Secret")
	if secret != os.Getenv("ADMIN_SECRET") {
		log.Printf("Unauthorized stock update attempt from IP: %s", c.ClientIP())
		views.UnauthorizedResponse(c, "Valid admin secret required")
		return
	}

	// ตรวจสอบ product ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid product ID for stock update: %s", idStr)
		views.ValidationErrorResponse(c, "Product ID must be a valid number")
		return
	}

	// ตรวจสอบ request body
	type StockRequest struct {
		Stock int `json:"stock"`
	}
	var req StockRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Invalid JSON in stock update request: %v", err)
		views.ValidationErrorResponse(c, err.Error())
		return
	}

	// ตรวจสอบค่า stock
	if req.Stock < 0 {
		views.ValidationErrorResponse(c, "Stock cannot be negative")
		return
	}

	// อัปเดต stock
	if err := ctrl.repo.UpdateStock(id, req.Stock); err != nil {
		log.Printf("Error updating stock for product ID %d: %v", id, err)
		views.InternalServerErrorResponse(c, "Failed to update stock")
		return
	}

	log.Printf("Stock updated successfully for product ID %d to %d", id, req.Stock)
	response := views.FormatStockUpdate(id, req.Stock, "Stock updated successfully")
	views.SuccessResponse(c, 200, response)
}
