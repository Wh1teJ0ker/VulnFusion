package nuclei

import (
	"VulnFusion/internal/log"
	"os/exec"
	"strconv"
)

// Options 扫描配置结构体，可按需扩展
type Options struct {
	Target     string   // 单个目标
	ListFile   string   // 目标列表文件
	Templates  []string // 模板路径
	Tags       []string // 标签过滤
	Output     string   // 输出文件路径
	JsonExport string   // JSON 导出路径
	Severity   []string // 严重程度过滤
	Silent     bool     // 是否静默输出
	Timeout    int      // 超时时间（秒）
}

// ScanWithOptions 使用结构体组合参数执行 nuclei 扫描
func ScanWithOptions(opt Options) error {
	args := buildArgs(opt)
	cmd := exec.Command(GetNucleiPath(), args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("nuclei 执行失败: %v\n%s", err, string(output))
		return err
	}
	log.Info("nuclei 扫描完成:\n%s", string(output))
	return nil
}

// buildArgs 将 Options 结构体转为 nuclei 命令行参数
func buildArgs(opt Options) []string {
	var args []string

	if opt.Target != "" {
		args = append(args, "-u", opt.Target)
	} else if opt.ListFile != "" {
		args = append(args, "-l", opt.ListFile)
	}

	for _, t := range opt.Templates {
		args = append(args, "-t", t)
	}
	for _, tag := range opt.Tags {
		args = append(args, "-tags", tag)
	}
	for _, sev := range opt.Severity {
		args = append(args, "-s", sev)
	}
	if opt.Output != "" {
		args = append(args, "-o", opt.Output)
	}
	if opt.JsonExport != "" {
		args = append(args, "-json-export", opt.JsonExport)
	}
	if opt.Silent {
		args = append(args, "-silent")
	}
	if opt.Timeout > 0 {
		args = append(args, "-timeout", strconv.Itoa(opt.Timeout))
	}
	return args
}
