package repositories

import (
	"database/sql"
	"ecommerce/internal/models"
	"errors"
)

type OrderRepository interface {
	ProcessCheckout(req models.CheckoutReq) (int, int, error)
	GetAllOrders() ([]models.Order, error)
	CheckIdempotencyKey(key string) (bool, error)
	SaveIdempotencyKey(key string) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CheckIdempotencyKey(key string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM idempotency_keys WHERE key = $1)`
	err := r.db.QueryRow(query, key).Scan(&exists)
	return exists, err
}

func (r *orderRepository) SaveIdempotencyKey(key string) error {
	query := `INSERT INTO idempotency_keys (key) VALUES ($1)`
	_, err := r.db.Exec(query, key)
	return err
}

func (r *orderRepository) ProcessCheckout(req models.CheckoutReq) (int, int, error) {
	// ตรวจสอบ idempotency key
	exists, err := r.CheckIdempotencyKey(req.IdempotencyKey)
	if err != nil {
		return 0, 0, err
	}
	if exists {
		return 0, 0, errors.New("duplicate")
	}

	// เริ่ม transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	// ดึงข้อมูล cart items
	cartItemsQuery := `
		SELECT ci.product_id, ci.qty, p.price_cents, p.stock
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
	`

	rows, err := tx.Query(cartItemsQuery, req.CartID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	type cartItem struct {
		ProductID  int
		Qty        int
		PriceCents int
		Stock      int
	}

	var items []cartItem
	totalCents := 0

	for rows.Next() {
		var item cartItem
		err := rows.Scan(&item.ProductID, &item.Qty, &item.PriceCents, &item.Stock)
		if err != nil {
			return 0, 0, err
		}

		// ตรวจสอบสต็อก
		if item.Stock < item.Qty {
			return 0, 0, errors.New("out of stock")
		}

		totalCents += item.PriceCents * item.Qty
		items = append(items, item)
	}

	// สร้าง order
	var orderID int
	orderQuery := `
		INSERT INTO orders (user_id, total_cents, cart_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = tx.QueryRow(orderQuery, req.UserID, totalCents, req.CartID).Scan(&orderID)
	if err != nil {
		return 0, 0, err
	}

	// สร้าง order items และอัปเดตสต็อก
	for _, item := range items {
		// เพิ่ม order item
		orderItemQuery := `
			INSERT INTO order_items (order_id, product_id, price_cents, qty)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(orderItemQuery, orderID, item.ProductID, item.PriceCents, item.Qty)
		if err != nil {
			return 0, 0, err
		}

		// อัปเดตสต็อก
		updateStockQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2`
		_, err = tx.Exec(updateStockQuery, item.Qty, item.ProductID)
		if err != nil {
			return 0, 0, err
		}
	}

	// บันทึก idempotency key
	_, err = tx.Exec(`INSERT INTO idempotency_keys (key) VALUES ($1)`, req.IdempotencyKey)
	if err != nil {
		return 0, 0, err
	}

	// ลบ cart items
	_, err = tx.Exec(`DELETE FROM cart_items WHERE cart_id = $1`, req.CartID)
	if err != nil {
		return 0, 0, err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, 0, err
	}

	return orderID, totalCents, nil
}

func (r *orderRepository) GetAllOrders() ([]models.Order, error) {
	query := `
		SELECT id, total_cents, status, created_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.TotalCents, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
