package views

import (
	"ecommerce/internal/models"
)

// ProductListResponse represents the response for product list
type ProductListResponse struct {
	Items []models.Product `json:"items"`
	Total int              `json:"total"`
}

// ProductResponse represents the response for a single product
type ProductResponse struct {
	Product *models.Product `json:"product"`
}

// StockUpdateResponse represents the response for stock update
type StockUpdateResponse struct {
	ID      int    `json:"id"`
	Stock   int    `json:"stock"`
	Message string `json:"message"`
}

// FormatProductList formats product list response
func FormatProductList(products []models.Product) ProductListResponse {
	if products == nil {
		products = []models.Product{}
	}

	return ProductListResponse{
		Items: products,
		Total: len(products),
	}
}

// FormatProduct formats single product response
func FormatProduct(product *models.Product) ProductResponse {
	return ProductResponse{
		Product: product,
	}
}

// FormatStockUpdate formats stock update response
func FormatStockUpdate(id, stock int, message string) StockUpdateResponse {
	return StockUpdateResponse{
		ID:      id,
		Stock:   stock,
		Message: message,
	}
}
