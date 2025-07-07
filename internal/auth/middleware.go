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

		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			BuildErrorResponse(ctx, http.StatusUnauthorized, "TOKEN_INVALID", err.Error())
			ctx.Abort()
			return
		}

		InjectClaimsToContext(ctx, claims)
		ctx.Next()
	}
}

func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken, err := ExtractRefreshTokenFromRequest(ctx)
		if err != nil {
			BuildErrorResponse(ctx, http.StatusUnauthorized, "REFRESH_TOKEN_MISSING", err.Error())
			ctx.Abort()
			return
		}

		claims, err := auth.ParseToken(refreshToken)
		if err != nil {
			BuildErrorResponse(ctx, http.StatusUnauthorized, "REFRESH_TOKEN_INVALID", err.Error())
			ctx.Abort()
			return
		}

		// 自动生成新的 Access Token
		newToken, err := auth.GenerateToken(claims.UserID, claims.Username, claims.Role, auth.DefaultTokenTTL())
		if err != nil {
			BuildErrorResponse(ctx, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "生成新 token 失败")
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token": newToken,
		})
		ctx.Abort()
	}
}

func InjectClaimsToContext(ctx *gin.Context, claims *auth.CustomClaims) {
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
