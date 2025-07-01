// scanner/scanner.go

package scanner

import (
	"VulnFusion/internal/log"
	"VulnFusion/internal/nuclei"
	"VulnFusion/internal/plugin"
	"path/filepath"
)

// Task 通用扫描任务结构体（支持插件或 nuclei 模式）
type Task struct {
	Target        string          // 扫描目标
	Type          string          // plugin / nuclei
	Name          string          // 插件名称或 nuclei 模板名
	NucleiOptions *nuclei.Options // 可选：用于扩展自定义 nuclei 参数
}

// Dispatch 分发调度任务，自动适配 nuclei 或插件
func Dispatch(task Task) error {
	log.Info("[Scanner] 接收到任务：Target=%s Type=%s Name=%s", task.Target, task.Type, task.Name)

	switch task.Type {
	case "nuclei":
		log.Info("[Scanner] 启动 nuclei 模板扫描：%s", task.Name)

		var opt nuclei.Options
		if task.NucleiOptions != nil {
			// 使用自定义 options
			opt = *task.NucleiOptions
		} else {
			// 默认构造 options
			opt = nuclei.Options{
				Target:    task.Target,
				Templates: []string{filepath.Join("./data/templates", task.Name)},
				Silent:    true,
				Timeout:   10,
			}
		}

		if err := nuclei.ScanWithOptions(opt); err != nil {
			log.Error("[Scanner] nuclei 扫描失败: %v", err)
			return err
		}

	case "plugin":
		log.Info("[Scanner] 启动插件任务：%s", task.Name)
		if err := plugin.Run(task.Name, task.Target); err != nil {
			log.Error("[Scanner] 插件执行失败: %v", err)
			return err
		}

	default:
		log.Warn("[Scanner] 不支持的任务类型: %s", task.Type)
	}
	return nil
}
