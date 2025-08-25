package services

import (
	"errors"
	"fmt"
	"log"
	"os"

	"ecommerce/internal/models"
	"ecommerce/internal/repositories"
)

// Service Layer คือชั้นกลางระหว่าง Controller และ Repository
// ทำหน้าที่จัดการ business logic, validation, และการประมวลผลข้อมูล
type ProductService interface {
	GetAllProducts() ([]models.Product, error)                      //  ดึงสินค้าทั้งหมด
	GetProductByID(id int) (*models.Product, error)                 //  ดึงสินค้าตาม ID
	UpdateProductStock(id int, stock int, adminSecret string) error //  อัปเดตสต็อกสินค้า
}

// productService struct - ตัวจริงที่ implement ProductService interface
type productService struct {
	repo repositories.ProductRepository // 🗄️ Repository สำหรับเข้าถึงข้อมูล
}

// Factory function สำหรับสร้าง service
func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// GetAllProducts - ดึงสินค้าทั้งหมด (business logic สำหรับการดึงข้อมูลสินค้า)
func (s *productService) GetAllProducts() ([]models.Product, error) {
	log.Println("Service: Starting to retrieve all products")

	// เรียก repository เพื่อดึงข้อมูล
	products, err := s.repo.GetAll()
	if err != nil {
		log.Printf("❌ Service: Failed to retrieve products from repository: %v", err)
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}

	log.Printf("✅ Service: Successfully retrieved %d products", len(products))

	// สามารถเพิ่ม business logic เพิ่มเติมที่นี่ เช่น:
	// - กรองสินค้าที่หมดสต็อก
	// - เพิ่มข้อมูลราคาพิเศษ
	// - คำนวณส่วนลด

	return products, nil
}

// 🔍 GetProductByID - ดึงสินค้าตาม ID พร้อม business logic
func (s *productService) GetProductByID(id int) (*models.Product, error) {
	log.Printf("🔍 Service: Starting to retrieve product with ID: %d", id)

	// ตรวจสอบ ID ว่าถูกต้องหรือไม่
	if id <= 0 {
		log.Printf("❌ Service: Invalid product ID: %d", id)
		return nil, errors.New("product ID must be greater than 0")
	}

	// เรียก repository เพื่อดึงข้อมูล
	product, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("❌ Service: Failed to retrieve product ID %d from repository: %v", id, err)

		// แปลง error จาก repository ให้เป็น business error
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	log.Printf("✅ Service: Successfully retrieved product: %s (ID: %d)", product.Name, product.ID)

	// สามารถเพิ่ม business logic เพิ่มเติมที่นี่ เช่น:
	// - ตรวจสอบว่าสินค้ายังมีขายอยู่หรือไม่
	// - เพิ่มข้อมูลโปรโมชั่น
	// - คำนวณราคาสุทธิ

	return product, nil
}

// 📊 UpdateProductStock - อัปเดตสต็อกสินค้าพร้อม business logic และ security
func (s *productService) UpdateProductStock(id int, stock int, adminSecret string) error {
	log.Printf("📊 Service: Starting stock update for product ID: %d, new stock: %d", id, stock)

	// 🔐 ตรวจสอบ admin authentication
	if err := s.validateAdminAccess(adminSecret); err != nil {
		log.Printf("🚫 Service: Admin authentication failed: %v", err)
		return err
	}

	// ✅ ตรวจสอบ business rules
	if err := s.validateStockUpdate(id, stock); err != nil {
		log.Printf("❌ Service: Stock validation failed: %v", err)
		return err
	}

	// 🔍 ตรวจสอบว่าสินค้ามีอยู่จริงหรือไม่
	existingProduct, err := s.repo.GetByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Printf("❌ Service: Product ID %d not found", id)
			return fmt.Errorf("product with ID %d not found", id)
		}
		log.Printf("❌ Service: Error checking product existence: %v", err)
		return fmt.Errorf("failed to verify product existence: %w", err)
	}

	// 📝 Log การเปลี่ยนแปลง stock
	log.Printf("📈 Service: Stock change for %s (ID: %d): %d → %d",
		existingProduct.Name, id, existingProduct.Stock, stock)

	// เรียก repository เพื่ออัปเดต
	if err := s.repo.UpdateStock(id, stock); err != nil {
		log.Printf("❌ Service: Failed to update stock in repository: %v", err)
		return fmt.Errorf("failed to update stock: %w", err)
	}

	log.Printf("✅ Service: Successfully updated stock for product ID %d to %d", id, stock)

	// สามารถเพิ่ม business logic เพิ่มเติมที่นี่ เช่น:
	// - ส่งการแจ้งเตือนเมื่อสต็อกต่ำ
	// - บันทึก audit log
	// - อัปเดตข้อมูลการขายที่เกี่ยวข้อง

	return nil
}

// 🔐 validateAdminAccess - ตรวจสอบสิทธิ์ admin
func (s *productService) validateAdminAccess(adminSecret string) error {
	expectedSecret := os.Getenv("ADMIN_SECRET")
	if expectedSecret == "" {
		return errors.New("admin secret not configured")
	}

	if adminSecret != expectedSecret {
		return errors.New("invalid admin credentials")
	}

	return nil
}

// ✅ validateStockUpdate - ตรวจสอบ business rules สำหรับการอัปเดตสต็อก
func (s *productService) validateStockUpdate(id int, stock int) error {
	// ตรวจสอบ ID
	if id <= 0 {
		return errors.New("product ID must be greater than 0")
	}

	// ตรวจสอบ stock value
	if stock < 0 {
		return errors.New("stock cannot be negative")
	}

	// สามารถเพิ่ม business rules เพิ่มเติม เช่น:
	// - ตรวจสอบสต็อกสูงสุดที่อนุญาต
	// - ตรวจสอบว่าสินค้ายังขายอยู่หรือไม่
	// - ตรวจสอบสิทธิ์การแก้ไขสินค้าประเภทนี้

	return nil
}

// 💡 หมายเหตุสำคัญ:
// 1. Service Layer ทำหน้าที่จัดการ business logic และ validation
// 2. แยก technical logic (database) ออกจาก business logic
// 3. ทำให้ code ง่ายต่อการทดสอบ (testable)
// 4. Controller จะเรียกใช้ Service แทนการเรียก Repository โดยตรง
// 5. Service สามารถเรียกใช้หลาย Repository ได้ถ้าจำเป็น
