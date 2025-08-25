package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 🎨 StandardResponse - โครงสร้าง JSON response มาตรฐาน
// ทุก API response จะใช้รูปแบบนี้เพื่อความสม่ำเสมอ
type StandardResponse struct {
	Success bool        `json:"success"`         // ✅ true = สำเร็จ, false = ผิดพลาด
	Data    interface{} `json:"data,omitempty"`  // 📦 ข้อมูลที่ต้องการส่งกลับ (เฉพาะเมื่อสำเร็จ)
	Error   *ErrorInfo  `json:"error,omitempty"` // 🚨 ข้อมูล error (เฉพาะเมื่อผิดพลาด)
}

// 🚨 ErrorInfo - โครงสร้างข้อมูล error ที่ส่งกลับ
type ErrorInfo struct {
	Code    string `json:"code"`              // 🏷️ รหัส error (เช่น "VALIDATION_ERROR")
	Message string `json:"message"`           // 💬 ข้อความ error ที่อ่านง่าย
	Details string `json:"details,omitempty"` // 📝 รายละเอียดเพิ่มเติม (ถ้ามี)
}

// ✅ SuccessResponse - ส่ง response เมื่อสำเร็จ
// ใช้สำหรับส่งข้อมูลกลับเมื่อ API ทำงานสำเร็จ
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, StandardResponse{
		Success: true, // ✅ บอกว่าสำเร็จ
		Data:    data, // 📦 ข้อมูลที่ต้องการส่งกลับ
	})
}

// 🚨 ErrorResponse - ส่ง response เมื่อเกิดข้อผิดพลาด (ฟังก์ชันหลัก)
// เป็นฟังก์ชันพื้นฐานที่ฟังก์ชัน error อื่นๆ จะเรียกใช้
func ErrorResponse(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, StandardResponse{
		Success: false, // ❌ บอกว่าผิดพลาด
		Error: &ErrorInfo{
			Code:    code,    // รหัส error
			Message: message, // ข้อความหลัก
			Details: details, // รายละเอียดเพิ่มเติม
		},
	})
}

// 🔍 ValidationErrorResponse - ส่ง error เมื่อข้อมูลที่ส่งมาไม่ถูกต้อง
// ใช้เมื่อ: user ส่งข้อมูลผิดรูปแบบ เช่น ID ไม่ใช่ตัวเลข, JSON ผิดรูปแบบ
// HTTP Status: 400 Bad Request
func ValidationErrorResponse(c *gin.Context, details string) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request data", details)
}

// 🔍 NotFoundResponse - ส่ง error เมื่อไม่พบข้อมูลที่ต้องการ
// ใช้เมื่อ: หาข้อมูลด้วย ID แล้วไม่เจอในฐานข้อมูล
// HTTP Status: 404 Not Found
func NotFoundResponse(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", resource+" not found", "")
}

// 🔧 InternalServerErrorResponse - ส่ง error เมื่อเกิดปัญหาใน server
// ใช้เมื่อ: database ล่ม, โค้ดมีบัค, หรือ error ที่ไม่คาดหวัง
// HTTP Status: 500 Internal Server Error
func InternalServerErrorResponse(c *gin.Context, details string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", details)
}

// 🔐 UnauthorizedResponse - ส่ง error เมื่อไม่ได้รับอนุญาต
// ใช้เมื่อ: ไม่มี token, token หมดอายุ, หรือ admin secret ผิด
// HTTP Status: 401 Unauthorized
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message, "")
}

// ⚡ ConflictResponse - ส่ง error เมื่อข้อมูลขัดแย้ง
// ใช้เมื่อ: พยายามสร้างข้อมูลที่มีอยู่แล้ว เช่น email ซ้ำ, SKU ซ้ำ
// HTTP Status: 409 Conflict
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message, "")
}

// 📋 UnprocessableEntityResponse - ส่ง error เมื่อข้อมูลถูกต้องแต่ไม่สามารถประมวลผลได้
// ใช้เมื่อ: ข้อมูลถูกรูปแบบแต่ไม่ผ่าน business logic เช่น สต็อกไม่พอ, ราคาติดลบ
// HTTP Status: 422 Unprocessable Entity
func UnprocessableEntityResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, "")
}

// 💡 ตัวอย่างการใช้งาน:

// ✅ ส่ง response สำเร็จ
// SuccessResponse(c, 200, map[string]string{"message": "Hello World"})
// ผลลัพธ์: {"success": true, "data": {"message": "Hello World"}}

// 🚨 ส่ง validation error
// ValidationErrorResponse(c, "Product ID must be a number")
// ผลลัพธ์: {
//   "success": false,
//   "error": {
//     "code": "VALIDATION_ERROR",
//     "message": "Invalid request data",
//     "details": "Product ID must be a number"
//   }
// }

// 🔍 ส่ง not found error
// NotFoundResponse(c, "Product")
// ผลลัพธ์: {
//   "success": false,
//   "error": {
//     "code": "NOT_FOUND",
//     "message": "Product not found"
//   }
// }
