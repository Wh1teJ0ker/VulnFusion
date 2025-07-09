package models

import (
	"VulnFusion/internal/db"
	"time"
)

// 任务状态常量定义
const (
	StatusPending = "pending"
	StatusRunning = "running"
	StatusDone    = "done"
	StatusFailed  = "failed"
)

type Task struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`        // 所属用户
	Target    string    `gorm:"not null"`        // 扫描目标
	Template  string    `gorm:"not null"`        // nuclei 模板名称
	CreatedAt time.Time `gorm:"autoCreateTime"`  // 创建时间
	Status    string    `gorm:"default:pending"` // 状态：pending、running、done、failed
}

// CreateTask 创建新任务记录
func CreateTask(task *Task) error {
	return db.GetDB().Create(task).Error
}

// GetTaskByID 根据任务 ID 查询任务详情
func GetTaskByID(id uint) (*Task, error) {
	var task Task
	if err := db.GetDB().First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// ListTasksByUserID 列出指定用户创建的所有任务
func ListTasksByUserID(userID uint) ([]Task, error) {
	var tasks []Task
	err := db.GetDB().Where("user_id = ?", userID).Order("created_at desc").Find(&tasks).Error
	return tasks, err
}

// ListAllTasks 列出系统中全部任务记录
func ListAllTasks() ([]Task, error) {
	var tasks []Task
	err := db.GetDB().Order("created_at desc").Find(&tasks).Error
	return tasks, err
}

// DeleteTaskByID 根据任务 ID 删除指定任务
func DeleteTaskByID(id uint) error {
	return db.GetDB().Delete(&Task{}, id).Error
}

// BatchDeleteTasks 批量删除任务（支持可选用户ID条件）
func BatchDeleteTasks(ids []uint, userID *uint) error {
	query := db.GetDB()
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	return query.Where("id IN ?", ids).Delete(&Task{}).Error
}

// UpdateTaskStatus 更新任务状态
func UpdateTaskStatus(taskID uint, status string) error {
	return db.GetDB().Model(&Task{}).Where("id = ?", taskID).Update("status", status).Error
}
