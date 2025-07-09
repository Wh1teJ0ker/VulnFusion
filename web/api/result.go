package api

import (
	"net/http"
	"strconv"

	"VulnFusion/internal/auth"
	"VulnFusion/internal/models"
	"github.com/gin-gonic/gin"
)

// HandleListResultsByTask 获取指定任务的扫描结果
// @Summary 获取任务扫描结果
// @Description 获取某个任务 ID 下所有扫描结果（用户或管理员）
// @Tags Result
// @Produce json
// @Param task_id path int true "任务 ID"
// @Success 200 {array} models.Result "扫描结果列表"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 403 {object} map[string]string "无权限访问"
// @Failure 500 {object} map[string]string "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/results/task/{task_id} [get]
func HandleListResultsByTask(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	taskID, err := strconv.Atoi(ctx.Param("task_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "任务 ID 错误"})
		return
	}

	task, err := models.GetTaskByID(uint(taskID))
	if err != nil || (task.UserID != claims.UserID && claims.Role != "admin") {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权限访问此任务结果"})
		return
	}

	results, err := models.ListResultsByTaskID(uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取扫描结果失败"})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

// HandleListAllResults 获取全系统扫描结果（管理员）
// @Summary 获取所有扫描结果
// @Description 管理员查看所有任务的扫描结果
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Result "所有扫描结果"
// @Failure 403 {object} map[string]string "权限不足"
// @Failure 500 {object} map[string]string "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/results [get]
func HandleListAllResults(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	if claims.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	results, err := models.ListAllResults()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取结果失败"})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

// HandleGetResultDetail 获取扫描结果详情
// @Summary 查看单个扫描结果
// @Description 根据结果 ID 查询具体的漏洞信息
// @Tags Result
// @Produce json
// @Param id path int true "扫描结果 ID"
// @Success 200 {object} models.Result "扫描结果详情"
// @Failure 400 {object} map[string]string "参数格式错误"
// @Failure 404 {object} map[string]string "结果不存在"
// @Router /api/v1/results/{id} [get]
func HandleGetResultDetail(ctx *gin.Context) {
	resultID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID 格式错误"})
		return
	}

	result, err := models.GetResultByID(uint(resultID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "结果不存在"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// HandleDeleteResultsByTask 删除任务的所有扫描结果
// @Summary 删除任务结果
// @Description 删除指定任务下所有扫描结果（用户本人或管理员）
// @Tags Result
// @Produce json
// @Param task_id path int true "任务 ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 403 {object} map[string]string "无权限"
// @Failure 500 {object} map[string]string "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/results/task/{task_id} [delete]
func HandleDeleteResultsByTask(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	taskID, err := strconv.Atoi(ctx.Param("task_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "任务 ID 错误"})
		return
	}

	task, err := models.GetTaskByID(uint(taskID))
	if err != nil || (task.UserID != claims.UserID && claims.Role != "admin") {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此任务结果"})
		return
	}

	if err := models.DeleteResultsByTaskID(uint(taskID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// HandleExportResults 导出任务结果
// @Summary 导出任务结果
// @Description 导出某个任务的扫描结果（JSON 格式）
// @Tags Result
// @Produce json
// @Param task_id path int true "任务 ID"
// @Success 200 {object} map[string]interface{} "任务及其结果数据"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "导出失败"
// @Security ApiKeyAuth
// @Router /api/v1/results/export/{task_id} [get]
func HandleExportResults(ctx *gin.Context) {
	taskID, err := strconv.Atoi(ctx.Param("task_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "任务 ID 错误"})
		return
	}

	results, err := models.ListResultsByTaskID(uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "导出失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"results": results,
	})
}
