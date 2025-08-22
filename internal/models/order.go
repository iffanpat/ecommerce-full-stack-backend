package models

type CheckoutReq struct {
	CartID         int    `json:"cart_id"`
	UserID         *int   `json:"user_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

type Order struct {
	ID         int    `json:"id"`
	TotalCents int    `json:"total_cents"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}
