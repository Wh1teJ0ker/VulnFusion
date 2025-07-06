package main

import (
	"VulnFusion/cmd"
	"VulnFusion/internal/config"
	"VulnFusion/internal/log"
	"VulnFusion/internal/nuclei"
	"VulnFusion/internal/storage"
	"VulnFusion/web/middleware"
	"os"
	"path/filepath"
)

// 打印使用说明
func printUsage() {
	log.Info(`Usage:
  ./VulnFusion web       # 启动 Web UI
  ./VulnFusion scan ...  # 命令行扫描
  ./VulnFusion init      # 初始化依赖`)
}

// 初始化数据目录、日志、配置与数据库
func initDataDir() {
	// 创建日志目录
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		log.Fatal("创建日志目录失败: %v", err)
	}

	// 初始化日志系统
	logPath := filepath.Join(config.LogDir, "vulnfusion.log")
	if err := log.InitLogger(logPath, "info"); err != nil {
		log.Fatal("日志系统初始化失败: %v", err)
	}
	log.Info("日志系统初始化完成")

	// JWT密钥管理
	middleware.InitJWTSecret()
	log.Info("JWT 密钥初始化完成")

	// 加载配置
	config.LoadConfig()
	log.Info("配置加载完成")

	// 初始化数据库
	storage.InitDB()
	log.Info("数据库初始化完成")
	log.Info("初始化目录完成: %s", config.BaseDataDir)
}

func main() {
	initDataDir()

	if len(os.Args) < 2 {
		log.Warn("缺少启动参数")
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "web":
		log.Info("启动 Web 模式")
		cmd.StartWeb()

	case "scan":
		log.Info("启动命令行扫描")
		cmd.StartCLI(os.Args[2:])

	case "init":
		log.Info("执行初始化任务")
		if err := nuclei.InitNuclei(); err != nil {
			log.Fatal("nuclei 初始化失败: %v", err)
		}
		log.Info("nuclei 初始化完成")

	default:
		log.Warn("未知启动模式: %s", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}
