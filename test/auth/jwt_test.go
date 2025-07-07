// jwt_test.go
package auth

import (
	"encoding/base64"
	"strconv"
	"testing"
	"time"

	auth "VulnFusion/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret, err := auth.LoadOrGenerateJWTSecret()
	assert.Nil(t, err)
	assert.NotEmpty(t, secret)

	tokenStr, err := auth.GenerateToken(1, "admin", "super", time.Minute*15)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := auth.ParseToken(tokenStr)
	assert.Nil(t, err)
	assert.Equal(t, 1, claims.UserID)
	assert.Equal(t, "admin", claims.Username)
	assert.Equal(t, "super", claims.Role)
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenStr, err := auth.GenerateRefreshToken(1, time.Hour*24*7)
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenStr)

	parsedToken, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.LoadOrGenerateJWTSecret()
	})
	assert.Nil(t, err)
	assert.True(t, parsedToken.Valid)

	claims := parsedToken.Claims.(*jwt.RegisteredClaims)

	// 解码 subject 为 userID
	decodedSub, err := base64.StdEncoding.DecodeString(claims.Subject)
	assert.Nil(t, err)

	userID, err := strconv.Atoi(string(decodedSub))
	assert.Nil(t, err)
	assert.Equal(t, 1, userID)
}

func TestBlacklistToken(t *testing.T) {
	jti := "test-jti"

	// 应该不在黑名单
	isBlacklisted := auth.IsTokenBlacklisted(jti)
	assert.False(t, isBlacklisted)

	// 加入黑名单
	err := auth.AddTokenToBlacklist(jti, time.Minute*30)
	assert.Nil(t, err)

	// 再次验证
	isBlacklisted = auth.IsTokenBlacklisted(jti)
	assert.True(t, isBlacklisted)
}
