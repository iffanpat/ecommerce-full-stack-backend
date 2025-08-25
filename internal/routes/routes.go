package routes

import (
	"ecommerce/internal/controllers"

	"github.com/gin-gonic/gin"
)

// 🛣️ SetupRoutes - กำหนด API routes ทั้งหมด
// ฟังก์ชันนี้จะรับ gin.Engine และ controllers มาเพื่อเซต routing
func SetupRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
) {
	// 🏥 Health check endpoint - ตรวจสอบว่า server ทำงานปกติหรือไม่
	// GET /health จะ return สถานะของ server
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",                 // ✅ server ทำงานปกติ
			"message":      "Server is running",  // 💬 ข้อความ
			"architecture": "MVC Pattern",        // 🏗️ รูปแบบ architecture
			"apis":         []string{"products"}, // 📋 APIs ที่มีให้ใช้
		})
	})

	// 🆕 API version 1 group - จัดกลุ่ม APIs เวอร์ชัน 1
	// ทุก API ในกลุ่มนี้จะมี prefix "/api/v1"
	v1 := r.Group("/api/v1")
	{
		// 📦 Product routes group - จัดกลุ่ม APIs เกี่ยวกับสินค้า
		// ทุก API ในกลุ่มนี้จะมี prefix "/api/v1/products"
		products := v1.Group("/products")
		{
			products.GET("/", productController.ListProducts)

			// GET /api/v1/products/:id - ดูสินค้าตาม ID
			// ตัวอย่าง: GET /api/v1/products/123
			products.GET("/:id", productController.GetProduct)

			// PUT /api/v1/products/:id/stock - อัปเดตสต็อกสินค้า
			// ตัวอย่าง: PUT /api/v1/products/123/stock
			// ต้องมี admin secret ใน header
			products.PUT("/:id/stock", productController.UpdateStock)
		}
	}

	// 🔄 Legacy routes - เก็บ APIs เดิมเพื่อ backward compatibility
	// เรียกฟังก์ชันแยกต่างหากเพื่อความเป็นระเบียบ
	setupLegacyRoutes(r, productController)
}

// 🔄 setupLegacyRoutes - เซต APIs แบบเดิมเพื่อไม่ให้ระบบที่มีอยู่เสีย
// APIs เหล่านี้จะไม่มี version prefix
func setupLegacyRoutes(
	r *gin.Engine,
	productController *controllers.ProductController,
) {
	// 📦 Product routes only - เฉพาะ APIs สินค้า

	// GET /products - ดูสินค้าทั้งหมด (แบบเดิม)
	r.GET("/products", productController.ListProducts)

	// GET /products/:id - ดูสินค้าตาม ID (แบบเดิม)
	// ตัวอย่าง: GET /products/123
	r.GET("/products/:id", productController.GetProduct)

	// PUT /products/:id/stock - อัปเดตสต็อกสินค้า (แบบเดิม)
	// ตัวอย่าง: PUT /products/123/stock
	r.PUT("/products/:id/stock", productController.UpdateStock)
}

// 💡 หมายเหตุ:
// 1. ใช้ r.Group() เพื่อจัดกลุ่ม routes ที่มี prefix เหมือนกัน
// 2. ใช้ parameter ":id" เพื่อรับค่าจาก URL path
// 3. แยก legacy routes ออกมาเพื่อง่ายต่อการจัดการ
// 4. ใช้ Health check endpoint เพื่อ monitoring

// 🛠️ ตัวอย่างการใช้งาน:
//
// Health Check:
// GET http://localhost:8080/health
//
// ดูสินค้าทั้งหมด:
// GET http://localhost:8080/products (legacy)
// GET http://localhost:8080/api/v1/products/ (v1)
//
// ดูสินค้าตาม ID:
// GET http://localhost:8080/products/1 (legacy)
// GET http://localhost:8080/api/v1/products/1 (v1)
//
// อัปเดตสต็อก:
// PUT http://localhost:8080/products/1/stock (legacy)
// PUT http://localhost:8080/api/v1/products/1/stock (v1)
// Headers: X-Admin-Secret: your_secret
// Body: {"stock": 100}
