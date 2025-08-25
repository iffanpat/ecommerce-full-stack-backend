# 🎓 คู่มือเริ่มต้น MVC สำหรับผู้เรียน Go

## 🎯 เป้าหมาย
เข้าใจ MVC architecture และสามารถเขียน Go API ได้ด้วยตัวเอง

## 📚 สิ่งที่ต้องเรียนรู้ก่อน
- Go basics (variables, functions, structs)
- HTTP basics (GET, POST, PUT, DELETE)
- JSON format
- SQL basics

## 🚀 เริ่มต้นเรียนรู้ (5 ขั้นตอน)

### ขั้นตอนที่ 1: เข้าใจ Struct และ JSON 🎯

```go
// Model - โครงสร้างข้อมูล
type Product struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// ใน Go จะแปลงเป็น JSON อัตโนมัติ:
// {"id": 1, "name": "iPhone"}
```

**การบ้าน:** ลองสร้าง struct `User` ที่มี ID, Name, Email

### ขั้นตอนที่ 2: เข้าใจ Repository Pattern 🗄️

```go
// Interface - สัญญาว่าต้องมี method อะไรบ้าง
type ProductRepository interface {
    GetAll() ([]Product, error)
}

// Implementation - ตัวจริงที่ทำงาน
type productRepository struct {
    db *sql.DB
}

func (r *productRepository) GetAll() ([]Product, error) {
    // ดึงข้อมูลจากฐานข้อมูล
}
```

**การบ้าน:** ลองเพิ่ม method `GetByID(id int) (*Product, error)`

### ขั้นตอนที่ 3: เข้าใจ Controller Pattern 🎮

```go
func (ctrl *ProductController) ListProducts(c *gin.Context) {
    // 1. เรียก Repository
    products, err := ctrl.repo.GetAll()
    
    // 2. ตรวจสอบ Error
    if err != nil {
        c.JSON(500, gin.H{"error": "เกิดข้อผิดพลาด"})
        return
    }
    
    // 3. ส่ง Response
    c.JSON(200, products)
}
```

**การบ้าน:** ลองเพิ่ม validation ตรวจสอบว่า products ไม่เป็น empty

### ขั้นตอนที่ 4: เข้าใจ View Pattern 🎨

```go
func FormatProductList(products []Product) map[string]interface{} {
    return map[string]interface{}{
        "items": products,
        "total": len(products),
        "message": "ดึงข้อมูลสำเร็จ",
    }
}
```

**การบ้าน:** ลองเพิ่ม timestamp เข้าไปใน response

### ขั้นตอนที่ 5: เข้าใจ Routes Pattern 🛣️

```go
func SetupRoutes(r *gin.Engine, ctrl *ProductController) {
    api := r.Group("/api/v1")
    {
        api.GET("/products", ctrl.ListProducts)
        api.GET("/products/:id", ctrl.GetProduct)
    }
}
```

**การบ้าน:** ลองเพิ่ม route สำหรับ POST /api/v1/products

## 🔧 แบบฝึกหัด: สร้าง User API

### ไฟล์ที่ต้องสร้าง:

1. **models/user.go**
```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

2. **repositories/user_repository.go**
```go
type UserRepository interface {
    GetAll() ([]models.User, error)
    GetByID(id int) (*models.User, error)
}
```

3. **controllers/user_controller.go**
```go
func (ctrl *UserController) ListUsers(c *gin.Context) {
    // TODO: ดึงข้อมูล users ทั้งหมด
}
```

4. **views/user_view.go**
```go
func FormatUserList(users []models.User) UserListResponse {
    // TODO: จัดรูปแบบ response
}
```

5. **routes/routes.go** (เพิ่ม user routes)
```go
users := v1.Group("/users")
{
    users.GET("/", userController.ListUsers)
    users.GET("/:id", userController.GetUser)
}
```

## 📖 ลำดับการศึกษา

### สัปดาห์ที่ 1: พื้นฐาน
- [ ] อ่าน `internal/models/product.go`
- [ ] เข้าใจ struct และ JSON tags
- [ ] ลองแก้ไข struct เพิ่มฟิลด์ใหม่

### สัปดาห์ที่ 2: Repository
- [ ] อ่าน `internal/repositories/product_repository.go`
- [ ] เข้าใจ interface และ implementation
- [ ] เข้าใจ SQL queries

### สัปดาห์ที่ 3: Controller
- [ ] อ่าน `internal/controllers/product_controller.go`
- [ ] เข้าใจการจัดการ HTTP requests
- [ ] เข้าใจ error handling

### สัปดาห์ที่ 4: View & Routes
- [ ] อ่าน `internal/views/` ทั้งหมด
- [ ] อ่าน `internal/routes/routes.go`
- [ ] เข้าใจการจัดรูปแบบ response

### สัปดาห์ที่ 5: ปฏิบัติ
- [ ] ลองสร้าง User API ตามแบบฝึกหัด
- [ ] ทดสอบด้วย Postman
- [ ] เพิ่มฟีเจอร์ validation

## 🧪 การทดสอบ

### ใช้ curl:
```bash
# ดูสินค้าทั้งหมด
curl http://localhost:8080/api/v1/products/

# ดูสินค้าตาม ID
curl http://localhost:8080/api/v1/products/1

# อัปเดตสต็อก
curl -X PUT http://localhost:8080/api/v1/products/1/stock \
  -H "Content-Type: application/json" \
  -H "X-Admin-Secret: your_secret" \
  -d '{"stock": 100}'
```

### ใช้ Postman:
1. เปิด Postman
2. สร้าง Collection ชื่อ "E-commerce API"
3. เพิ่ม requests ตามตัวอย่างด้านบน
4. ลองเรียกดู response

## 🐛 การ Debug

### ดู Log:
```bash
# เริ่ม server และดู logs
go run cmd/server/main.go
```

### เทคนิค Debug:
1. เพิ่ม `log.Printf()` ในโค้ด
2. ใช้ Postman ส่ง request
3. ดู Console logs
4. ตรวจสอบ Database

## 🎯 เป้าหมายสุดท้าย

หลังจากเรียนจบคู่มือนี้ คุณจะสามารถ:
- ✅ เข้าใจ MVC architecture
- ✅ สร้าง API ใหม่ได้ด้วยตัวเอง
- ✅ จัดการ Database operations
- ✅ Handle HTTP requests/responses
- ✅ เขียน Go code ที่มีโครงสร้างดี

## 📋 Checklist ความเข้าใจ

- [ ] เข้าใจแล้วว่า Model คืออะไร
- [ ] เข้าใจแล้วว่า Repository คืออะไร
- [ ] เข้าใจแล้วว่า Controller คืออะไร
- [ ] เข้าใจแล้วว่า View คืออะไร
- [ ] เข้าใจแล้วว่า Routes คืออะไร
- [ ] สามารถอ่านโค้ดเดิมได้
- [ ] สามารถแก้ไขโค้ดเดิมได้
- [ ] สามารถเพิ่มฟีเจอร์ใหม่ได้

## 🆘 ต้องการความช่วยเหลือ?

1. **อ่าน Error messages** ให้ละเอียด
2. **ดู Logs** ใน console
3. **เช็ค Database** ว่ามีข้อมูลหรือไม่
4. **ใช้ Postman** ทดสอบ API
5. **Google search** error message

**Remember: Programming คือการแก้ปัญหา ไม่ใช่การจำ! 🧠** 