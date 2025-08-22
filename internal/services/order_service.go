package services

import (
	"database/sql"
	"ecommerce/internal/db"
	"ecommerce/internal/models"
)

type OrderService struct {
	db *sql.DB
}

func NewOrderService(database *sql.DB) *OrderService {
	return &OrderService{db: database}
}

func (s *OrderService) ProcessCheckout(req models.CheckoutReq) (int, int, error) {
	return db.ProcessCheckout(s.db, req)
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return db.GetOrders(s.db)
}
