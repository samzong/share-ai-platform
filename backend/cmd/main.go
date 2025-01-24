package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/samzong/share-ai-platform/docs" // 导入 swagger docs
	"github.com/samzong/share-ai-platform/internal/api"
	"github.com/samzong/share-ai-platform/internal/database"
	"github.com/spf13/viper"
)

// @title           Share AI Platform API
// @version         1.0
// @description     This is the API server for Share AI Platform.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(viper.GetString("server.mode"))

	// 初始化路由
	r := api.SetupRouter()

	// 打印所有注册的路由
	for _, route := range r.Routes() {
		log.Printf("Route: %s %s\n", route.Method, route.Path)
	}

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
	viper.AddConfigPath("./backend/config")            // 本地开发路径
	viper.AddConfigPath("./config")                    // Docker 路径
	viper.AddConfigPath(filepath.Join("..", "config")) // 相对于 cmd 目录的路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}
