package plugin

import (
	"os"
)

// ListPlugins 列出所有插件名称
func ListPlugins() ([]string, error) {
	entries, err := os.ReadDir("./data/plugins")
	if err != nil {
		return nil, err
	}
	var result []string
	for _, entry := range entries {
		if !entry.IsDir() {
			result = append(result, entry.Name())
		}
	}
	return result, nil
}
