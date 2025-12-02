package main

import (
	"ProtectedArea/internal/handler"
	"ProtectedArea/internal/router"
	"ProtectedArea/internal/service"
	"ProtectedArea/internal/store"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. 初始化数据库连接
	dsn := "root:123456@tcp(127.0.0.1:3306)/protected_area?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 2. 依赖注入 (层层组装)
	// Store 依赖 DB
	natureStore := store.NewNatureStore(db)
	// Service 依赖 Store
	natureService := service.NewNatureService(natureStore)
	// Handler 依赖 Service
	natureHandler := handler.NewNatureHandler(natureService)

	// 3. 初始化路由
	r := router.InitRouter(natureHandler)

	// 4. 启动服务
	//log.Println("服务启动在 :8080 端口...")
	r.Run(":9094")
}
