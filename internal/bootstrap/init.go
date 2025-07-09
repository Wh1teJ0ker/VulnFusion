package bootstrap

import (
	"VulnFusion/internal/config"
	"VulnFusion/internal/db"
	"VulnFusion/internal/log"
	"VulnFusion/internal/models"
	"VulnFusion/internal/scanner"
	"VulnFusion/internal/utils"
)

// InitializeSystem 执行系统初始化流程，包括日志、数据库、管理员账号、nuclei 初始化等
func InitializeSystem() error {
	// 初始化日志系统
	log.InitLogger("dev", "debug")

	// 初始化数据库
	dbPath := config.GetDBPath()
	_, err := db.InitDatabase(dbPath)
	if err != nil {
		log.Error("数据库初始化失败: %v", err)
		return err
	}

	// 初始化管理员账号
	if err := InitializeAdmin(); err != nil {
		log.Error("初始化管理员失败: %v", err)
		return err
	}

	// 初始化 nuclei 扫描器
	if err := scanner.InitNuclei(); err != nil {
		log.Error("初始化 nuclei 环境失败: %v", err)
		return err
	}

	log.Info("系统初始化完成")
	return nil
}

// InitializeAdmin 检查是否存在管理员账号，不存在则创建默认账号
func InitializeAdmin() error {
	admin, err := models.GetUserByUsername("admin")
	if err == nil && admin != nil {
		log.Info("已存在管理员账号，无需初始化")
		return nil
	}

	log.Info("未检测到管理员账号，正在创建默认管理员账户")
	hashed, err := utils.HashPassword("admin123")
	if err != nil {
		log.Error("密码加密失败: %v", err)
		return err
	}

	user := &models.User{
		Username: "admin",
		Password: hashed,
		Role:     "admin",
	}

	err = models.CreateUser(user)
	if err != nil {
		log.Error("创建管理员账号失败: %v", err)
		return err
	}

	log.Info("默认管理员账号创建成功：用户名 admin，密码 admin123")
	return nil
}
