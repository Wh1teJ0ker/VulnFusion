package middleware

import (
	"VulnFusion/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTAuthMiddleware 验证 Access Token 并注入 claims
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供 token"})
			return
		}

		// ✅ 去除 Bearer 前缀
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效或过期的 token"})
			return
		}
		InjectClaimsToContext(ctx, claims)
		ctx.Next()
	}
}

// RequireAdmin 仅允许管理员角色访问
func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			return
		}
		if claims.(*auth.CustomClaims).Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足，需管理员身份"})
			return
		}
		ctx.Next()
	}
}

// RequireRole 允许多个角色访问，只要命中任一即可通过
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			return
		}
		userRole := claims.(*auth.CustomClaims).Role
		for _, role := range roles {
			if userRole == role {
				ctx.Next()
				return
			}
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "权限不足"})
	}
}

// InjectClaimsToContext 将解析后的 claims 注入 Gin Context
func InjectClaimsToContext(ctx *gin.Context, claims *auth.CustomClaims) {
	ctx.Set("claims", claims)
}
