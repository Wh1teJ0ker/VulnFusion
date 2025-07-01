package cmd

import (
	"VulnFusion/internal/log"
	"VulnFusion/internal/plugin"
	"VulnFusion/internal/scanner"
	"flag"
	"os"
)

// StartCLI 启动 CLI 模式并解析参数
func StartCLI(args []string) {
	log.Info("[CLI] 命令行模式启动")

	// 自定义 flag 集
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	target := flags.String("target", "", "目标地址，例如 https://example.com")
	taskType := flags.String("type", "plugin", "任务类型（plugin / nuclei）")
	name := flags.String("name", "", "模板名或插件名")
	list := flags.Bool("list", false, "列出可用插件")

	// 解析参数
	if err := flags.Parse(args); err != nil {
		log.Error("参数解析失败: %v", err)
		return
	}

	// 展示插件列表
	if *list {
		plugins, err := plugin.ListPlugins()
		if err != nil {
			log.Error("读取插件列表失败: %v", err)
			return
		}
		log.Info("可用插件列表：")
		for _, p := range plugins {
			log.Info(" - %s", p)
		}
		return
	}

	// 参数校验
	if *target == "" || *name == "" {
		log.Warn("参数缺失，请使用 -target 和 -name 指定目标和模板/插件名")
		flags.Usage()
		return
	}

	// 构建任务
	task := scanner.Task{
		Target: *target,
		Type:   *taskType,
		Name:   *name,
	}

	log.Info("[CLI] 调度任务: %+v", task)
	if err := scanner.Dispatch(task); err != nil {
		log.Error("任务执行失败: %v", err)
	} else {
		log.Info("任务执行完成")
	}
}
