package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	AppName string `yaml:"app_name"`
	Port    int    `yaml:"port"`
	Env     string `yaml:"env"`

	JWT struct {
		Secret          string        `yaml:"secret"`
		AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
		RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	} `yaml:"jwt"`

	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"admin"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
}

var Global Config

// LoadConfig 从指定路径加载 config.yaml
func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &Global)
}

// DefaultTokenTTL 获取默认的 AccessToken 有效期
func DefaultTokenTTL() time.Duration {
	if Global.JWT.AccessTokenTTL > 0 {
		return Global.JWT.AccessTokenTTL
	}
	return 15 * time.Minute
}

// DefaultRefreshTokenTTL 获取默认的 RefreshToken 有效期
func DefaultRefreshTokenTTL() time.Duration {
	if Global.JWT.RefreshTokenTTL > 0 {
		return Global.JWT.RefreshTokenTTL
	}
	return 7 * 24 * time.Hour
}

// GetJWTSecret 返回 JWT 密钥字符串
func GetJWTSecret() string {
	return Global.JWT.Secret
}

// GetDBPath 返回 SQLite 数据库路径
func GetDBPath() string {
	return Global.Database.Path
}
