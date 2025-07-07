package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ExtractTokenFromHeader(ctx)
		if err != nil {
			BuildErrorResponse(ctx, http.StatusUnauthorized, "TOKEN_MISSING", err.Error())
			ctx.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			BuildErrorResponse(ctx, http.StatusUnauthorized, "TOKEN_INVALID", err.Error())
			ctx.Abort()
			return
		}

		InjectClaimsToContext(ctx, claims)
		ctx.Next()
	}
}

func InjectClaimsToContext(ctx *gin.Context, claims *CustomClaims) {
	ctx.Set("user_id", claims.UserID)
	ctx.Set("username", claims.Username)
	ctx.Set("role", claims.Role)
}

func ExtractTokenFromHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", gin.Error{Err: http.ErrNoCookie, Type: gin.ErrorTypeBind, Meta: "Authorization 头格式应为 Bearer <token>"}
	}
	return authHeader[7:], nil
}

func ExtractRefreshTokenFromRequest(ctx *gin.Context) (string, error) {
	if token := ctx.PostForm("refresh_token"); token != "" {
		return token, nil
	}
	if token := ctx.GetHeader("X-Refresh-Token"); token != "" {
		return token, nil
	}
	if cookie, err := ctx.Cookie("refresh_token"); err == nil {
		return cookie, nil
	}
	return "", gin.Error{Err: http.ErrNoCookie, Type: gin.ErrorTypeBind, Meta: "未提供有效的 refresh token"}
}

func BuildErrorResponse(ctx *gin.Context, statusCode int, code string, message string) {
	ctx.JSON(statusCode, gin.H{
		"code":    code,
		"message": message,
	})
}
