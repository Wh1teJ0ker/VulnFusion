package api

// RegisterRequest 用户注册请求参数
type RegisterRequest struct {
	Username string `json:"username" example:"admin"`  // 用户名
	Password string `json:"password" example:"123456"` // 密码
	Role     string `json:"role" example:"user"`       // 用户角色：user / admin
}

// LoginRequest 用户登录请求参数
type LoginRequest struct {
	Username string `json:"username" example:"admin"`  // 用户名
	Password string `json:"password" example:"123456"` // 密码
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"<refresh-token>"` // 刷新 token
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Target   string `json:"target" example:"https://example.com"` // 目标地址
	Template string `json:"template" example:"cves/2021/*.yaml"`  // Nuclei 模板路径
}

// BatchDeleteRequest 批量删除任务请求
type BatchDeleteRequest struct {
	IDs []uint `json:"ids"`
}

// UpdateStatusRequest 更新任务状态请求
type UpdateStatusRequest struct {
	ID     uint   `json:"id" example:"1"`        // 任务 ID
	Status string `json:"status" example:"done"` // 新状态：pending / running / done / failed
}
