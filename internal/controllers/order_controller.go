package controllers

import (
	"log"

	"ecommerce/internal/models"
	"ecommerce/internal/repositories"
	"ecommerce/internal/views"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	repo repositories.OrderRepository
}

func NewOrderController(repo repositories.OrderRepository) *OrderController {
	return &OrderController{repo: repo}
}

func (ctrl *OrderController) Checkout(c *gin.Context) {
	var req models.CheckoutReq
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Invalid JSON in checkout request: %v", err)
		views.ValidationErrorResponse(c, err.Error())
		return
	}

	// ตรวจสอบ required fields
	if req.CartID <= 0 {
		views.ValidationErrorResponse(c, "Cart ID is required and must be greater than 0")
		return
	}

	if req.IdempotencyKey == "" {
		views.ValidationErrorResponse(c, "Idempotency key is required")
		return
	}

	orderID, total, err := ctrl.repo.ProcessCheckout(req)
	if err != nil {
		log.Printf("Error processing checkout: %v", err)

		switch err.Error() {
		case "duplicate":
			views.ConflictResponse(c, "Duplicate request - order already processed")
		case "out of stock":
			views.UnprocessableEntityResponse(c, "One or more items are out of stock")
		default:
			views.InternalServerErrorResponse(c, "Failed to process checkout")
		}
		return
	}

	response := views.FormatCheckout(orderID, total)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *OrderController) ListOrders(c *gin.Context) {
	orders, err := ctrl.repo.GetAllOrders()
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		views.InternalServerErrorResponse(c, "Failed to retrieve orders")
		return
	}

	response := views.FormatOrderList(orders)
	views.SuccessResponse(c, 200, response)
}
