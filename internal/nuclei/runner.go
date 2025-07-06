package nuclei

import (
	"VulnFusion/internal/config"
	"VulnFusion/internal/log"
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

// Options 简化后的扫描配置
type Options struct {
	Target    string   // 单个目标
	Templates []string // 模板路径（必须）
}

type Result struct {
	Raw        string `json:"-"`
	TemplateID string `json:"template.go-id"`
	MatchedAt  string `json:"matched-at"`
	Info       struct {
		Name     string `json:"name"`
		Severity string `json:"severity"`
	} `json:"info"`
}

// ScanWithOptions 执行 nuclei 并返回每一条 JSON 格式的扫描结果
func ScanWithOptions(opt Options) ([]Result, error) {
	args := buildArgs(opt)
	cmd := exec.Command(config.NucleiBin, args...)

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	output := outBuf.String()

	if err != nil {
		log.Error("nuclei 执行失败: %v\n%s", err, output)
		return nil, err
	}

	var results []Result
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var res Result
		if err := json.Unmarshal([]byte(line), &res); err == nil {
			res.Raw = line
			//log.Info("[漏洞] URL: %s | 名称: %s | 危险等级: %s", res.MatchedAt, res.Info.Name, res.Info.Severity)
			results = append(results, res)
		}
	}

	return results, nil
}

// buildArgs 构造 nuclei 命令参数
func buildArgs(opt Options) []string {
	var args []string

	if opt.Target != "" {
		args = append(args, "-u", opt.Target)
	}
	for _, t := range opt.Templates {
		args = append(args, "-t", t)
	}

	args = append(args, "-silent", "-jsonl")
	return args
}
