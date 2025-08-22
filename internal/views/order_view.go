package views

import (
	"ecommerce/internal/models"
)

// CheckoutResponse represents the response for checkout
type CheckoutResponse struct {
	OrderID    int `json:"order_id"`
	TotalCents int `json:"total_cents"`
}

// OrderListResponse represents the response for order list
type OrderListResponse struct {
	Items []models.Order `json:"items"`
	Total int            `json:"total"`
}

// FormatCheckout formats checkout response
func FormatCheckout(orderID, totalCents int) CheckoutResponse {
	return CheckoutResponse{
		OrderID:    orderID,
		TotalCents: totalCents,
	}
}

// FormatOrderList formats order list response
func FormatOrderList(orders []models.Order) OrderListResponse {
	if orders == nil {
		orders = []models.Order{}
	}

	return OrderListResponse{
		Items: orders,
		Total: len(orders),
	}
}
