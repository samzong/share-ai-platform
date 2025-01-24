package main

import (
	"log"
	"path/filepath"

	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/samzong/share-ai-platform/internal/models"
	"github.com/spf13/viper"
)

func init() {
	// 设置配置文件路径
	viper.SetConfigName("migrate")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config") // 相对于当前目录的上级configs目录
	viper.AddConfigPath("config")       // 当前目录的configs目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// 打印当前配置文件路径
	configFile := viper.ConfigFileUsed()
	absPath, _ := filepath.Abs(configFile)
	log.Printf("Using config file: %s", absPath)
}

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
		&models.Image{},
		&models.Label{},
		&models.Collection{},
	); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
