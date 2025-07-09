package auth

import (
	"VulnFusion/internal/config"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"sync"
)

type CustomClaims struct {
	UserID   uint   `json:"user_id"` // ✅ 改为 uint
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte
var blacklist = make(map[string]int64)
var blacklistMu sync.RWMutex
var once sync.Once

func LoadOrGenerateJWTSecret() ([]byte, error) {
	var err error
	once.Do(func() {
		cfgSecret := config.GetJWTSecret()
		if cfgSecret != "" {
			jwtSecret = []byte(cfgSecret)
		} else {
			secret := make([]byte, 32)
			_, e := rand.Read(secret)
			if e != nil {
				err = e
				return
			}
			jwtSecret = secret
		}
	})
	if err != nil {
		return nil, err
	}
	return jwtSecret, nil
}

func GenerateToken(userID uint, username string, role string, duration time.Duration) (string, error) {
	secret, err := LoadOrGenerateJWTSecret()
	if err != nil {
		return "", err
	}
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        generateJTI(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(userID uint, duration time.Duration) (string, error) {
	secret, err := LoadOrGenerateJWTSecret()
	if err != nil {
		return "", err
	}
	claims := jwt.RegisteredClaims{
		Subject:   base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(userID)))),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        generateJTI(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	secret, err := LoadOrGenerateJWTSecret()
	if err != nil {
		return nil, err
	}
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	if IsTokenBlacklisted(claims.ID) {
		return nil, errors.New("token is blacklisted")
	}
	return claims, nil
}

func IsTokenBlacklisted(jti string) bool {
	blacklistMu.RLock()
	expiry, exists := blacklist[jti]
	blacklistMu.RUnlock()
	return exists && time.Now().Unix() < expiry
}

func AddTokenToBlacklist(jti string, expiration time.Duration) error {
	blacklistMu.Lock()
	blacklist[jti] = time.Now().Add(expiration).Unix()
	blacklistMu.Unlock()
	return nil
}

func generateJTI() string {
	raw := make([]byte, 16)
	_, _ = rand.Read(raw)
	return base64.URLEncoding.EncodeToString(raw)
}
