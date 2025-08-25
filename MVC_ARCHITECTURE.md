# MVCS Architecture Guide สำหรับผู้เริ่มต้น Go 🚀

## 🎯 MVCS คืออะไร?

**MVCS (Model-View-Controller-Service)** เป็นแบบแผนการออกแบบโปรแกรม ที่แยกโค้ดออกเป็น 4 ส่วนหลัก:

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Model     │    │   View      │    │ Controller  │    │   Service   │
│ (ข้อมูล)    │    │ (การแสดงผล) │    │ (ควบคุม)    │    │ (ธุรกิจ)     │
│             │    │             │    │             │    │             │
│ - โครงสร้าง │    │ - JSON      │    │ - รับ HTTP  │    │ - Logic     │
│ - ข้อมูล    │    │ - Response  │    │ - ประมวลผล │    │ - Validation│
│ - Database  │    │ - Format    │    │ - เรียก    │    │ - Security  │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

## 🔄 การทำงานของ MVCS (Data Flow)

```
1. HTTP Request (เช่น GET /products)
       ↓
2. Router ส่งไปหา Controller
       ↓
3. Controller ประมวลผลพื้นฐาน และเรียก Service
       ↓
4. Service จัดการ Business Logic และเรียก Repository
       ↓
5. Repository ไปดึงข้อมูลจาก Database
       ↓
6. Repository ส่งข้อมูลกลับมาให้ Service
       ↓
7. Service ประมวลผล Business Logic และส่งกลับให้ Controller
       ↓
8. Controller ใช้ View จัดรูปแบบข้อมูล
       ↓
9. View ส่ง JSON Response กลับไปหา Client
```

## 📁 โครงสร้างไฟล์ (อธิบายง่ายๆ)

```
backend/
├── cmd/server/main.go          # 🚪 ประตูหลักของโปรแกรม (เริ่มต้นที่นี่)
│
├── internal/                   # 📦 โค้ดหลักของเรา
│   ├── models/                 # 🎯 Model: โครงสร้างข้อมูล
│   │   └── product.go          #     ► struct ของสินค้า
│   │
│   ├── repositories/           # 🗄️ Repository: จัดการฐานข้อมูล
│   │   └── product_repository.go #    ► ดึง/เซฟข้อมูลสินค้า
│   │
│   ├── services/               # 🏢 Service: จัดการ Business Logic
│   │   └── product_service.go  #     ► ตรรกะทางธุรกิจ, validation, security
│   │
│   ├── controllers/            # 🎮 Controller: ควบคุมการทำงาน
│   │   └── product_controller.go #    ► จัดการ request เกี่ยวกับสินค้า
│   │
│   ├── views/                  # 🎨 View: จัดรูปแบบ response
│   │   ├── response.go         #     ► format response มาตรฐาน
│   │   └── product_view.go     #     ► format response สินค้า
│   │
│   └── routes/                 # 🛣️ Routes: กำหนดเส้นทาง API
│       └── routes.go           #     ► กำหนด URL endpoints
│
└── sql/                        # 🗃️ Database Schema
    └── 001_init.sql            #     ► โครงสร้างตาราง
```

## 🧱 อธิบายแต่ละชั้น (Layer)

### 1. 🎯 Model Layer - ข้อมูล

**หน้าที่:** กำหนดโครงสร้างข้อมูล (struct)

```go
// internal/models/product.go
type Product struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Price       int    `json:"price_cents"`
    Stock       int    `json:"stock"`
}
```

**เข้าใจง่ายๆ:** เปรียบเสมือน "แบบฟอร์ม" ที่กำหนดว่าข้อมูลสินค้าต้องมีอะไรบ้าง

### 2. 🗄️ Repository Layer - ฐานข้อมูล

**หน้าที่:** จัดการการดึง/เซฟข้อมูลจากฐานข้อมูล

```go
// internal/repositories/product_repository.go
type ProductRepository interface {
    GetAll() ([]Product, error)        // ดึงสินค้าทั้งหมด
    GetByID(id int) (*Product, error)  // ดึงสินค้าตาม ID
    UpdateStock(id, stock int) error   // อัปเดตสต็อก
}
```

**เข้าใจง่ายๆ:** เปรียบเสมือน "พนักงานคลัง" ที่รู้วิธีเข้าไปหยิบ/วางข้อมูลในฐานข้อมูล

### 3. 🏢 Service Layer - ตรรกะทางธุรกิจ

**หน้าที่:** จัดการ Business Logic, Validation, Security และ Authentication

```go
// internal/services/product_service.go
type ProductService interface {
    GetAllProducts() ([]Product, error)                         // ดึงสินค้าทั้งหมด
    GetProductByID(id int) (*Product, error)                   // ดึงสินค้าตาม ID
    UpdateProductStock(id, stock int, adminSecret string) error // อัปเดตสต็อก
}
```

**เข้าใจง่ายๆ:** เปรียบเสมือน "ผู้เชี่ยวชาญ" ที่ตรวจสอบกฎระเบียบธุรกิจ, ความปลอดภัย และความถูกต้องของข้อมูล

**Service ทำอะไรบ้าง:**
- ✅ ตรวจสอบสิทธิ์ Admin (Authentication)
- ✅ ตรวจสอบความถูกต้องของข้อมูล (Validation)
- ✅ จัดการ Business Rules (เช่น สต็อกห้ามติดลบ)
- ✅ แปลง Error จาก Repository ให้เป็น Business Error
- ✅ Log การทำงานเพื่อ Audit

### 4. 🎮 Controller Layer - ควบคุม

**หน้าที่:** รับ HTTP request และประสานงาน (ไม่มี Business Logic)

```go
// internal/controllers/product_controller.go
func (ctrl *ProductController) ListProducts(c *gin.Context) {
    // 1. เรียก Service ให้ดึงข้อมูล
    products, err := ctrl.service.GetAllProducts()
    if err != nil {
        // 2. ถ้าเกิดข้อผิดพลาด ส่ง error response
        views.InternalServerErrorResponse(c, "ดึงข้อมูลไม่ได้")
        return
    }
    
    // 3. ใช้ View จัดรูปแบบข้อมูล
    response := views.FormatProductList(products)
    
    // 4. ส่ง response กลับไป
    views.SuccessResponse(c, 200, response)
}
```

**เข้าใจง่ายๆ:** เปรียบเสมือน "พนักงานต้อนรับ" ที่รับคำสั่งจากลูกค้า แล้วไปติดต่อผู้เชี่ยวชาญ

### 5. 🎨 View Layer - การแสดงผล

**หน้าที่:** จัดรูปแบบ JSON response

```go
// internal/views/product_view.go
func FormatProductList(products []Product) ProductListResponse {
    return ProductListResponse{
        Items: products,
        Total: len(products),
    }
}
```

**เข้าใจง่ายๆ:** เปรียบเสมือน "พนักงานแพ็คของ" ที่จัดข้อมูลให้อยู่ในรูปแบบที่สวยงาม

### 6. 🛣️ Routes Layer - เส้นทาง

**หน้าที่:** กำหนดว่า URL ไหนไปหา Controller ไหน

```go
// internal/routes/routes.go
func SetupRoutes(r *gin.Engine, productController *ProductController) {
    v1 := r.Group("/api/v1")
    {
        products := v1.Group("/products")
        {
            products.GET("/", productController.ListProducts)      // GET /api/v1/products/
            products.GET("/:id", productController.GetProduct)     // GET /api/v1/products/123
            products.PUT("/:id/stock", productController.UpdateStock) // PUT /api/v1/products/123/stock
        }
    }
}
```

**เข้าใจง่ายๆ:** เปรียบเสมือน "ป้ายบอกทาง" ที่บอกว่าถ้าลูกค้ามาทางไหน ให้ไปหาพนักงานต้อนรับคนไหน

## 🚀 ตัวอย่างการทำงานจริง

### สถานการณ์: ลูกค้าต้องการดูสินค้าทั้งหมด

1. **Client ส่ง Request:** `GET /api/v1/products/`

2. **Router ตัดสินใจ:** ดู routes.go → ส่งไปหา `productController.ListProducts`

3. **Controller ทำงาน:**
   ```go
   func (ctrl *ProductController) ListProducts(c *gin.Context) {
       // เรียก Service (ไม่ใช่ Repository โดยตรง)
       products, err := ctrl.service.GetAllProducts()
       
       // จัดรูปแบบด้วย View
       response := views.FormatProductList(products)
       
       // ส่ง response
       views.SuccessResponse(c, 200, response)
   }
   ```

4. **Service ทำงาน:**
   ```go
   func (s *productService) GetAllProducts() ([]Product, error) {
       log.Println("🔍 Service: Starting to retrieve all products")
       
       // เรียก Repository
       products, err := s.repo.GetAll()
       if err != nil {
           return nil, fmt.Errorf("failed to retrieve products: %w", err)
       }
       
       // สามารถเพิ่ม business logic ที่นี่ เช่น:
       // - กรองสินค้าที่หมดสต็อก
       // - เพิ่มข้อมูลราคาพิเศษ
       
       return products, nil
   }
   ```

5. **Repository ทำงาน:**
   ```go
   func (r *productRepository) GetAll() ([]Product, error) {
       // Query ฐานข้อมูล
       query := "SELECT id, name, price_cents, stock FROM products"
       rows, err := r.db.Query(query)
       // ... ประมวลผลข้อมูล
       return products, nil
   }
   ```

6. **View จัดรูปแบบ:**
   ```go
   func FormatProductList(products []Product) ProductListResponse {
       return ProductListResponse{
           Items: products,
           Total: len(products),
       }
   }
   ```

7. **Response ที่ได้:**
   ```json
   {
     "success": true,
     "data": {
       "items": [
         {"id": 1, "name": "iPhone", "price_cents": 2500000, "stock": 10},
         {"id": 2, "name": "Samsung", "price_cents": 1800000, "stock": 5}
       ],
       "total": 2
     }
   }
   ```

## 🔄 การทำงานของ UpdateStock ที่ซับซ้อนขึ้น

### สถานการณ์: Admin ต้องการอัปเดตสต็อกสินค้า

1. **Request:** `PUT /api/v1/products/1/stock` + Header: `X-Admin-Secret: secret123`

2. **Controller:** รับ request และส่งต่อไป Service
   ```go
   adminSecret := c.GetHeader("X-Admin-Secret")
   err := ctrl.service.UpdateProductStock(id, req.Stock, adminSecret)
   ```

3. **Service ทำงานหนัก:**
   ```go
   func (s *productService) UpdateProductStock(id, stock int, adminSecret string) error {
       // 🔐 ตรวจสอบสิทธิ์ Admin
       if err := s.validateAdminAccess(adminSecret); err != nil {
           return err
       }
       
       // ✅ ตรวจสอบ Business Rules
       if stock < 0 {
           return errors.New("stock cannot be negative")
       }
       
       // 🔍 ตรวจสอบว่าสินค้ามีอยู่จริง
       existingProduct, err := s.repo.GetByID(id)
       if err != nil {
           return fmt.Errorf("product not found")
       }
       
       // 📝 Log การเปลี่ยนแปลง
       log.Printf("Stock change: %s (%d → %d)", existingProduct.Name, existingProduct.Stock, stock)
       
       // เรียก Repository
       return s.repo.UpdateStock(id, stock)
   }
   ```

4. **Repository:** ทำงานกับฐานข้อมูลตามปกติ

## 💡 ทำไมต้องเพิ่ม Service Layer?

### ✅ ข้อดีของ Service Layer

1. **แยก Business Logic** - Controller ไม่ต้องรู้กฎธุรกิจซับซ้อน
2. **ความปลอดภัย** - จัดการ Authentication/Authorization ในที่เดียว
3. **Validation ครบถ้วน** - ตรวจสอบข้อมูลอย่างละเอียด
4. **ทดสอบง่าย** - สามารถเทส Business Logic แยกจาก HTTP
5. **ใช้ซ้ำได้** - Service สามารถเรียกจาก Controller หลายตัว
6. **Log และ Audit** - ติดตามการทำงานได้ละเอียด

### 📋 เปรียบเทียบก่อนและหลังมี Service

**ก่อนมี Service (Controller เรียก Repository โดยตรง):**
```go
func (ctrl *ProductController) UpdateStock(c *gin.Context) {
    // ❌ Business Logic ปะปนใน Controller
    secret := c.GetHeader("X-Admin-Secret")
    if secret != os.Getenv("ADMIN_SECRET") {
        views.UnauthorizedResponse(c, "Invalid secret")
        return
    }
    
    // ❌ Validation ใน Controller
    if req.Stock < 0 {
        views.ValidationErrorResponse(c, "Stock cannot be negative")
        return
    }
    
    // เรียก Repository
    err := ctrl.repo.UpdateStock(id, req.Stock)
}
```

**หลังมี Service (Controller เรียก Service):**
```go
func (ctrl *ProductController) UpdateStock(c *gin.Context) {
    // ✅ Controller เบา เป็นแค่ตัวประสาน
    adminSecret := c.GetHeader("X-Admin-Secret")
    
    // ✅ Service จัดการทุกอย่าง
    err := ctrl.service.UpdateProductStock(id, req.Stock, adminSecret)
    
    // ✅ แค่จัดการ Error Response
    if err != nil {
        // จัดการ error response ตาม error type
    }
}
```

## 🎓 เริ่มต้นเรียนรู้ MVCS

### ขั้นตอนที่ 1: เข้าใจ Model
1. เปิดไฟล์ `internal/models/product.go`
2. ดูโครงสร้าง `Product struct`
3. เข้าใจ JSON tags

### ขั้นตอนที่ 2: เข้าใจ Repository  
1. เปิดไฟล์ `internal/repositories/product_repository.go`
2. ดู interface และ implementation
3. เข้าใจ SQL queries

### ขั้นตอนที่ 3: เข้าใจ Service (ใหม่!)
1. เปิดไฟล์ `internal/services/product_service.go`
2. ดู Business Logic และ Validation
3. เข้าใจการจัดการ Authentication
4. ศึกษา Error Handling

### ขั้นตอนที่ 4: เข้าใจ Controller
1. เปิดไฟล์ `internal/controllers/product_controller.go` 
2. ตามการทำงานของแต่ละ function
3. เข้าใจว่า Controller เป็นแค่ตัวประสาน

### ขั้นตอนที่ 5: เข้าใจ View
1. เปิดไฟล์ `internal/views/response.go`
2. ดูการจัดรูปแบบ response
3. เข้าใจ error handling

### ขั้นตอนที่ 6: เข้าใจ Routes
1. เปิดไฟล์ `internal/routes/routes.go`
2. ดู URL mapping
3. เข้าใจ API versioning

## 🛠️ API Endpoints ที่มีให้ใช้

### Products API
```bash
# ดูสินค้าทั้งหมด
GET /api/v1/products/

# ดูสินค้าตาม ID  
GET /api/v1/products/1

# อัปเดตสต็อกสินค้า (ต้องใช้ admin secret)
PUT /api/v1/products/1/stock
Headers: X-Admin-Secret: your_admin_secret
Body: {"stock": 100}

# Health check
GET /health
```

### Response Format มาตรฐาน

**สำเร็จ:**
```json
{
  "success": true,
  "data": {
    // ข้อมูลที่ต้องการ
  }
}
```

**ผิดพลาด:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "ข้อความที่เข้าใจง่าย",
    "details": "รายละเอียดเพิ่มเติม"
  }
}
```

## 🎯 เทคนิคการเรียนรู้

1. **เริ่มจาก main.go** - ดูการเริ่มต้นโปรแกรม และการสร้าง Service
2. **ตาม flow การทำงาน** - เริ่มจาก router → controller → service → repository → view
3. **ทดลองแก้โค้ด** - เปลี่ยน Business Logic ใน Service
4. **ใช้ Postman ทดสอบ** - ลองเรียก API ดูผลลัพธ์
5. **อ่าน logs** - ดู console เพื่อเข้าใจการทำงานของ Service

## 🔧 การต่อยอด

เมื่อเข้าใจ MVCS แล้ว สามารถเพิ่มฟีเจอร์ใหม่ได้โดย:

1. **เพิ่ม Model** - สร้าง struct ใหม่ใน `models/`
2. **เพิ่ม Repository** - สร้าง interface และ implementation ใหม่
3. **เพิ่ม Service** - สร้าง service ใหม่ที่จัดการ Business Logic
4. **เพิ่ม Controller** - สร้าง controller ใหม่ที่ใช้ service
5. **เพิ่ม View** - สร้าง response formatter
6. **เพิ่ม Routes** - เพิ่ม URL endpoints ใหม่

---

## 📚 เอกสารเพิ่มเติม

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [SQL Tutorial](https://www.w3schools.com/sql/)

**Happy Coding! 🚀** 