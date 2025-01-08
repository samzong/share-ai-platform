package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samzong/share-ai-platform/internal/api/handlers"
	"github.com/samzong/share-ai-platform/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter sets up the router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}  // 前端开发服务器地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置静态文件服务
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			userHandler := handlers.NewUserHandler()
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), userHandler.Logout)
		}

		// 用户相关路由
		users := api.Group("/users", middleware.AuthMiddleware())
		{
			userHandler := handlers.NewUserHandler()
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.PUT("/:id", userHandler.UpdateUser)
			users.PUT("/:id/role", userHandler.UpdateUserRole)
			users.GET("", userHandler.ListUsers)
		}

		// 镜像相关路由
		imageHandler := handlers.NewImageHandler()
		images := api.Group("/images")
		{
			images.GET("", imageHandler.ListImages)
			images.GET("/:id", imageHandler.GetImage)
			
			// 需要认证的路由
			auth := images.Use(middleware.AuthMiddleware())
			{
				auth.POST("/:id/collect", imageHandler.CollectImage)
				auth.DELETE("/:id/collect", imageHandler.UncollectImage)
			}
		}

		// 组织相关路由
		orgs := api.Group("/orgs")
		{
			// 公共镜像路由
			orgs.GET("/public/images", imageHandler.ListImages)
			
			// 需要认证的路由
			auth := orgs.Use(middleware.AuthMiddleware())
			{
				auth.POST("/:org_id/images", imageHandler.CreateImage)
				auth.PUT("/:org_id/images/:id", imageHandler.UpdateImage)
				auth.DELETE("/:org_id/images/:id", imageHandler.DeleteImage)
			}
		}

		// 收藏夹路由
		favorites := api.Group("/favorites").Use(middleware.AuthMiddleware())
		{
			favorites.GET("", imageHandler.ListFavorites)
		}
	}

	return r
} 