package handlers

import (
	"database/sql"
	"strconv"

	"ecommerce/internal/models"
	"ecommerce/internal/services"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) UpsertCart(c *gin.Context) {
	var req models.UpsertCartReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad json"})
		return
	}

	cartID, err := h.service.UpsertCart(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"cart_id": cartID})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	cartID, _ := strconv.Atoi(c.Param("cid"))
	var req models.AddItemReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad json"})
		return
	}

	if err := h.service.AddItem(cartID, req); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	cartID, _ := strconv.Atoi(c.Param("cid"))
	itemID, _ := strconv.Atoi(c.Param("iid"))
	var req models.AddItemReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bad json"})
		return
	}

	err := h.service.UpdateItem(cartID, itemID, req)
	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	cartID, _ := strconv.Atoi(c.Param("cid"))
	itemID, _ := strconv.Atoi(c.Param("iid"))

	if err := h.service.RemoveItem(cartID, itemID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *CartHandler) ListItems(c *gin.Context) {
	cartID, _ := strconv.Atoi(c.Param("cid"))
	items, err := h.service.GetCartItems(cartID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"items": items})
}
