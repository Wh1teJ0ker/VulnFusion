package api

import (
	"VulnFusion/internal/scanner"
	"net/http"
	"strconv"

	"VulnFusion/internal/auth"
	"VulnFusion/internal/log"
	"VulnFusion/internal/models"
	"github.com/gin-gonic/gin"
)

// HandleCreateTask 创建扫描任务
// @Summary 创建扫描任务
// @Description 创建新的扫描任务并异步执行，需提供目标和模板
// @Tags Task
// @Accept json
// @Produce json
// @Param data body api.CreateTaskRequest true "任务创建参数"
// @Success 200 {object} map[string]interface{} "任务创建成功，返回任务 ID"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "任务创建失败"
// @Security ApiKeyAuth
// @Router /api/v1/tasks [post]
func HandleCreateTask(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	var req struct {
		Target   string `json:"target"`
		Template string `json:"template"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("创建任务参数解析失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	task := &models.Task{
		UserID:   claims.UserID,
		Target:   req.Target,
		Template: req.Template,
		Status:   "pending", // 初始状态
	}

	if err := models.CreateTask(task); err != nil {
		log.Error("任务创建失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "任务创建失败"})
		return
	}

	// 异步执行扫描
	go func(taskID uint, target, template string) {
		_ = models.UpdateTaskStatus(taskID, "running")

		output, err := scanner.RunScanTask(scanner.ScanOptions{
			Target:     target,
			Template:   template,
			JsonOutput: true,
			Silent:     true,
		})
		if err != nil {
			log.Error("任务 %d 扫描失败: %v", taskID, err)
			_ = models.UpdateTaskStatus(taskID, "failed")
			return
		}

		parsed, err := scanner.ParseNucleiResult(output)
		if err != nil {
			log.Warn("任务 %d 扫描结果解析失败: %v", taskID, err)
			_ = models.UpdateTaskStatus(taskID, "done")
			return
		}

		for _, p := range parsed {
			res := &models.Result{
				TaskID:        taskID,
				Target:        p.Matched,
				Vulnerability: p.Info.Name,
				Severity:      p.Info.Severity,
				Detail:        string(output), // 可选：改为 JSON 单行或详情部分
			}
			_ = models.SaveScanResult(res)
		}

		_ = models.UpdateTaskStatus(taskID, "done")
	}(task.ID, req.Target, req.Template)

	ctx.JSON(http.StatusOK, gin.H{"message": "任务创建成功", "task_id": task.ID})
}

// HandleGetTaskByID 获取任务详情
// @Summary 获取任务详情
// @Description 根据任务 ID 获取对应任务内容，需权限校验
// @Tags Task
// @Produce json
// @Param id path int true "任务 ID"
// @Success 200 {object} models.Task "任务信息"
// @Failure 400 {object} map[string]string "ID 错误"
// @Failure 403 {object} map[string]string "无权限访问"
// @Failure 404 {object} map[string]string "任务不存在"
// @Security ApiKeyAuth
// @Router /api/v1/tasks/{id} [get]
func HandleGetTaskByID(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID 格式错误"})
		return
	}

	task, err := models.GetTaskByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if claims.Role != "admin" && task.UserID != claims.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权限访问此任务"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// HandleListMyTasks 获取当前用户的任务列表
// @Summary 获取我的任务
// @Description 返回当前用户创建的所有扫描任务
// @Tags Task
// @Produce json
// @Success 200 {array} models.Task "任务列表"
// @Failure 500 {object} map[string]string "获取失败"
// @Security ApiKeyAuth
// @Router /api/v1/tasks [get]
func HandleListMyTasks(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	tasks, err := models.ListTasksByUserID(claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

// HandleListAllTasks 管理员获取全系统任务列表
// @Summary 获取所有任务（管理员）
// @Description 管理员可查看所有用户的任务信息
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Task "所有任务列表"
// @Failure 403 {object} map[string]string "权限不足"
// @Failure 500 {object} map[string]string "获取失败"
// @Security ApiKeyAuth
// @Router /api/v1/admin/tasks [get]
func HandleListAllTasks(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	if claims.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}
	tasks, err := models.ListAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		return
	}
	ctx.JSON(http.StatusOK, tasks)
}

// HandleDeleteTaskByID 删除任务
// @Summary 删除任务
// @Description 根据任务 ID 删除，需本人或管理员权限
// @Tags Task
// @Produce json
// @Param id path int true "任务 ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 403 {object} map[string]string "权限不足"
// @Failure 404 {object} map[string]string "任务不存在"
// @Failure 500 {object} map[string]string "删除失败"
// @Security ApiKeyAuth
// @Router /api/v1/tasks/{id} [delete]
func HandleDeleteTaskByID(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID 格式错误"})
		return
	}

	task, err := models.GetTaskByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if claims.Role != "admin" && task.UserID != claims.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此任务"})
		return
	}

	if err := models.DeleteTaskByID(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "任务已删除"})
}

// HandleBatchDeleteTasks 批量删除任务
// @Summary 批量删除任务
// @Description 批量删除指定任务 ID 列表，需管理员或本人权限
// @Tags Task
// @Accept json
// @Produce json
// @Param data body api.BatchDeleteRequest true "要删除的任务 ID 列表"
// @Success 200 {object} map[string]string "批量删除成功"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "删除失败"
// @Security ApiKeyAuth
// @Router /api/v1/tasks/batch_delete [post]
func HandleBatchDeleteTasks(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*auth.CustomClaims)
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var userID *uint
	if claims.Role != "admin" {
		userID = &claims.UserID
	}

	if err := models.BatchDeleteTasks(req.IDs, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "批量删除成功"})
}

// HandleUpdateTaskStatus 更新任务状态
// @Summary 更新任务状态
// @Description 手动更新任务状态（如 pending, running, done 等）
// @Tags Task
// @Accept json
// @Produce json
// @Param data body api.UpdateStatusRequest true "状态更新参数"
// @Success 200 {object} map[string]string "更新成功"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "更新失败"
// @Security ApiKeyAuth
// @Router /api/v1/tasks/status [post]
func HandleUpdateTaskStatus(ctx *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := models.UpdateTaskStatus(req.ID, req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "任务状态已更新"})
}
