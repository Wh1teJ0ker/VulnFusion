package router

import (
	"VulnFusion/web/api"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// 全局中间件：可添加 CORS、日志、异常恢复等
	r.Use(gin.Recovery())

	apiV1 := r.Group("/api/v1")

	// 公共接口（无需登录）
	{
		apiV1.POST("/auth/login", api.HandleLogin)
		apiV1.POST("/auth/register", api.HandleRegister)
	}
	//
	//	// 鉴权接口（需登录）
	//	protected := apiV1.Group("")
	//	protected.Use(middleware.AuthMiddleware())
	//	{
	//		// 用户相关
	//		protected.GET("/auth/me", api.HandleGetCurrentUser)
	//		protected.POST("/auth/logout", api.HandleLogout)
	//
	//		// 扫描管理
	//		protected.POST("/scan", api.HandleScanTask)             // 新建任务
	//		protected.GET("/scan/:id", api.HandleGetScanResult)     // 获取任务详情
	//		protected.GET("/scan", api.HandleListScanTasks)         // 获取任务列表
	//		protected.DELETE("/scan/:id", api.HandleDeleteScanTask) // 删除任务
	//
	//		// 模板管理
	//		protected.GET("/templates", api.HandleListTemplates)
	//		protected.POST("/templates/upload", api.HandleUploadTemplate)
	//		protected.DELETE("/templates/:id", api.HandleDeleteTemplate)
	//
	//		// 系统状态
	//		protected.GET("/status", api.HandleSystemStatus)
	//	}
}
