package cmd

import (
	"VulnFusion/internal/log"
	"VulnFusion/internal/scanner"
	"VulnFusion/internal/utils"
	"flag"
	"fmt"
	"strings"
)

func StartCLI(args []string) {
	log.Info("[CLI] 命令行模式启动")

	fs := flag.NewFlagSet("scan", flag.ExitOnError)

	var singleTarget string
	var listFile string
	var name string
	var concurrency int

	fs.StringVar(&singleTarget, "u", "", "单个目标地址（如 http://example.com）")
	fs.StringVar(&listFile, "t", "", "目标地址列表文件路径")
	fs.StringVar(&name, "n", "", "模板文件名（如 test.yaml）")
	fs.IntVar(&concurrency, "c", 5, "并发数（默认5）")

	if err := fs.Parse(args); err != nil {
		log.Error("参数解析失败: %v", err)
		return
	}

	if name == "" || (singleTarget == "" && listFile == "") || (singleTarget != "" && listFile != "") {
		log.Warn("缺少必要参数，或 -u 与 -t 同时指定")
		fmt.Println("Usage:\n  ./VulnFusion scan -u http://target.com -n test.yaml\n  ./VulnFusion scan -t targets.txt -n test.yaml [-c 10]")
		return
	}

	var targets []string
	if singleTarget != "" {
		targets = []string{strings.TrimSpace(singleTarget)}
	} else {
		var err error
		targets, err = utils.LoadTargetsFromFile(listFile)
		if err != nil {
			log.Error("读取目标文件失败: %v", err)
			return
		}
	}

	task := scanner.Task{
		Targets:     targets,
		Name:        name,
		Concurrency: concurrency,
	}

	if err := scanner.Dispatch(task); err != nil {
		log.Error("任务执行失败: %v", err)
	}
}
