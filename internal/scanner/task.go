package scanner

type Task struct {
	Targets     []string // 多个目标 URL
	Name        string   // 模板文件名（如 test.yaml）
	Concurrency int      // 并发数，0 表示默认
}
