package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"VulnFusion/internal/config"
	"VulnFusion/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
)

type CustomClaims struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// JWTConfig 用于加载或写入 config.yaml 中的 jwt_secret 字段
type JWTConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
}

// InitJWTSecret 初始化 JWT 密钥（优先使用配置文件，失败则生成并写入）
func InitJWTSecret() {
	// 尝试读取配置文件
	if _, err := os.Stat(config.ConfigPath); err == nil {
		data, err := os.ReadFile(config.ConfigPath)
		if err == nil {
			var cfg JWTConfig
			if err := yaml.Unmarshal(data, &cfg); err == nil && cfg.JWTSecret != "" {
				jwtSecret = []byte(cfg.JWTSecret)
				log.Info("[Auth] 从 config.yaml 加载 JWT 密钥成功")
				return
			}
		}
		log.Warn("[Auth] 尝试读取 config.yaml 中 JWT 密钥失败，将生成新密钥")
	} else {
		log.Warn("[Auth] 未找到 config.yaml 文件，将生成新密钥")
	}

	// 生成强随机密钥
	secret, err := generateSecureRandom(32)
	if err != nil {
		log.Fatal("生成随机 JWT 密钥失败: %v", err)
	}
	jwtSecret = []byte(secret)
	log.Info("[Auth] 已生成临时 JWT 密钥")

	// 写入 config.yaml 文件
	newCfg := JWTConfig{JWTSecret: secret}
	yamlData, err := yaml.Marshal(&newCfg)
	if err != nil {
		log.Warn("[Auth] 配置序列化失败: %v", err)
		return
	}

	if err := os.WriteFile(config.ConfigPath, yamlData, 0644); err != nil {
		log.Warn("[Auth] 写入 config.yaml 失败: %v", err)
		return
	}
	log.Info("[Auth] 已将 JWT 密钥写入 config.yaml")
}

// AuthMiddleware 鉴权中间件，解析 JWT 并注入用户信息
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "缺少 Authorization 头"})
			return
		}

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "无效或过期的 token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func GenerateToken(userID string, role string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// generateSecureRandom 生成指定长度的 Base64 随机字符串
func generateSecureRandom(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("无法读取随机源")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateTokenWithTTL(userID interface{}, ttlSeconds int) (string, error) {
	claims := jwt.MapClaims{
		"user": userID,
		"exp":  time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
