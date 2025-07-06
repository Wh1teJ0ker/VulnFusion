package api

import (
	"VulnFusion/internal/log"
	"VulnFusion/internal/storage"
	"VulnFusion/web/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const tokenExpireSeconds = 3600

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// HandleRegister 用户注册处理
func HandleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("[Register] 参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if storage.UserExists(req.Username) {
		log.Warn("[Register] 用户已存在: %s", req.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "用户已存在"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("[Register] 密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误"})
		return
	}

	if err := storage.CreateUser(req.Username, string(hashedPwd)); err != nil {
		log.Error("[Register] 创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	log.Info("[Register] 注册成功: %s", req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleLogin 用户登录处理（JWT 写入 HttpOnly Cookie，前端不返回 token）
func HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("[Login] 参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := storage.GetUserByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		log.Warn("[Login] 用户名或密码错误: %s", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := middleware.GenerateTokenWithTTL(user.ID, tokenExpireSeconds)
	if err != nil {
		log.Error("[Login] 生成 Token 失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		tokenExpireSeconds, // ⏱ 3600 秒
		"/",
		"",    // domain
		false, // secure，可按是否 HTTPS 修改
		true,  // httpOnly
	)

	log.Info("[Login] 登录成功: %s", req.Username)
	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}
