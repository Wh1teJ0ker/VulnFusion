package auth

import (
	"VulnFusion/internal/auth"
	"VulnFusion/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExtractTokenFromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.Request.Header.Set("Authorization", "Bearer test-token-123")

	token, err := auth.ExtractTokenFromHeader(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "test-token-123", token)
}

func TestExtractTokenFromHeader_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{Header: make(http.Header)}
	ctx.Request.Header.Set("Authorization", "InvalidFormat")

	_, err := auth.ExtractTokenFromHeader(ctx)
	assert.NotNil(t, err)
}

func TestExtractRefreshTokenFromRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 设置 POST 表单值
	body := strings.NewReader("refresh_token=refresh123")
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx.Request = req
	_ = ctx.Request.ParseForm()

	token, err := auth.ExtractRefreshTokenFromRequest(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "refresh123", token)
}

func TestBuildErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	auth.BuildErrorResponse(ctx, http.StatusUnauthorized, "ERR_CODE", "Unauthorized")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "ERR_CODE")
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestJWTAuthMiddleware_Success(t *testing.T) {
	_ = config.LoadConfig("config.yaml") // 假设已存在

	// 生成 token
	token, _ := auth.GenerateToken(1, "admin", "super", time.Minute*5)

	r := gin.New()
	r.Use(auth.JWTAuthMiddleware())
	r.GET("/test", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

func TestJWTAuthMiddleware_Failure(t *testing.T) {
	r := gin.New()
	r.Use(auth.JWTAuthMiddleware())
	r.GET("/test", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	req := httptest.NewRequest("GET", "/test", nil) // 无 Authorization
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "TOKEN_MISSING")
}
