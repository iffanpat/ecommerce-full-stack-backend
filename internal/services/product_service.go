package services

import (
	"database/sql"
	"ecommerce/internal/db"
	"ecommerce/internal/models"
)

type ProductService struct {
	db *sql.DB
}

func NewProductService(database *sql.DB) *ProductService {
	return &ProductService{db: database}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return db.GetProducts(s.db)
}

func (s *ProductService) GetProduct(id int) (*models.Product, error) {
	return db.GetProductByID(s.db, id)
}

func (s *ProductService) UpdateStock(id int, stock int) error {
	return db.UpdateProductStock(s.db, id, stock)
}
