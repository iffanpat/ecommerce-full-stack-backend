# E-commerce Backend API

Go-based e-commerce backend API ที่ใช้ Gin framework, PostgreSQL และ **MVC Architecture Pattern**

## 🏗️ Architecture

โปรเจคนี้ใช้ **MVC (Model-View-Controller)** pattern เพื่อให้โค้ดมีโครงสร้างที่ชัดเจน:

- **Model** (`internal/models/`) - โครงสร้างข้อมูล
- **View** (`internal/views/`) - จัดรูปแบบ Response
- **Controller** (`internal/controllers/`) - จัดการ HTTP Requests
- **Repository** (`internal/repositories/`) - Data Access Layer
- **Routes** (`internal/routes/`) - API Routing

📖 ดูรายละเอียดเพิ่มเติมใน [MVC_ARCHITECTURE.md](MVC_ARCHITECTURE.md)

## ข้อกำหนดระบบ

- Go 1.22 หรือใหม่กว่า
- PostgreSQL 12 หรือใหม่กว่า

## การติดตั้งและรัน Local

### 1. ติดตั้ง Dependencies

```bash
go mod tidy
```

### 2. ตั้งค่า PostgreSQL

#### วิธีที่ 1: ใช้ Docker (แนะนำ)
```bash
docker run --name ecommerce-postgres -e POSTGRES_DB=ecommerce -e POSTGRES_USER=username -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15
```

#### วิธีที่ 2: ติดตั้ง PostgreSQL แบบ Local
1. ติดตั้ง PostgreSQL
2. สร้างฐานข้อมูล:
```sql
CREATE DATABASE ecommerce;
CREATE USER username WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE ecommerce TO username;
```

### 3. ตั้งค่า Environment Variables

แก้ไขไฟล์ `.env` ให้ตรงกับการตั้งค่าฐานข้อมูลของคุณ:

```env
DATABASE_URL=postgres://username:password@localhost:5432/ecommerce?sslmode=disable
PORT=8080
GIN_MODE=debug
```

### 4. สร้างตารางและข้อมูลตัวอย่าง

```bash
psql -d ecommerce -f sql/001_init.sql
```

หรือใช้ connection string:
```bash
psql "postgres://username:password@localhost:5432/ecommerce" -f sql/001_init.sql
```

### 5. รันแอปพลิเคชัน

#### วิธีที่ 1: ใช้ go run
```bash
go run cmd/server/main.go
```

#### วิธีที่ 2: Build และรัน
```bash
go build -o server cmd/server/main.go
./server
```

#### วิธีที่ 3: ใช้ Air สำหรับ Hot Reload (แนะนำสำหรับ Development)
```bash
# ติดตั้ง Air
go install github.com/cosmtrek/air@latest

# รัน
air
```

### 6. ทดสอบ API

API จะรันที่ `http://localhost:8080`

#### Health Check
```bash
curl http://localhost:8080/health
```

#### ดูสินค้าทั้งหมด
```bash
curl http://localhost:8080/products
```

## API Endpoints

### 🔗 Legacy API (Backward Compatible)
- `GET /products` - ดูสินค้าทั้งหมด
- `GET /products/:id` - ดูสินค้าตาม ID
- `PUT /products/:id/stock` - อัปเดตสต็อกสินค้า
- `POST /carts` - สร้างตะกร้าสินค้า
- `GET /carts/:cid/items` - ดูสินค้าในตะกร้า
- `POST /carts/:cid/items` - เพิ่มสินค้าในตะกร้า
- `PATCH /carts/:cid/items/:iid` - อัปเดตจำนวนสินค้าในตะกร้า
- `DELETE /carts/:cid/items/:iid` - ลบสินค้าจากตะกร้า
- `POST /checkout` - สั่งซื้อสินค้า
- `GET /orders` - ดูรายการสั่งซื้อ

### ✨ New API v1 (Recommended)
- `GET /api/v1/products/` - ดูสินค้าทั้งหมด
- `GET /api/v1/products/:id` - ดูสินค้าตาม ID
- `PUT /api/v1/products/:id/stock` - อัปเดตสต็อกสินค้า
- `POST /api/v1/carts/` - สร้างตะกร้าสินค้า
- `GET /api/v1/carts/:cid/items` - ดูสินค้าในตะกร้า
- `POST /api/v1/carts/:cid/items` - เพิ่มสินค้าในตะกร้า
- `PATCH /api/v1/carts/:cid/items/:iid` - อัปเดตจำนวนสินค้าในตะกร้า
- `DELETE /api/v1/carts/:cid/items/:iid` - ลบสินค้าจากตะกร้า
- `POST /api/v1/orders/checkout` - สั่งซื้อสินค้า
- `GET /api/v1/orders/` - ดูรายการสั่งซื้อ

### 🏥 Health Check
- `GET /health` - ตรวจสอบสถานะเซิร์ฟเวอร์

## การแก้ไขปัญหาที่พบบ่อย

### 1. Database Connection Error
- ตรวจสอบว่า PostgreSQL รันอยู่
- ตรวจสอบ `DATABASE_URL` ในไฟล์ `.env`
- ตรวจสอบ username/password

### 2. Port Already in Use
- เปลี่ยน `PORT` ในไฟล์ `.env`
- หรือหยุดโปรเซสที่ใช้ port 8080

### 3. Module Not Found
- รัน `go mod tidy` เพื่อดาวน์โหลด dependencies 