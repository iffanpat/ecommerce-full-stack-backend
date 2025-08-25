package views

import (
	"ecommerce/internal/models"
)

// 📋 ProductListResponse - รูปแบบ response สำหรับรายการสินค้า
// ใช้เมื่อส่งรายการสินค้าทั้งหมดกลับไป (GET /products)
type ProductListResponse struct {
	Items []models.Product `json:"items"` // 📦 รายการสินค้าทั้งหมด
	Total int              `json:"total"` // 🔢 จำนวนสินค้าทั้งหมด
}

// 🔍 ProductResponse - รูปแบบ response สำหรับสินค้าเดี่ยว
// ใช้เมื่อส่งข้อมูลสินค้าหนึ่งรายการกลับไป (GET /products/:id)
type ProductResponse struct {
	Product *models.Product `json:"product"` // 🎯 ข้อมูลสินค้า (ใช้ pointer เพื่อประหยัด memory)
}

// 📊 StockUpdateResponse - รูปแบบ response สำหรับการอัปเดตสต็อก
// ใช้เมื่ออัปเดตสต็อกสินค้าสำเร็จ (PUT /products/:id/stock)
type StockUpdateResponse struct {
	ID      int    `json:"id"`      // 🆔 ID ของสินค้าที่อัปเดต
	Stock   int    `json:"stock"`   // 📊 จำนวนสต็อกใหม่
	Message string `json:"message"` // 💬 ข้อความยืนยัน
}

// 📋 FormatProductList - จัดรูปแบบรายการสินค้า
// แปลง []models.Product เป็น ProductListResponse
func FormatProductList(products []models.Product) ProductListResponse {
	// 🛡️ ป้องกัน nil pointer - ถ้า products เป็น nil ให้สร้าง empty slice
	if products == nil {
		products = []models.Product{}
	}

	return ProductListResponse{
		Items: products,      // รายการสินค้า
		Total: len(products), // นับจำนวนสินค้า
	}
}

// 🔍 FormatProduct - จัดรูปแบบสินค้าเดี่ยว
// แปลง *models.Product เป็น ProductResponse
// func FormatProduct(product *models.Product) ProductResponse {
// 	return ProductResponse{
// 		Product: product, // ใส่ข้อมูลสินค้าเข้าไป
// 	}
// }

// 📊 FormatStockUpdate - จัดรูปแบบ response การอัปเดตสต็อก
// สร้าง response เมื่ออัปเดตสต็อกสำเร็จ
func FormatStockUpdate(id, stock int, message string) StockUpdateResponse {
	return StockUpdateResponse{
		ID:      id,      // ID ของสินค้า
		Stock:   stock,   // จำนวนสต็อกใหม่
		Message: message, // ข้อความยืนยัน เช่น "Stock updated successfully"
	}
}

// 💡 ตัวอย่าง JSON response ที่ได้:

// 📋 รายการสินค้า (FormatProductList):
// {
//   "success": true,
//   "data": {
//     "items": [
//       {
//         "id": 1,
//         "sku": "PHONE001",
//         "name": "iPhone 15",
//         "description": "Latest iPhone",
//         "price_cents": 2500000,
//         "currency": "THB",
//         "stock": 10,
//         "image_url": "https://example.com/iphone.jpg"
//       }
//     ],
//     "total": 1
//   }
// }

// 🔍 สินค้าเดี่ยว (FormatProduct):
// {
//   "success": true,
//   "data": {
//     "product": {
//       "id": 1,
//       "sku": "PHONE001",
//       "name": "iPhone 15",
//       ...
//     }
//   }
// }

// 📊 อัปเดตสต็อก (FormatStockUpdate):
// {
//   "success": true,
//   "data": {
//     "id": 1,
//     "stock": 15,
//     "message": "Stock updated successfully"
//   }
// }
