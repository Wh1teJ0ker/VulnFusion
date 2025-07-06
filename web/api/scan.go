package api

//
//import (
//	"VulnFusion/internal/scanner"
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//type ScanRequest struct {
//	Targets     []string `json:"targets"`
//	Template    string   `json:"template.go"`
//	Concurrency int      `json:"concurrency"`
//}
//
//func HandleScanTask(c *gin.Context) {
//	var req ScanRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
//		return
//	}
//
//	err := scanner.Dispatch(scanner.Task{
//		Name:        req.Template,
//		Targets:     req.Targets,
//		Concurrency: req.Concurrency,
//	})
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "扫描任务已启动"})
//}
//
//func HandleStatus(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{"status": "服务运行中"})
//}
