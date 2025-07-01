package main

import (
	"VulnFusion/internal/log"
)

func main() {
	// 初始化日志系统（放在程序最早位置）
	err := log.InitLogger("data/logs/vulnfusion.log", "debug")
	if err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}

	// 调用日志输出函数
	log.Debug("调试模式启动，参数为：%v", []string{"web", "cli"})
	log.Info("程序初始化完成")
	log.Warn("配置文件缺失，使用默认配置")
	log.Error("连接数据库失败: %s", "timeout")
	log.Fatal("致命错误：%s", "配置非法")
}
