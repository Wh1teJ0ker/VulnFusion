package models

import (
	"VulnFusion/internal/db"
	"time"
)

type Result struct {
	ID            uint      `gorm:"primaryKey"`
	TaskID        uint      `gorm:"not null"`       // 所属任务 ID
	Target        string    `gorm:"not null"`       // 受影响目标
	Vulnerability string    `gorm:"not null"`       // 漏洞名称或标识
	Severity      string    `gorm:"default:medium"` // 风险等级：low / medium / high / critical
	Detail        string    `gorm:"type:text"`      // 详细信息（原始输出或解析后的内容）
	Timestamp     time.Time `gorm:"autoCreateTime"` // 记录时间
}

// SaveScanResult 保存单条扫描结果
func SaveScanResult(result *Result) error {
	return db.GetDB().Create(result).Error
}

// ListResultsByTaskID 根据任务 ID 获取所有扫描结果
func ListResultsByTaskID(taskID uint) ([]Result, error) {
	var results []Result
	err := db.GetDB().Where("task_id = ?", taskID).Find(&results).Error
	return results, err
}

// ListAllResults 获取系统所有扫描结果（管理员）
func ListAllResults() ([]Result, error) {
	var results []Result
	err := db.GetDB().Find(&results).Error
	return results, err
}

// GetResultByID 根据结果 ID 获取详细信息
func GetResultByID(resultID uint) (*Result, error) {
	var result Result
	err := db.GetDB().First(&result, resultID).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteResultsByTaskID 删除指定任务的所有结果记录
func DeleteResultsByTaskID(taskID uint) error {
	return db.GetDB().Where("task_id = ?", taskID).Delete(&Result{}).Error
}
