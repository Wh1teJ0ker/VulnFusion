package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"VulnFusion/internal/log"
)

var nucleiPath = "./data/bin/nuclei"

// InitNuclei 初始化 nuclei 环境（自动安装 + 版本检测）
func InitNuclei() error {
	absPath, _ := filepath.Abs(nucleiPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Warn("未检测到 nuclei，将自动安装到本地路径：%s", absPath)

		// 创建 data/bin 目录
		if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
			log.Error("创建 nuclei 安装目录失败: %v", err)
			return err
		}

		// 设置 GOBIN 路径并执行安装
		cmd := exec.Command("go", "install", "-v", "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest")
		cmd.Env = append(os.Environ(), fmt.Sprintf("GOBIN=%s", filepath.Dir(absPath)))

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("nuclei 安装失败: %v\n%s", err, string(output))
			return err
		}
		log.Info("nuclei 安装成功，路径：%s", absPath)
	} else {
		log.Info("已检测到 nuclei 可执行文件：%s", absPath)
	}

	// 打印 nuclei 版本
	cmd := exec.Command(absPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("获取 nuclei 版本失败: %v", err)
		return err
	}
	log.Info("nuclei 版本：%s", strings.TrimSpace(string(output)))
	return nil
}

// GetNucleiPath 返回 nuclei 执行路径
func GetNucleiPath() string {
	return nucleiPath
}
