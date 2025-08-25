package repositories

import (
	"database/sql"
	"ecommerce/internal/models"
)

// Interface คือ "สัญญา" ที่กำหนดว่าต้องมี method อะไรบ้าง
type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	UpdateStock(id int, stock int) error
}

// productRepository struct - ตัวจริงที่ทำงานกับฐานข้อมูล
// struct นี้จะ implement ProductRepository interface
type productRepository struct {
	db *sql.DB // การเชื่อมต่อฐานข้อมูล
}

// รับ database connection เข้ามา และ return ProductRepository interface
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// GetAll - ดึงสินค้าทั้งหมดจากฐานข้อมูล
func (r *productRepository) GetAll() ([]models.Product, error) {
	// SQL Query สำหรับดึงข้อมูลสินค้าพร้อมรูปภาพ
	// ใช้ LEFT JOIN เพื่อให้แสดงสินค้าแม้ว่าจะไม่มีรูปก็ได้
	// COALESCE จะ return empty string ถ้า url เป็น NULL
	query := `
		SELECT p.id, p.sku, p.name, p.description, p.price_cents, p.currency, p.stock, 
		       COALESCE(pi.url, '') as image_url
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
		ORDER BY p.id
	`

	// Execute query และรับ rows กลับมา
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err // ถ้า query ผิดพลาด return error
	}
	defer rows.Close() // ปิด rows เมื่อเสร็จสิ้น (สำคัญมาก!)

	// สร้าง slice เพื่อเก็บสินค้าทั้งหมด
	var products []models.Product

	// วนลูปอ่านข้อมูลแต่ละแถว
	for rows.Next() {
		var p models.Product
		// อ่านข้อมูลจาก row ไปใส่ใน struct
		// ลำดับต้องตรงกับ SELECT statement
		err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description,
			&p.PriceCents, &p.Currency, &p.Stock, &p.ImageURL)
		if err != nil {
			return nil, err // ถ้าอ่านข้อมูลผิดพลาด return error
		}
		// เพิ่มสินค้าลงใน slice
		products = append(products, p)
	}

	return products, nil // return สินค้าทั้งหมด
}

// 🔍 GetByID - ดึงสินค้าตาม ID
func (r *productRepository) GetByID(id int) (*models.Product, error) {
	// 📝 SQL Query สำหรับดึงสินค้าเฉพาะ ID
	// ใช้ $1 เป็น placeholder สำหรับ parameter (ป้องกัน SQL injection)
	query := `
		SELECT p.id, p.sku, p.name, p.description, p.price_cents, p.currency, p.stock,
		       COALESCE(pi.url, '') as image_url
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
		WHERE p.id = $1
	`

	var p models.Product
	// 🎯 QueryRow ใช้สำหรับ query ที่คาดว่าจะได้ผลลัพธ์ 1 แถว
	// ส่ง id เป็น parameter ไปแทนที่ $1
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.SKU, &p.Name, &p.Description,
		&p.PriceCents, &p.Currency, &p.Stock, &p.ImageURL)
	if err != nil {
		return nil, err // ถ้าไม่เจอหรือเกิดข้อผิดพลาด return error
	}

	return &p, nil // ✅ return pointer ของสินค้า
}

// 📊 UpdateStock - อัปเดตสต็อกสินค้า
func (r *productRepository) UpdateStock(id int, stock int) error {
	// 📝 SQL Query สำหรับอัปเดตสต็อก
	// ใช้ $1, $2 เป็น placeholder สำหรับ parameters
	query := `UPDATE products SET stock = $1 WHERE id = $2`

	// 🔧 Execute update query
	// Exec ใช้สำหรับ query ที่ไม่ return ข้อมูล (INSERT, UPDATE, DELETE)
	_, err := r.db.Exec(query, stock, id)
	return err // return error (จะเป็น nil ถ้าสำเร็จ)
}

// 💡 หมายเหตุ:
// 1. ใช้ pointer receiver (r *productRepository) เพื่อหลีกเลี่ยงการ copy struct
// 2. ใช้ parameterized query ($1, $2) เพื่อป้องกัน SQL injection
// 3. ใช้ defer rows.Close() เพื่อให้แน่ใจว่า connection จะถูกปิด
// 4. ตรวจสอบ error ทุกครั้งที่มีการเรียกใช้ database operation
