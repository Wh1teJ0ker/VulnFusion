package scanner

import (
	"VulnFusion/internal/config"
	"errors"
	"os/exec"
	"strings"

	"VulnFusion/internal/log"
)

// ScanOptions 描述 nuclei 扫描参数
type ScanOptions struct {
	Target     string   // 单个目标地址（-u）或文件路径（-l）
	Template   string   // 模板路径或目录（-t）
	Silent     bool     // 静默模式（-silent）
	JsonOutput bool     // 是否启用 JSONL 输出（-jsonl）
	CustomArgs []string // 用户自定义附加参数
}

// RunScanTask 执行 nuclei 扫描任务并返回 stdout 输出
func RunScanTask(options ScanOptions) ([]byte, error) {
	if err := ValidateScanOptions(options); err != nil {
		log.Error("参数校验失败: %v", err)
		return nil, err
	}

	cmd, err := BuildNucleiCommand(options)
	if err != nil {
		log.Error("命令构建失败: %v", err)
		return nil, err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("nuclei 执行失败: %v\n%s", err, string(output))
		return output, err
	}

	log.Info("nuclei 扫描完成，输出结果: %s", string(output))
	return output, nil
}

// BuildNucleiCommand 根据参数构造 nuclei 命令
func BuildNucleiCommand(options ScanOptions) (*exec.Cmd, error) {
	args := BuildCommandArgs(options)
	log.Info("构建命令参数: %v", args)
	cmd := exec.Command(GetNucleiPath(), args...)
	return cmd, nil
}

// BuildCommandArgs 构建 nuclei CLI 参数数组
func BuildCommandArgs(options ScanOptions) []string {
	var args []string

	if strings.HasPrefix(options.Target, "http") {
		args = append(args, "-u", options.Target)
	} else {
		args = append(args, "-l", options.Target)
	}

	if options.Template != "" {
		// 前端只传模板文件名，如 test.yaml
		templateDir := config.GetNucleiTemplatePath() // 比如返回 "./templates"
		fullPath := templateDir + "/" + options.Template
		args = append(args, "-t", fullPath)
	} else {
		defaultPath := config.GetNucleiTemplatePath() // 可保留默认模板兜底
		if defaultPath != "" {
			args = append(args, "-t", defaultPath)
		}
	}

	if options.JsonOutput {
		args = append(args, "-jsonl")
	}

	if options.Silent {
		args = append(args, "-silent")
	}

	if len(options.CustomArgs) > 0 {
		args = append(args, options.CustomArgs...)
	}

	return args
}

// ValidateScanOptions 检查参数是否合法
func ValidateScanOptions(opt ScanOptions) error {
	if opt.Target == "" {
		err := errors.New("Target 不能为空")
		log.Error(err.Error())
		return err
	}
	if opt.Template == "" {
		err := errors.New("Template 不能为空")
		log.Error(err.Error())
		return err
	}
	return nil
}
