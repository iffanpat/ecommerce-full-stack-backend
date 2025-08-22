package services

import (
	"database/sql"
	"ecommerce/internal/db"
	"ecommerce/internal/models"
	"errors"
)

type CartService struct {
	db *sql.DB
}

func NewCartService(database *sql.DB) *CartService {
	return &CartService{db: database}
}

func (s *CartService) UpsertCart(req models.UpsertCartReq) (int, error) {
	if req.UserID == nil && req.GuestToken == nil {
		return 0, errors.New("user_id or guest_token required")
	}
	return db.UpsertCart(s.db, req)
}

func (s *CartService) AddItem(cartID int, req models.AddItemReq) error {
	return db.AddCartItem(s.db, cartID, req)
}

func (s *CartService) UpdateItem(cartID, itemID int, req models.AddItemReq) error {
	return db.UpdateCartItem(s.db, cartID, itemID, req)
}

func (s *CartService) RemoveItem(cartID, itemID int) error {
	return db.RemoveCartItem(s.db, cartID, itemID)
}

func (s *CartService) GetCartItems(cartID int) ([]models.CartItem, error) {
	return db.GetCartItems(s.db, cartID)
}
