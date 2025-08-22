package models

type UpsertCartReq struct {
	UserID     *int    `json:"user_id"`
	GuestToken *string `json:"guest_token"`
}

type AddItemReq struct {
	ProductID int `json:"product_id"`
	Qty       int `json:"qty"`
}

type CartItem struct {
	ID         int    `json:"id"`
	ProductID  int    `json:"product_id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Qty        int    `json:"qty"`
}
