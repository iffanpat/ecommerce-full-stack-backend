package repositories

import (
	"database/sql"
	"ecommerce/internal/models"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	UpdateStock(id int, stock int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]models.Product, error) {
	query := `
		SELECT p.id, p.sku, p.name, p.description, p.price_cents, p.currency, p.stock, 
		       COALESCE(pi.url, '') as image_url
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
		ORDER BY p.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description,
			&p.PriceCents, &p.Currency, &p.Stock, &p.ImageURL)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *productRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT p.id, p.sku, p.name, p.description, p.price_cents, p.currency, p.stock,
		       COALESCE(pi.url, '') as image_url
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
		WHERE p.id = $1
	`

	var p models.Product
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Description,
		&p.PriceCents, &p.Currency, &p.Stock, &p.ImageURL)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) UpdateStock(id int, stock int) error {
	query := `UPDATE products SET stock = $1 WHERE id = $2`
	_, err := r.db.Exec(query, stock, id)
	return err
}
