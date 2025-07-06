package storage

// createTables 初始化数据库中的所有核心表结构
func createTables() {
	// 用户表：存储注册用户的信息
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,         -- 用户唯一 ID
		username TEXT NOT NULL UNIQUE,                -- 用户名，唯一
		password TEXT NOT NULL,                       -- 哈希加密后的密码
		role TEXT NOT NULL DEFAULT 'user',            -- 用户角色（如：admin/user/viewer）
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP -- 创建时间
	);`

	// 模板表：用于记录 YAML 模板的元信息，文件本体保存在磁盘
	createTemplateTable := `
	CREATE TABLE IF NOT EXISTS templates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,         -- 模板唯一 ID
		name TEXT NOT NULL UNIQUE,                    -- 模板名称（如 python-rce）
		path TEXT NOT NULL,                           -- YAML 文件路径
		description TEXT,                             -- 模板描述（如“Python 代码执行”）
		tags TEXT,                                     -- 模板标签（如 rce,test）
		file_hash TEXT,                                -- 文件哈希值（防止重复上传）
		uploaded_by INTEGER,                          -- 上传者的用户 ID
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 上传时间
		FOREIGN KEY(uploaded_by) REFERENCES users(id) -- 外键关联用户
	);`

	// 扫描任务表：每个任务记录一个目标和一个模板，执行后写入原始结果和摘要
	createScanTaskTable := `
	CREATE TABLE IF NOT EXISTS scan_tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,         -- 任务唯一 ID
		user_id INTEGER NOT NULL,                     -- 所属用户 ID
		target TEXT NOT NULL,                         -- 扫描目标（URL 或 IP）
		template_id INTEGER NOT NULL,                 -- 使用的模板 ID（引用 templates）
		status TEXT DEFAULT 'pending',                -- 任务状态：pending、running、success、failed
		result TEXT,                                   -- 原始 nuclei 扫描结果（JSON）
		summary TEXT,                                  -- 简要信息摘要（如高危数量等）
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
		FOREIGN KEY(user_id) REFERENCES users(id),     -- 外键关联用户
		FOREIGN KEY(template_id) REFERENCES templates(id) -- 外键关联模板
	);`

	// 执行建表语句
	statements := []string{
		createUserTable,
		createTemplateTable,
		createScanTaskTable,
	}

	for _, stmt := range statements {
		if _, err := DB.Exec(stmt); err != nil {
			panic("数据库建表失败: " + err.Error())
		}
	}
}
