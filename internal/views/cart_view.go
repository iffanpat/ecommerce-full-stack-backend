package views

import (
	"ecommerce/internal/models"
)

// CartResponse represents the response for cart creation
type CartResponse struct {
	CartID int `json:"cart_id"`
}

// CartItemsResponse represents the response for cart items
type CartItemsResponse struct {
	Items     []models.CartItem `json:"items"`
	Total     int               `json:"total"`
	TotalCost int               `json:"total_cost"`
}

// ItemActionResponse represents the response for item actions (add, update, remove)
type ItemActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	CartID  int    `json:"cart_id"`
	ItemID  int    `json:"item_id,omitempty"`
}

// FormatCart formats cart creation response
func FormatCart(cartID int) CartResponse {
	return CartResponse{
		CartID: cartID,
	}
}

// FormatCartItems formats cart items response
func FormatCartItems(items []models.CartItem) CartItemsResponse {
	if items == nil {
		items = []models.CartItem{}
	}

	totalCost := 0
	for _, item := range items {
		totalCost += item.PriceCents * item.Qty
	}

	return CartItemsResponse{
		Items:     items,
		Total:     len(items),
		TotalCost: totalCost,
	}
}

// FormatItemAction formats item action response
func FormatItemAction(success bool, message string, cartID, itemID int) ItemActionResponse {
	response := ItemActionResponse{
		Success: success,
		Message: message,
		CartID:  cartID,
	}

	if itemID > 0 {
		response.ItemID = itemID
	}

	return response
}
