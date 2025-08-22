# MVC Architecture Documentation

## Overview

โปรเจค E-commerce Backend ได้ถูกปรับโครงสร้างให้ใช้ **MVC (Model-View-Controller)** pattern เพื่อให้โค้ดมีความเป็นระเบียบ แยกหน้าที่ชัดเจน และง่ายต่อการบำรุงรักษา

## Architecture Layers

### 1. Model Layer (`internal/models/`)
**หน้าที่:** จัดการโครงสร้างข้อมูล (Data Structures)

- `product.go` - โครงสร้างข้อมูลสินค้า
- `cart.go` - โครงสร้างข้อมูลตะกร้าสินค้า
- `order.go` - โครงสร้างข้อมูลคำสั่งซื้อ

### 2. View Layer (`internal/views/`)
**หน้าที่:** จัดการรูปแบบการตอบกลับ (Response Formatting)

- `response.go` - Standard API response formats
- `product_view.go` - Product-specific response formatters
- `cart_view.go` - Cart-specific response formatters  
- `order_view.go` - Order-specific response formatters

### 3. Controller Layer (`internal/controllers/`)
**หน้าที่:** จัดการ HTTP requests และเรียกใช้ Repository

- `product_controller.go` - จัดการ requests เกี่ยวกับสินค้า
- `cart_controller.go` - จัดการ requests เกี่ยวกับตะกร้าสินค้า
- `order_controller.go` - จัดการ requests เกี่ยวกับคำสั่งซื้อ

### 4. Repository Layer (`internal/repositories/`)
**หน้าที่:** จัดการการเข้าถึงข้อมูล (Data Access Layer)

- `product_repository.go` - การเข้าถึงข้อมูลสินค้า
- `cart_repository.go` - การเข้าถึงข้อมูลตะกร้าสินค้า
- `order_repository.go` - การเข้าถึงข้อมูลคำสั่งซื้อ

### 5. Routes Layer (`internal/routes/`)
**หน้าที่:** จัดการ API routing

- `routes.go` - กำหนด API endpoints และ routing logic

## Data Flow

```
HTTP Request → Controller → Repository → Database
                    ↓
HTTP Response ← View ← Controller ← Repository
```

1. **Request เข้ามา:** HTTP request เข้ามาที่ Controller
2. **Controller ประมวลผล:** ตรวจสอบข้อมูล และเรียกใช้ Repository
3. **Repository ทำงาน:** ดำเนินการกับฐานข้อมูล
4. **View จัดรูปแบบ:** จัดรูปแบบข้อมูลสำหรับตอบกลับ
5. **Response ส่งกลับ:** ส่ง JSON response กลับไปยัง client

## API Endpoints

### Legacy Endpoints (Backward Compatible)
- `GET /products` - ดูสินค้าทั้งหมด
- `GET /products/:id` - ดูสินค้าตาม ID
- `PUT /products/:id/stock` - อัปเดตสต็อกสินค้า
- `POST /carts` - สร้างตะกร้าสินค้า
- `GET /carts/:cid/items` - ดูสินค้าในตะกร้า
- `POST /carts/:cid/items` - เพิ่มสินค้าในตะกร้า
- `PATCH /carts/:cid/items/:iid` - อัปเดตจำนวนสินค้า
- `DELETE /carts/:cid/items/:iid` - ลบสินค้าจากตะกร้า
- `POST /checkout` - สั่งซื้อสินค้า
- `GET /orders` - ดูรายการสั่งซื้อ

### New API v1 Endpoints
- `GET /api/v1/products/` - ดูสินค้าทั้งหมด
- `GET /api/v1/products/:id` - ดูสินค้าตาม ID
- `PUT /api/v1/products/:id/stock` - อัปเดตสต็อกสินค้า
- `POST /api/v1/carts/` - สร้างตะกร้าสินค้า
- `GET /api/v1/carts/:cid/items` - ดูสินค้าในตะกร้า
- `POST /api/v1/carts/:cid/items` - เพิ่มสินค้าในตะกร้า
- `PATCH /api/v1/carts/:cid/items/:iid` - อัปเดตจำนวนสินค้า
- `DELETE /api/v1/carts/:cid/items/:iid` - ลบสินค้าจากตะกร้า
- `POST /api/v1/orders/checkout` - สั่งซื้อสินค้า
- `GET /api/v1/orders/` - ดูรายการสั่งซื้อ

## Response Format

### Success Response
```json
{
  "success": true,
  "data": {
    // response data here
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message",
    "details": "Additional details"
  }
}
```

## Benefits of MVC Architecture

### 1. **Separation of Concerns**
- แต่ละ layer มีหน้าที่ที่ชัดเจน
- ง่ายต่อการแก้ไขและบำรุงรักษา

### 2. **Reusability**
- Repository สามารถใช้ซ้ำได้ในหลาย Controller
- View formatters สามารถใช้ซ้ำได้

### 3. **Testability**
- แต่ละ layer สามารถ test แยกกันได้
- Mock dependencies ได้ง่าย

### 4. **Scalability**
- เพิ่มฟีเจอร์ใหม่ได้ง่าย
- แยก API versioning ได้

### 5. **Maintainability**
- โค้ดมีโครงสร้างที่ชัดเจน
- ง่ายต่อการหา bug และแก้ไข

## Directory Structure

```
backend/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── controllers/            # MVC Controllers
│   │   ├── product_controller.go
│   │   ├── cart_controller.go
│   │   └── order_controller.go
│   ├── views/                  # Response Formatters
│   │   ├── response.go
│   │   ├── product_view.go
│   │   ├── cart_view.go
│   │   └── order_view.go
│   ├── repositories/           # Data Access Layer
│   │   ├── product_repository.go
│   │   ├── cart_repository.go
│   │   └── order_repository.go
│   ├── models/                 # Data Models
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── order.go
│   ├── routes/                 # Route Configuration
│   │   └── routes.go
│   └── db/                     # Database Layer
│       ├── db.go
│       └── *_queries.go (legacy)
└── sql/                        # Database Schema
    └── 001_init.sql
```

## Migration Guide

### จากโครงสร้างเดิม
- `handlers/` → `controllers/` (เปลี่ยนชื่อและปรับโครงสร้าง)
- `services/` → `repositories/` (ย้าย business logic ไป controllers)
- เพิ่ม `views/` สำหรับ response formatting
- เพิ่ม `routes/` สำหรับ route management

### Backward Compatibility
- API endpoints เดิมยังใช้งานได้
- Response format ปรับใหม่แต่ข้อมูลเหมือนเดิม
- Database schema ไม่เปลี่ยนแปลง 