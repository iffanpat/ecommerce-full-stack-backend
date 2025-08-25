package controllers

import (
	"log"
	"strconv"
	"strings"

	"ecommerce/internal/services"
	"ecommerce/internal/views"

	"github.com/gin-gonic/gin"
)

// ProductController struct - ตัวควบคุมสำหรับจัดการ HTTP requests เกี่ยวกับสินค้า
// Controller เป็นตัวกลางระหว่าง HTTP request และ business logic (ผ่าน Service)
type ProductController struct {
	service services.ProductService // Service สำหรับจัดการ business logic
}

// NewProductController - Factory function สำหรับสร้าง controller
// รับ service เข้ามาเพื่อให้ controller สามารถเรียกใช้ business logic ได้
func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{service: service}
}

// รับ gin.Context ที่มีข้อมูล HTTP request และใช้สำหรับส่ง response
func (ctrl *ProductController) ListProducts(c *gin.Context) {
	log.Println("🎮 Controller: Handling list products request")

	// เรียก Service เพื่อดึงข้อมูลสินค้าทั้งหมด
	// Service จะจัดการ business logic และเรียก Repository
	products, err := ctrl.service.GetAllProducts()
	if err != nil {
		// ถ้าเกิดข้อผิดพลาด ให้ log error และส่ง error response
		log.Printf("❌ Controller: Error getting products: %v", err)
		views.InternalServerErrorResponse(c, "Failed to retrieve products")
		return // หยุดการทำงาน
	}

	// ใช้ View จัดรูปแบบข้อมูลสำหรับ response
	response := views.FormatProductList(products)

	// ส่ง success response กลับไป
	log.Printf("✅ Controller: Successfully returned %d products", len(products))
	views.SuccessResponse(c, 200, response)
}

// GetProduct - Handler สำหรับ GET /products/:id (ดูสินค้าตาม ID)
func (ctrl *ProductController) GetProduct(c *gin.Context) {
	log.Println("🎮 Controller: Handling get product request")

	// Step 1: ดึง parameter "id" จาก URL path
	// ตัวอย่าง: GET /products/123 จะได้ idStr = "123"
	idStr := c.Param("id")
	log.Printf("Controller: Requested product ID: %s", idStr)

	// Step 2: แปลง string เป็น integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// ถ้าแปลงไม่ได้ (เช่น ส่งมา "abc") ให้ส่ง validation error
		log.Printf("❌ Controller: Invalid product ID format: %s", idStr)
		views.ValidationErrorResponse(c, "Product ID must be a valid number")
		return
	}

	// Step 3: เรียก Service เพื่อหาสินค้าด้วย ID
	// Service จะจัดการ business logic, validation และเรียก Repository
	product, err := ctrl.service.GetProductByID(id)
	if err != nil {
		log.Printf("❌ Controller: Error getting product ID %d: %v", id, err)
		// ตรวจสอบประเภท error จาก Service
		if strings.Contains(err.Error(), "not found") {
			// ไม่เจอสินค้า = 404 Not Found
			views.NotFoundResponse(c, "Product")
		} else if strings.Contains(err.Error(), "must be greater than 0") {
			// Validation error = 400 Bad Request
			views.ValidationErrorResponse(c, err.Error())
		} else {
			// error อื่นๆ = 500 Internal Server Error
			views.InternalServerErrorResponse(c, "Failed to retrieve product")
		}
		return
	}

	// Step 4: ส่ง success response
	log.Printf("✅ Controller: Successfully returned product: %s (ID: %d)", product.Name, product.ID)
	views.SuccessResponse(c, 200, product)
}

// UpdateStock - Handler สำหรับ PUT /products/:id/stock (อัปเดตสต็อกสินค้า)
func (ctrl *ProductController) UpdateStock(c *gin.Context) {
	log.Println("🎮 Controller: Handling update stock request")

	// Step 1: ดึง product ID จาก URL parameter
	idStr := c.Param("id")
	log.Printf("📦 Controller: Updating stock for product ID: %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("❌ Controller: Invalid product ID format: %s", idStr)
		views.ValidationErrorResponse(c, "Product ID must be a valid number")
		return
	}

	// Step 2: ดึง admin secret จาก header
	adminSecret := c.GetHeader("X-Admin-Secret")
	if adminSecret == "" {
		log.Printf("🚫 Controller: Missing admin secret header from IP: %s", c.ClientIP())
		views.UnauthorizedResponse(c, "Admin secret header required")
		return
	}

	// Step 3: อ่าน JSON request body
	// สร้าง struct สำหรับรับข้อมูล stock
	type StockRequest struct {
		Stock int `json:"stock"` // จำนวนสต็อกใหม่
	}
	var req StockRequest

	// แปลง JSON จาก request body ไปเป็น struct
	if err := c.BindJSON(&req); err != nil {
		log.Printf("❌ Controller: Invalid JSON in stock update request: %v", err)
		views.ValidationErrorResponse(c, "Invalid JSON format")
		return
	}

	// Step 4: เรียก Service เพื่ออัปเดตสต็อก
	// Service จะจัดการ business logic, validation, authentication และเรียก Repository
	if err := ctrl.service.UpdateProductStock(id, req.Stock, adminSecret); err != nil {
		log.Printf("❌ Controller: Error updating stock: %v", err)

		// ตรวจสอบประเภท error จาก Service
		if strings.Contains(err.Error(), "invalid admin credentials") ||
			strings.Contains(err.Error(), "admin secret not configured") {
			// Authentication error = 401 Unauthorized
			views.UnauthorizedResponse(c, "Invalid admin credentials")
		} else if strings.Contains(err.Error(), "not found") {
			// Product not found = 404 Not Found
			views.NotFoundResponse(c, "Product")
		} else if strings.Contains(err.Error(), "must be greater than 0") ||
			strings.Contains(err.Error(), "cannot be negative") {
			// Validation error = 400 Bad Request
			views.ValidationErrorResponse(c, err.Error())
		} else {
			// error อื่นๆ = 500 Internal Server Error
			views.InternalServerErrorResponse(c, "Failed to update stock")
		}
		return
	}

	// Step 5: อัปเดตสำเร็จ
	log.Printf("✅ Controller: Successfully updated stock for product ID %d to %d", id, req.Stock)

	// จัดรูปแบบ response
	response := views.FormatStockUpdate(id, req.Stock, "Stock updated successfully")

	// ส่ง success response
	views.SuccessResponse(c, 200, response)
}

// หมายเหตุสำคัญ:
// 1. Controller ไม่มี business logic ซับซ้อน เป็นแค่ตัวประสานงาน
// 2. ใช้ Service Layer จัดการ business logic แทน
// 3. ต้อง validate input พื้นฐานก่อนส่งต่อไป Service
// 4. ใช้ Views เพื่อจัดรูปแบบ response ให้สม่ำเสมอ
// 5. Log ทุกขั้นตอนเพื่อง่ายต่อการ debug
// 6. Return เร็วๆ เมื่อเกิด error (Guard Clause pattern)
// 7. Service จะจัดการ error handling และ business validation
