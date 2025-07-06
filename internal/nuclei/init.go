package nuclei

import (
	"VulnFusion/internal/config"
	"VulnFusion/internal/log"
	"archive/zip"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var nucleiPath = config.NucleiBin

func InitNuclei() error {
	absPath, _ := filepath.Abs(nucleiPath)

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Warn("未检测到 nuclei，将尝试自动安装到: %s", absPath)

		if checkGoExists() {
			log.Info("检测到 go 环境，尝试使用 go install 安装 nuclei...")
			if installViaGo(absPath) == nil {
				return printVersion(absPath)
			}
			log.Warn("go install 安装失败，尝试使用 release 安装方式")
		} else {
			log.Warn("系统未安装 go，将直接使用 release 安装方式")
		}

		// 尝试使用 release 下载方式
		if err := installViaRelease(absPath); err != nil {
			log.Error("nuclei 安装失败（release 模式）: %v", err)
			return err
		}

		log.Info("nuclei 安装成功（release 模式），路径: %s", absPath)
	} else {
		log.Info("已检测到 nuclei：%s", absPath)
	}

	return printVersion(absPath)
}

// 尝试 go install 安装
func installViaGo(absPath string) error {
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		log.Error("创建目录失败: %v", err)
		return err
	}
	cmd := exec.Command("go", "install", "-v", "github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest")
	cmd.Env = append(os.Environ(), "GOBIN="+filepath.Dir(absPath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("go install nuclei 失败: %v", err)
		log.Debug("输出:\n%s", string(output))
		return err
	}
	return nil
}

// 使用 release zip 下载并解压 nuclei
func installViaRelease(destPath string) error {
	url := buildDownloadURL()
	log.Info("从 release 下载 nuclei: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tmpZip := filepath.Join(os.TempDir(), "nuclei_release.zip")
	out, err := os.Create(tmpZip)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	// 解压并找到 nuclei 可执行文件
	reader, err := zip.OpenReader(tmpZip)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Name == "nuclei" || f.Name == "nuclei.exe" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, 0755)
			if err != nil {
				return err
			}
			defer out.Close()
			io.Copy(out, rc)
			return nil
		}
	}
	return os.ErrNotExist
}

// 构造 GitHub Release 下载链接（适配平台）
func buildDownloadURL() string {
	platform := runtime.GOOS
	arch := runtime.GOARCH
	return "https://github.com/projectdiscovery/nuclei/releases/latest/download/nuclei_" + platform + "_" + arch + ".zip"
}

// 打印 nuclei 版本
func printVersion(path string) error {
	cmd := exec.Command(path, "-version")
	output, err := cmd.CombinedOutput()
	if err == nil {
		log.Info("nuclei 版本: %s", strings.TrimSpace(string(output)))
	}
	return err
}

func checkGoExists() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func GetNucleiPath() string {
	return nucleiPath
}
