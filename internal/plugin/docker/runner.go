package docker

import (
	"VulnFusion/internal/log"
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"strings"
)

// Run 执行插件（自动构建镜像 + 运行容器）
func Run(pluginName string, target string) error {
	log.Info("[Plugin] 准备运行插件: %s", pluginName)

	// 提取插件名，去除 .py
	name := strings.TrimSuffix(pluginName, ".py")
	imageName := fmt.Sprintf("vulnfusion_plugin_%s", name)

	// 构建镜像（如果镜像不存在）
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Error("Docker 客户端创建失败: %v", err)
		return err
	}

	if !checkImageExists(cli, ctx, imageName) {
		log.Warn("镜像 %s 不存在，开始构建...", imageName)
		if err := BuildImage(name); err != nil {
			log.Error("镜像构建失败: %v", err)
			return err
		}
		log.Info("镜像构建完成: %s", imageName)
	} else {
		log.Debug("镜像已存在: %s", imageName)
	}

	// 创建容器
	log.Info("开始创建容器，镜像: %s，参数: %s", imageName, target)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{target},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		log.Error("容器创建失败: %v", err)
		return err
	}
	log.Debug("容器创建成功: %s", resp.ID)

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Error("容器启动失败: %v", err)
		return err
	}
	log.Info("容器启动成功: %s", resp.ID)

	// 读取容器日志
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		log.Error("日志读取失败: %v", err)
		return err
	}
	defer out.Close()

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		log.Info("[插件输出] %s", scanner.Text())
	}

	// 等待容器退出
	waitCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case result := <-waitCh:
		log.Info("插件容器退出，状态码: %d", result.StatusCode)
	case err := <-errCh:
		log.Error("等待容器退出失败: %v", err)
	}

	// 清理容器
	if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Warn("容器清理失败: %v", err)
	} else {
		log.Debug("容器已清理: %s", resp.ID)
	}

	log.Info("[Plugin] 插件执行完成: %s", pluginName)
	return nil
}

// checkImageExists 检查镜像是否存在
func checkImageExists(cli *client.Client, ctx context.Context, imageName string) bool {
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		log.Warn("检查镜像存在性失败: %v", err)
		return false
	}
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				return true
			}
		}
	}
	return false
}
