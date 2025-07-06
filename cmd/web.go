package cmd

import (
	"VulnFusion/web/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartWeb() {
	r := gin.Default()

	// 注册后端 API 路由
	router.Register(r)

	//// 静态资源（暂不启用，后续打包或调试前端时再开启）
	// staticDir := "web/ui/dist"
	// r.Static("/", staticDir)
	// r.NoRoute(func(c *gin.Context) {
	//     c.File(staticDir + "/index.html")
	// })

	fmt.Println("[Web] 启动 Web API: http://localhost:8000")
	if err := r.Run(":8000"); err != nil {
		panic("启动失败: " + err.Error())
	}
}
