package main

import (
	"VulnFusion/internal/bootstrap"
	"VulnFusion/internal/config"
	"VulnFusion/web/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	_ "VulnFusion/docs" // ✅ 加载 swag 生成的文档文件
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title VulnFusion API 文档
// @version 1.0
// @description 用于描述漏洞扫描平台的后端 API。
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 加载配置
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	if err := bootstrap.InitializeSystem(); err != nil {
		panic(err)
	}

	// 启动 Gin 引擎
	r := gin.Default()

	// 注册业务路由
	router.RegisterRoutes(r)

	// 注册 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 挂载静态资源在 /web，避免与 API 路由冲突
	r.StaticFS("/web", http.Dir("./frontend/dist"))

	// fallback 到前端首页（确保前端用的是 HTML5 History 模式路由）
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	// 启动服务
	addr := config.GetListenAddr() // e.g., ":8080"
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
