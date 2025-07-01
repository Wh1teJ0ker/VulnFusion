package plugin

import (
	"VulnFusion/internal/log"
	"VulnFusion/internal/plugin/docker"
	"fmt"
)

// Run 执行指定插件容器，使用目标作为参数
func Run(pluginName string, target string) error {
	imageName := fmt.Sprintf("vulnfusion_plugin_%s", pluginName)

	log.Info("准备运行插件: %s，目标: %s", pluginName, target)

	// 运行容器
	if err := docker.Run(imageName, target); err != nil {
		log.Error("插件执行失败: %v", err)
		return err
	}

	log.Info("插件 %s 执行完成", pluginName)
	return nil
}
