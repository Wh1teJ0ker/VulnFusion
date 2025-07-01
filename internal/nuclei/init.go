package nuclei

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"VulnFusion/internal/log"
)

var nucleiPath = "./data/bin/nuclei"

// InitNuclei 初始化 nuclei：优先使用项目内的 nuclei，若未安装则自动安装
func InitNuclei() error {
	absPath, _ := filepath.Abs(nucleiPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Warn("未检测到 nuclei，将自动安装到本地路径：%s", absPath)

		// 创建 data/bin 目录
		if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
			log.Error("创建本地 nuclei 安装目录失败: %v", err)
			return err
		}

		// 设置 GOPATH 路径，并安装 nuclei 到本地路径
		cmd := exec.Command("go", "install", "-v", "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest")
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("GOBIN=%s", filepath.Dir(absPath)),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("nuclei 安装失败: %v\n%s", err, string(output))
			return err
		}

		log.Info("nuclei 安装成功（私有路径），路径：%s", absPath)
	} else {
		log.Info("检测到 nuclei（私有路径）: %s", absPath)
	}

	// 打印版本
	cmd := exec.Command(absPath, "-version")
	output, err := cmd.CombinedOutput()
	if err == nil {
		version := strings.TrimSpace(string(output))
		log.Info("nuclei 版本: %s", version)
	}
	return nil
}

// GetNucleiPath 返回当前 nuclei 可执行文件路径
func GetNucleiPath() string {
	return nucleiPath
}
