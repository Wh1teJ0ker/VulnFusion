package scanner

import (
	"VulnFusion/internal/log"
	"encoding/json"
	"strings"
)

type NucleiResult struct {
	TemplateID string `json:"template.go-id"`
	MatchedAt  string `json:"matched-at"`
	Host       string `json:"host"`
	Severity   string `json:"info.severity"`
}

func PrintNucleiResults(jsonStr string) {
	lines := strings.Split(jsonStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			log.Warn("JSON 解析失败: %v", err)
			continue
		}

		log.Info("[Result] %s  %s  %s",
			result["template.go-id"],
			result["matched-at"],
			result["info"].(map[string]interface{})["severity"],
		)
	}
}
