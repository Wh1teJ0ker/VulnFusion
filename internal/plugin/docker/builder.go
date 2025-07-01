package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"os"
	"path/filepath"

	"VulnFusion/internal/log"
)

// BuildImage 构建插件镜像（从指定目录读取 Dockerfile + 源码）
func BuildImage(pluginName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Error("创建 Docker 客户端失败: %v", err)
		return err
	}

	pluginDir := filepath.Join("data", "plugins", "vulnfusion_plugin_"+pluginName)
	templateDockerfile := filepath.Join("data", "plugins", "python", "Dockerfile")
	targetDockerfile := filepath.Join(pluginDir, "Dockerfile")
	pluginSrc := filepath.Join(pluginDir, pluginName+".py")
	pluginDst := filepath.Join(pluginDir, "plugin.py") // 镜像入口文件固定为 plugin.py

	// 自动复制模板 Dockerfile
	if _, err := os.Stat(targetDockerfile); os.IsNotExist(err) {
		src, err := os.ReadFile(templateDockerfile)
		if err != nil {
			log.Error("无法读取模板 Dockerfile: %v", err)
			return err
		}
		if err := os.WriteFile(targetDockerfile, src, 0644); err != nil {
			log.Error("写入 Dockerfile 失败: %v", err)
			return err
		}
		log.Info("已复制模板 Dockerfile 到插件目录")
	}

	// 自动复制并重命名插件脚本
	if _, err := os.Stat(pluginDst); os.IsNotExist(err) {
		src, err := os.ReadFile(pluginSrc)
		if err != nil {
			log.Error("无法读取插件脚本: %v", err)
			return err
		}
		if err := os.WriteFile(pluginDst, src, 0644); err != nil {
			log.Error("写入 plugin.py 失败: %v", err)
			return err
		}
		log.Info("已复制插件脚本为 plugin.py")
	}

	// 构建 tar 上下文
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	files := []string{"Dockerfile", "plugin.py"}

	for _, fname := range files {
		fullPath := filepath.Join(pluginDir, fname)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			log.Error("读取文件失败 [%s]: %v", fullPath, err)
			return err
		}

		info, err := os.Stat(fullPath)
		if err != nil {
			log.Error("获取文件信息失败 [%s]: %v", fullPath, err)
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			log.Error("创建 tar header 失败 [%s]: %v", fname, err)
			return err
		}
		header.Name = fname

		if err := tw.WriteHeader(header); err != nil {
			log.Error("写入 tar header 失败 [%s]: %v", fname, err)
			return err
		}
		if _, err := tw.Write(data); err != nil {
			log.Error("写入文件数据失败 [%s]: %v", fname, err)
			return err
		}
	}
	_ = tw.Close()

	imageName := fmt.Sprintf("vulnfusion_plugin_%s", pluginName)
	buildResp, err := cli.ImageBuild(ctx, &buf, types.ImageBuildOptions{
		Tags:        []string{imageName},
		Remove:      true,
		ForceRemove: true,
	})
	if err != nil {
		log.Error("镜像构建失败: %v", err)
		return err
	}
	defer buildResp.Body.Close()

	log.Info("正在构建镜像 %s...", imageName)
	_, _ = io.Copy(os.Stdout, buildResp.Body)
	log.Info("镜像构建完成: %s", imageName)

	return nil
}
