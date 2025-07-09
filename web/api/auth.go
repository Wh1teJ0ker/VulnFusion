package api

import (
	"net/http"
	"strconv"
	"time"

	"VulnFusion/internal/auth"
	"VulnFusion/internal/log"
	"VulnFusion/internal/models"
	"VulnFusion/internal/utils"

	"github.com/gin-gonic/gin"
)

// HandleRegister 用户注册
// @Summary 用户注册
// @Description 使用用户名、密码和角色注册新用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body api.RegisterRequest true "注册参数"
// @Success 200 {object} map[string]string "注册成功"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 500 {object} map[string]string "服务器错误或用户名已存在"
// @Router /api/v1/auth/register [post]
func HandleRegister(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("注册参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error("注册密码加密失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: hashed,
		Role:     req.Role,
	}

	err = models.CreateUser(user)
	if err != nil {
		log.Error("创建用户失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户名已存在或创建失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// HandleLogin 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码进行登录，返回访问令牌与刷新令牌
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body api.LoginRequest true "登录参数"
// @Success 200 {object} map[string]string "包含 access_token 和 refresh_token"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 401 {object} map[string]string "用户名或密码错误"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /api/v1/auth/login [post]
func HandleLogin(ctx *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warn("登录参数绑定失败: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		log.Warn("用户不存在: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := auth.GenerateToken(uint(int(user.ID)), user.Username, user.Role, 15*time.Minute)
	if err != nil {
		log.Error("生成 Token 失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成 token 失败"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(uint(user.ID), 7*24*time.Hour)
	if err != nil {
		log.Error("生成刷新 Token 失败: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成 refresh token 失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// HandleRefreshToken 刷新 JWT Token
// @Summary 刷新 JWT Token
// @Description 使用 Refresh Token 获取新的访问令牌
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body api.RefreshRequest true "刷新令牌参数"
// @Success 200 {object} map[string]string "新的访问令牌"
// @Failure 400 {object} map[string]string "参数错误"
// @Failure 401 {object} map[string]string "无效或过期的刷新令牌"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /api/v1/auth/refresh [post]
func HandleRefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	claims, err := auth.ParseToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "无效或已过期的刷新令牌"})
		return
	}

	userID, _ := strconv.Atoi(claims.Subject)
	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	token, err := auth.GenerateToken(uint(user.ID), user.Username, user.Role, 15*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成 token 失败"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// HandleLogout 用户注销登录
// @Summary 用户注销登录
// @Description 将当前访问令牌加入黑名单，防止后续使用
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string "注销成功"
// @Failure 401 {object} map[string]string "未登录或 token 无效"
// @Failure 500 {object} map[string]string "系统错误"
// @Router /api/v1/auth/logout [post]
func HandleLogout(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或 token 无效"})
		return
	}

	customClaims, ok := claims.(*auth.CustomClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		return
	}

	exp := time.Until(customClaims.ExpiresAt.Time)
	auth.AddTokenToBlacklist(customClaims.ID, exp)

	ctx.JSON(http.StatusOK, gin.H{"message": "注销成功"})
}

// HandleGetCurrentUser 获取当前用户信息
// @Summary 获取当前登录用户信息
// @Description 从 JWT 中解析并返回用户 ID、用户名与角色
// @Tags User
// @Produce json
// @Success 200 {object} map[string]interface{} "包含用户 id、用户名、角色"
// @Failure 401 {object} map[string]string "未登录"
// @Failure 500 {object} map[string]string "系统错误或用户不存在"
// @Security ApiKeyAuth
// @Router /api/v1/user/info [get]
func HandleGetCurrentUser(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	customClaims, ok := claims.(*auth.CustomClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "系统错误"})
		return
	}

	user, err := models.GetUserByID(uint(customClaims.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
