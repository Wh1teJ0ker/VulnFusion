package main

import (
	"os"
	"path/filepath"

	"VulnFusion/cmd"
	"VulnFusion/internal/log"
)

func initDataDir() {
	baseDir := "data"
	logDir := filepath.Join(baseDir, "logs")

	// 创建 data 和 logs 目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("创建数据目录失败: %v", err)
	}

	// 初始化日志系统
	logPath := filepath.Join(logDir, "vulnfusion.log")
	if err := log.InitLogger(logPath, "info"); err != nil {
		log.Fatal("日志系统初始化失败: %v", err)
	}

	log.Info("初始化目录完成: %s", baseDir)
}

func main() {
	// 初始化目录和日志
	initDataDir()

	// 检查启动参数
	if len(os.Args) < 2 {
		log.Warn("启动参数缺失")
		log.Info("Usage:\n  ./VulnFusion web       # 启动 Web UI\n  ./VulnFusion scan ...  # 命令行扫描")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "web":
		log.Info("启动 Web 模式")
		cmd.StartWeb()
	case "scan":
		log.Info("启动命令行扫描")
		cmd.StartCLI(os.Args[2:])
	default:
		log.Warn("未知模式: %s", os.Args[1])
		log.Info("Usage:\n  ./VulnFusion web       # 启动 Web UI\n  ./VulnFusion scan ...  # 命令行扫描")
	}
}
