package router

import (
	"VulnFusion/web/api"
	"VulnFusion/web/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 初始化所有路由
func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORS())

	apiV1 := r.Group("/api/v1")

	// 公共接口（无需登录）
	apiV1.POST("/auth/login", api.HandleLogin)
	apiV1.POST("/auth/register", api.HandleRegister)
	apiV1.POST("/auth/refresh", api.HandleRefreshToken)

	// 登录后访问接口
	authGroup := apiV1.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		// 用户相关
		authGroup.GET("/user/info", api.HandleGetCurrentUser)
		authGroup.POST("/auth/logout", api.HandleLogout)

		// 扫描任务
		authGroup.POST("/tasks", api.HandleCreateTask)
		authGroup.GET("/tasks", api.HandleListMyTasks)
		authGroup.GET("/tasks/:id", api.HandleGetTaskByID)
		authGroup.DELETE("/tasks/:id", api.HandleDeleteTaskByID)
		authGroup.POST("/tasks/batch_delete", api.HandleBatchDeleteTasks)
		authGroup.POST("/tasks/status", api.HandleUpdateTaskStatus)

		// 扫描结果
		authGroup.GET("/results/task/:task_id", api.HandleListResultsByTask)
		authGroup.DELETE("/results/task/:task_id", api.HandleDeleteResultsByTask)
		authGroup.GET("/results/:id", api.HandleGetResultDetail)
		authGroup.GET("/results/export/:task_id", api.HandleExportResults)
	}

	// 管理员接口（需具备 admin 权限）
	adminGroup := apiV1.Group("/admin")
	adminGroup.Use(middleware.JWTAuthMiddleware(), middleware.RequireAdmin())
	{
		adminGroup.GET("/tasks", api.HandleListAllTasks)
		adminGroup.GET("/results", api.HandleListAllResults)
		adminGroup.GET("/users", api.HandleListAllUsers)
		adminGroup.DELETE("/users/:id", api.HandleDeleteUserByID)
		adminGroup.PUT("/users/:id", api.HandleUpdateUserByID)
		adminGroup.PUT("/users/:id/password", api.HandleResetPasswordByID) // ✅ 新增
	}

}
