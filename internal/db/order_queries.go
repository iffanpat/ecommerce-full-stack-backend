package db

import (
	"database/sql"
	"ecommerce/internal/models"
	"errors"
)

func ProcessCheckout(db *sql.DB, req models.CheckoutReq) (int, int, error) {
	// idempotency check
	if _, err := db.Exec(`INSERT INTO idempotency_keys(key) VALUES($1)`, req.IdempotencyKey); err != nil {
		return 0, 0, errors.New("duplicate")
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	rows, err := tx.Query(`
    SELECT ci.product_id, ci.qty, p.price_cents, p.stock
    FROM cart_items ci JOIN products p ON p.id=ci.product_id
    WHERE ci.cart_id=$1 FOR UPDATE`, req.CartID)
	if err != nil {
		return 0, 0, err
	}

	type line struct{ pid, qty, price int }
	var items []line
	total := 0
	for rows.Next() {
		var pid, qty, price, stock int
		if err = rows.Scan(&pid, &qty, &price, &stock); err != nil {
			return 0, 0, err
		}
		if stock < qty {
			err = errors.New("out of stock")
			return 0, 0, err
		}
		items = append(items, line{pid, qty, price})
		total += price * qty
	}
	rows.Close()

	// Update stock
	for _, it := range items {
		if _, err = tx.Exec(`UPDATE products SET stock = stock - $1 WHERE id=$2`, it.qty, it.pid); err != nil {
			return 0, 0, err
		}
	}

	// Create order
	var orderID int
	if err = tx.QueryRow(`INSERT INTO orders(user_id,total_cents,status,cart_id) VALUES($1,$2,'PAID',$3) RETURNING id`,
		req.UserID, total, req.CartID).Scan(&orderID); err != nil {
		return 0, 0, err
	}

	// Create order items
	for _, it := range items {
		if _, err = tx.Exec(`INSERT INTO order_items(order_id,product_id,price_cents,qty) VALUES($1,$2,$3,$4)`,
			orderID, it.pid, it.price, it.qty); err != nil {
			return 0, 0, err
		}
	}

	// Clear cart
	_, _ = tx.Exec(`DELETE FROM cart_items WHERE cart_id=$1`, req.CartID)
	err = nil
	return orderID, total, nil
}

func GetOrders(db *sql.DB) ([]models.Order, error) {
	rows, err := db.Query(`SELECT id,total_cents,status,created_at FROM orders ORDER BY id DESC LIMIT 50`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		rows.Scan(&order.ID, &order.TotalCents, &order.Status, &order.CreatedAt)
		orders = append(orders, order)
	}
	return orders, nil
}
