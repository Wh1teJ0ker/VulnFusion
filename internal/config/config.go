package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	JWTSecret string `yaml:"jwt_secret"`
	// 其他配置项
}

var GlobalConfig Config

func LoadConfig() error {
	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &GlobalConfig)
}
