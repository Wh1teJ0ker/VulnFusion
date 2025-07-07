// internal/storage/models/task.go
package models

import (
	"VulnFusion/internal/db"
	"time"
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
func GetTaskByID(id int) (*Task, error) {
	var task Task
	if err := db.GetDB().First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// ListTasksByUser 列出指定用户创建的所有任务
func ListTasksByUser(userID int) ([]Task, error) {
	var tasks []Task
	if err := db.GetDB().Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// ListAllTasks 列出系统中全部任务记录
func ListAllTasks() ([]Task, error) {
	var tasks []Task
	if err := db.GetDB().Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteTaskByID 根据任务 ID 删除指定任务
func DeleteTaskByID(id int) error {
	return db.GetDB().Delete(&Task{}, id).Error
}
