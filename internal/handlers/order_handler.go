package handlers

import (
	"ecommerce/internal/models"
	"ecommerce/internal/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	var req models.CheckoutReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad json"})
		return
	}

	orderID, total, err := h.service.ProcessCheckout(req)
	if err != nil {
		if err.Error() == "duplicate" {
			c.JSON(409, gin.H{"error": "duplicate"})
			return
		}
		if err.Error() == "out of stock" {
			c.JSON(422, gin.H{"error": "out of stock"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"order_id": orderID, "total_cents": total})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"items": orders})
}
