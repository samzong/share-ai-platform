package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(viper.GetString("server.mode"))

	// 创建 Gin 引擎
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = viper.GetStringSlice("server.cors.allowed_origins")
	config.AllowMethods = viper.GetStringSlice("server.cors.allowed_methods")
	config.AllowHeaders = viper.GetStringSlice("server.cors.allowed_headers")
	config.AllowCredentials = viper.GetBool("server.cors.allow_credentials")
	r.Use(cors.New(config))

	// 基本路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 启动服务器
	port := viper.GetString("server.port")
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func initConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// 添加配置文件的搜索路径
	viper.AddConfigPath("./backend/configs")     // 本地开发路径
	viper.AddConfigPath("./configs")             // Docker 路径
	viper.AddConfigPath(filepath.Join("..", "configs")) // 相对于 cmd 目录的路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
} 