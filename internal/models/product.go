package models

type Product struct {
	ID          int    `json:"id"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"` //  ราคาในหน่วยเซนต์ (เก็บเป็น int เพื่อหลีกเลี่ยงปัญหา floating point)
	Currency    string `json:"currency"`    //  สกุลเงิน (เช่น "THB", "USD")
	Stock       int    `json:"stock"`
	ImageURL    string `json:"image_url"`
}

//  คำอธิบาย JSON Tags:
// `json:"field_name"` - กำหนดชื่อฟิลด์เมื่อแปลงเป็น JSON
// เมื่อส่ง response กลับไป จะได้ JSON ในรูปแบบนี้:
// {
//   "id": 1,
//   "sku": "PHONE001",
//   "name": "iPhone 15",
//   "description": "Latest iPhone model",
//   "price_cents": 2500000,  // 25,000 บาท (25000 * 100)
//   "currency": "THB",
//   "stock": 10,
//   "image_url": "https://example.com/iphone15.jpg"
// }

//  ทำไมใช้ PriceCents แทน Price?
// เพราะการคำนวณเงินด้วย floating point (เช่น float64) อาจมีปัญหาความแม่นยำ
// เช่น 0.1 + 0.2 ไม่เท่ากับ 0.3 เสมอไป
// การเก็บเป็น cents (หรือ satang ในไทย) เป็น integer จะแม่นยำกว่า
//
// ตัวอย่าง:
// ราคา 25,000 บาท = 2,500,000 สตางค์ = PriceCents: 2500000
