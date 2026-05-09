package main

import (
	"fmt"
	"go-gin-demo-service/database"
	"go-gin-demo-service/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	// 用户路由
	users := r.Group("/users")
	{
		users.GET("", handler.GetUsers)
		users.GET("/:id", handler.GetUser)
		users.POST("", handler.CreateUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}

	// 登录
	r.POST("/login", handler.Login)

	// 资源下载
	r.GET("/files", handler.ListFiles)
	r.GET("/download/:filename", handler.Download)

	// App更新检查
	r.POST("/app/update", handler.CheckUpdate)

	// App管理
	appAPI := r.Group("/api/app")
	{
		appAPI.GET("", handler.GetAppList)
		appAPI.POST("/upload", handler.UploadApp)
		appAPI.DELETE("/:id", handler.DeleteApp)
	}

	// 模块管理
	r.GET("/admin", handler.AdminPage)
	api := r.Group("/api/modules")
	{
		api.GET("", handler.GetModules)
		api.POST("/upload", handler.UploadModule)
		api.PUT("/:id", handler.UpdateModule)
		api.DELETE("/:id", handler.DeleteModule)
		api.GET("/:id/history", handler.GetModuleHistory)
		api.GET("/:id/download", handler.DownloadModule)
	}

	// 分包模块
	r.GET("/api/chunk-modules", handler.GetChunkModuleList)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	fmt.Println("服务已启动: http://localhost:8080")
	r.Run("0.0.0.0:8080")
}
