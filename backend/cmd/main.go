package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/samzong/share-ai-platform/internal/api"
	"github.com/samzong/share-ai-platform/internal/database"
)

func init() {
	// 加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../configs")
	
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 设置 Gin 模式
	if viper.GetString("server.mode") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化 Redis 连接
	if err := database.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// 创建 Gin 实例
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = viper.GetStringSlice("cors.allowed_origins")
	config.AllowMethods = viper.GetStringSlice("cors.allowed_methods")
	config.AllowHeaders = viper.GetStringSlice("cors.allowed_headers")
	config.AllowCredentials = viper.GetBool("cors.allow_credentials")
	config.MaxAge = viper.GetDuration("cors.max_age")
	r.Use(cors.New(config))

	// 注册路由
	handler := api.NewHandler()
	handler.RegisterRoutes(r)

	// 启动服务器
	port := viper.GetString("server.port")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
} 