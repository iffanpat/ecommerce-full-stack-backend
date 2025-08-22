package db

import (
	"database/sql"
	"ecommerce/internal/models"
)

func UpsertCart(db *sql.DB, req models.UpsertCartReq) (int, error) {
	var id int
	if req.GuestToken != nil {
		// insert or no-op if exists
		if err := db.QueryRow(`
      INSERT INTO carts(guest_token) VALUES($1)
      ON CONFLICT (guest_token) DO UPDATE SET guest_token=EXCLUDED.guest_token
      RETURNING id`, *req.GuestToken).Scan(&id); err != nil {
			// try fetch existing
			_ = db.QueryRow(`SELECT id FROM carts WHERE guest_token=$1`, *req.GuestToken).Scan(&id)
		}
	} else if req.UserID != nil {
		if err := db.QueryRow(`
      INSERT INTO carts(user_id) VALUES($1)
      ON CONFLICT DO NOTHING
      RETURNING id`, *req.UserID).Scan(&id); err != nil {
			_ = db.QueryRow(`SELECT id FROM carts WHERE user_id=$1 ORDER BY id DESC LIMIT 1`, *req.UserID).Scan(&id)
		}
	}
	return id, nil
}

func AddCartItem(db *sql.DB, cartID int, req models.AddItemReq) error {
	qty := req.Qty
	if qty <= 0 {
		qty = 1
	}
	_, err := db.Exec(`INSERT INTO cart_items(cart_id, product_id, qty)
    VALUES($1,$2,$3)
    ON CONFLICT (cart_id, product_id) DO UPDATE SET qty = cart_items.qty + EXCLUDED.qty`, cartID, req.ProductID, qty)
	return err
}

func UpdateCartItem(db *sql.DB, cartID, itemID int, req models.AddItemReq) error {
	res, err := db.Exec(`UPDATE cart_items SET qty=$1 WHERE id=$2 AND cart_id=$3`, req.Qty, itemID, cartID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func RemoveCartItem(db *sql.DB, cartID, itemID int) error {
	_, err := db.Exec(`DELETE FROM cart_items WHERE id=$1 AND cart_id=$2`, itemID, cartID)
	return err
}

func GetCartItems(db *sql.DB, cartID int) ([]models.CartItem, error) {
	rows, err := db.Query(`
    SELECT ci.id, ci.qty, p.id as product_id, p.name, p.price_cents
    FROM cart_items ci
    JOIN products p ON p.id = ci.product_id
    WHERE ci.cart_id=$1
    ORDER BY ci.id`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(&item.ID, &item.Qty, &item.ProductID, &item.Name, &item.PriceCents); err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}
