package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/samzong/share-ai-platform/internal/api"
	"github.com/samzong/share-ai-platform/internal/database"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/samzong/share-ai-platform/docs"  // 导入 swagger docs
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

	// 创建 Gin 引擎
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}  // 前端开发服务器地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// 注册路由
	handler := api.NewHandler()
	handler.RegisterRoutes(r)

	// 添加 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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