package models

type Product struct {
	ID          int    `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"`
	Currency    string `json:"currency"`
	Stock       int    `json:"stock"`
	ImageURL    string `json:"image_url"`
}
