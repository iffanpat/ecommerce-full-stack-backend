package main

import (
	"log"
	"os"

	"ecommerce/internal/controllers"
	"ecommerce/internal/db"
	"ecommerce/internal/repositories"
	"ecommerce/internal/routes"
	"ecommerce/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// โหลดไฟล์ .env เพื่อเซต environment variables
	// ถ้าไม่มีไฟล์ .env ก็ไม่เป็นไร จะใช้ system environment แทน
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// เชื่อมต่อฐานข้อมูล
	// ฟังก์ชัน InitDB() จะใช้ connection string จาก environment variable
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err) // ถ้าเชื่อมต่อไม่ได้ ให้หยุดโปรแกรม
	}

	// Repository = ตัวจัดการฐานข้อมูล ทำหน้าที่ CRUD (Create, Read, Update, Delete)
	productRepo := repositories.NewProductRepository(database)

	// Service = ตัวจัดการ business logic, validation, และการประมวลผลข้อมูล
	// Service จะเรียกใช้ Repository เพื่อเข้าถึงข้อมูล
	productService := services.NewProductService(productRepo)

	// Controller = ตัวควบคุม รับ HTTP request และประมวลผล
	// ต้องส่ง service เข้าไปเพื่อให้ controller เรียกใช้ business logic ได้
	productController := controllers.NewProductController(productService)

	// สร้าง Gin router (HTTP server framework)
	r := gin.Default()

	// เซต CORS (Cross-Origin Resource Sharing)
	// เพื่อให้ frontend จากโดเมนอื่นเรียก API ได้
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // อนุญาตทุกโดเมน (ใช้ในการพัฒนา)
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "X-Admin-Secret"},
	}))

	// เซต Routes (API endpoints)
	// ส่ง router และ controller ไปให้ routes.go จัดการ
	routes.SetupRoutes(r, productController)

	// เริ่มต้น HTTP server
	// ใช้ PORT จาก environment variable หรือ 8080 เป็นค่าเริ่มต้น
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// แสดงข้อมูลเมื่อ server เริ่มทำงาน
	log.Printf("Server starting on port %s", port)
	log.Printf("Architecture: Controller → Service → Repository → Database")
	log.Printf("- Health check: http://localhost:%s/health", port)
	log.Printf("- Products (Legacy): http://localhost:%s/products", port)
	log.Printf("- Products (v1): http://localhost:%s/api/v1/products", port)

	// ฟังก์ชันนี้จะทำงานไปเรื่อยๆ จนกว่าจะหยุดโปรแกรม
	r.Run(":" + port)
}
