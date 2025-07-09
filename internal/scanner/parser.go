package scanner

import (
	"bufio"
	"encoding/json"
	"errors"
	"strings"

	"VulnFusion/internal/log"
)

// Result 描述 nuclei 扫描输出中的一条漏洞信息
type Result struct {
	TemplateID string `json:"templateID"`
	Info       struct {
		Name     string   `json:"name"`
		Severity string   `json:"severity"`
		Tags     []string `json:"tags"`
	} `json:"info"`
	Matched   string `json:"matched-at"`
	Timestamp string `json:"timestamp"`
}

// ParseNucleiResult 解析 nuclei JSONL 输出，返回结构化结果数组
func ParseNucleiResult(raw []byte) ([]Result, error) {
	var results []Result

	scanner := bufio.NewScanner(strings.NewReader(string(raw)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var r Result
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			log.Error("解析 JSONL 行失败: %v\n内容: %s", err, line)
			continue
		}
		results = append(results, r)
	}

	if err := scanner.Err(); err != nil {
		log.Error("读取 JSONL 输出失败: %v", err)
		return nil, err
	}
	if len(results) == 0 {
		err := errors.New("未解析出任何有效结果")
		log.Warn(err.Error())
		return nil, err
	}

	log.Info("成功解析 %d 条漏洞结果", len(results))
	return results, nil
}
