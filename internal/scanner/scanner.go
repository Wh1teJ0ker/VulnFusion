package scanner

import (
	"VulnFusion/internal/config"
	"VulnFusion/internal/log"
	"VulnFusion/internal/nuclei"
	"path/filepath"
	"sync"
)

const defaultConcurrency = 5

func Dispatch(task Task) error {
	concurrency := task.Concurrency
	if concurrency <= 0 {
		concurrency = defaultConcurrency
	}

	log.Info("[Scanner] 启动并发扫描：%d 个目标，模板: %s，并发数: %d",
		len(task.Targets), task.Name, concurrency)

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	templatePath := filepath.Join(config.TemplateDir, task.Name)

	for _, target := range task.Targets {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			sem <- struct{}{}
			log.Info("[Scanner] 扫描目标: %s", t)

			results, err := nuclei.ScanWithOptions(nuclei.Options{
				Target:    t,
				Templates: []string{templatePath},
			})
			if err != nil {
				log.Error("[Scanner] 扫描失败: %s => %v", t, err)
				<-sem
				return
			}

			for _, res := range results {
				// 控制台保留原始 Raw（可替换为保存功能）
				//fmt.Println(res.Raw)

				// 日志输出关键信息
				log.Info("[漏洞] URL: %s | 名称: %s | 危险等级: %s", res.MatchedAt, res.Info.Name, res.Info.Severity)
			}

			<-sem
		}(target)
	}

	wg.Wait()
	log.Info("[Scanner] 所有扫描任务完成")
	return nil
}
