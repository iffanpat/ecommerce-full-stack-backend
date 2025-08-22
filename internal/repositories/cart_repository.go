package repositories

import (
	"database/sql"
	"ecommerce/internal/models"
)

type CartRepository interface {
	UpsertCart(req models.UpsertCartReq) (int, error)
	GetCartItems(cartID int) ([]models.CartItem, error)
	AddItem(cartID int, req models.AddItemReq) error
	UpdateItemQuantity(cartID int, itemID int, qty int) error
	RemoveItem(cartID int, itemID int) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) UpsertCart(req models.UpsertCartReq) (int, error) {
	var cartID int
	var query string
	var args []interface{}

	if req.UserID != nil {
		query = `
			INSERT INTO carts (user_id) VALUES ($1)
			ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id
			RETURNING id
		`
		args = []interface{}{*req.UserID}
	} else if req.GuestToken != nil {
		query = `
			INSERT INTO carts (guest_token) VALUES ($1)
			ON CONFLICT (guest_token) DO UPDATE SET guest_token = EXCLUDED.guest_token
			RETURNING id
		`
		args = []interface{}{*req.GuestToken}
	} else {
		query = `INSERT INTO carts DEFAULT VALUES RETURNING id`
		args = []interface{}{}
	}

	err := r.db.QueryRow(query, args...).Scan(&cartID)
	return cartID, err
}

func (r *cartRepository) GetCartItems(cartID int) ([]models.CartItem, error) {
	query := `
		SELECT ci.id, ci.product_id, p.name, p.price_cents, ci.qty
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
		ORDER BY ci.id
	`

	rows, err := r.db.Query(query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.ProductID, &item.Name, &item.PriceCents, &item.Qty)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *cartRepository) AddItem(cartID int, req models.AddItemReq) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, qty)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty
	`
	_, err := r.db.Exec(query, cartID, req.ProductID, req.Qty)
	return err
}

func (r *cartRepository) UpdateItemQuantity(cartID int, itemID int, qty int) error {
	query := `UPDATE cart_items SET qty = $1 WHERE id = $2 AND cart_id = $3`
	_, err := r.db.Exec(query, qty, itemID, cartID)
	return err
}

func (r *cartRepository) RemoveItem(cartID int, itemID int) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`
	_, err := r.db.Exec(query, itemID, cartID)
	return err
}
