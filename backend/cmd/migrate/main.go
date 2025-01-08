package main

import (
	"log"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/models"
)

func main() {
	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// 获取数据库连接
	db := database.GetDB()

	// 自动迁移模型
	if err := db.AutoMigrate(
		&models.User{},
		// 在这里添加其他模型
	); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	log.Println("Database migration completed successfully!")
} 