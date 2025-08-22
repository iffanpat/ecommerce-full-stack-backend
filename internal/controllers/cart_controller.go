package controllers

import (
	"log"
	"strconv"

	"ecommerce/internal/models"
	"ecommerce/internal/repositories"
	"ecommerce/internal/views"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	repo repositories.CartRepository
}

func NewCartController(repo repositories.CartRepository) *CartController {
	return &CartController{repo: repo}
}

func (ctrl *CartController) UpsertCart(c *gin.Context) {
	var req models.UpsertCartReq
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Invalid JSON in cart upsert request: %v", err)
		views.ValidationErrorResponse(c, err.Error())
		return
	}

	cartID, err := ctrl.repo.UpsertCart(req)
	if err != nil {
		log.Printf("Error creating/updating cart: %v", err)
		views.InternalServerErrorResponse(c, "Failed to create or update cart")
		return
	}

	response := views.FormatCart(cartID)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *CartController) ListItems(c *gin.Context) {
	cartIDStr := c.Param("cid")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		log.Printf("Invalid cart ID: %s", cartIDStr)
		views.ValidationErrorResponse(c, "Cart ID must be a valid number")
		return
	}

	items, err := ctrl.repo.GetCartItems(cartID)
	if err != nil {
		log.Printf("Error getting cart items for cart ID %d: %v", cartID, err)
		views.InternalServerErrorResponse(c, "Failed to retrieve cart items")
		return
	}

	response := views.FormatCartItems(items)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *CartController) AddItem(c *gin.Context) {
	cartIDStr := c.Param("cid")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		log.Printf("Invalid cart ID: %s", cartIDStr)
		views.ValidationErrorResponse(c, "Cart ID must be a valid number")
		return
	}

	var req models.AddItemReq
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Invalid JSON in add item request: %v", err)
		views.ValidationErrorResponse(c, err.Error())
		return
	}

	// ตรวจสอบค่า quantity
	if req.Qty <= 0 {
		views.ValidationErrorResponse(c, "Quantity must be greater than 0")
		return
	}

	if err := ctrl.repo.AddItem(cartID, req); err != nil {
		log.Printf("Error adding item to cart ID %d: %v", cartID, err)
		views.InternalServerErrorResponse(c, "Failed to add item to cart")
		return
	}

	response := views.FormatItemAction(true, "Item added to cart successfully", cartID, 0)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *CartController) UpdateItem(c *gin.Context) {
	cartIDStr := c.Param("cid")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		log.Printf("Invalid cart ID: %s", cartIDStr)
		views.ValidationErrorResponse(c, "Cart ID must be a valid number")
		return
	}

	itemIDStr := c.Param("iid")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		log.Printf("Invalid item ID: %s", itemIDStr)
		views.ValidationErrorResponse(c, "Item ID must be a valid number")
		return
	}

	type UpdateItemRequest struct {
		Qty int `json:"qty"`
	}
	var req UpdateItemRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Invalid JSON in update item request: %v", err)
		views.ValidationErrorResponse(c, err.Error())
		return
	}

	// ตรวจสอบค่า quantity
	if req.Qty <= 0 {
		views.ValidationErrorResponse(c, "Quantity must be greater than 0")
		return
	}

	if err := ctrl.repo.UpdateItemQuantity(cartID, itemID, req.Qty); err != nil {
		log.Printf("Error updating item quantity for cart ID %d, item ID %d: %v", cartID, itemID, err)
		views.InternalServerErrorResponse(c, "Failed to update item quantity")
		return
	}

	response := views.FormatItemAction(true, "Item quantity updated successfully", cartID, itemID)
	views.SuccessResponse(c, 200, response)
}

func (ctrl *CartController) RemoveItem(c *gin.Context) {
	cartIDStr := c.Param("cid")
	cartID, err := strconv.Atoi(cartIDStr)
	if err != nil {
		log.Printf("Invalid cart ID: %s", cartIDStr)
		views.ValidationErrorResponse(c, "Cart ID must be a valid number")
		return
	}

	itemIDStr := c.Param("iid")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		log.Printf("Invalid item ID: %s", itemIDStr)
		views.ValidationErrorResponse(c, "Item ID must be a valid number")
		return
	}

	if err := ctrl.repo.RemoveItem(cartID, itemID); err != nil {
		log.Printf("Error removing item from cart ID %d, item ID %d: %v", cartID, itemID, err)
		views.InternalServerErrorResponse(c, "Failed to remove item from cart")
		return
	}

	response := views.FormatItemAction(true, "Item removed from cart successfully", cartID, itemID)
	views.SuccessResponse(c, 200, response)
}
