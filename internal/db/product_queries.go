package db

import (
	"database/sql"
	"ecommerce/internal/models"
)

func GetProducts(db *sql.DB) ([]models.Product, error) {
	rows, err := db.Query(`
    SELECT p.id, p.sku, p.name, COALESCE(p.description,''), p.price_cents, p.currency,
           p.stock,
           COALESCE((SELECT url FROM product_images WHERE product_id=p.id AND is_primary=true LIMIT 1),'') as image_url
    FROM products p
    ORDER BY p.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.Stock, &p.ImageURL); err == nil {
			products = append(products, p)
		}
	}
	return products, nil
}

func GetProductByID(db *sql.DB, id int) (*models.Product, error) {
	var p models.Product
	err := db.QueryRow(`
    SELECT p.id, p.sku, p.name, COALESCE(p.description,''), p.price_cents, p.currency, p.stock,
           COALESCE((SELECT url FROM product_images WHERE product_id=p.id AND is_primary=true LIMIT 1),'')
    FROM products p WHERE p.id=$1`, id).
		Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.Stock, &p.ImageURL)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProductStock(db *sql.DB, id int, stock int) error {
	_, err := db.Exec(`UPDATE products SET stock=$1 WHERE id=$2`, stock, id)
	return err
}
