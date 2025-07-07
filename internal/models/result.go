// internal/storage/models/result.go
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

// SaveScanResult 保存扫描结果
func SaveScanResult(result *Result) error {
	return db.GetDB().Create(result).Error
}

// ListResultsByTask 获取某个任务的所有扫描结果
func ListResultsByTask(taskID int) ([]Result, error) {
	var results []Result
	if err := db.GetDB().Where("task_id = ?", taskID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// ListAllResults 获取所有扫描结果（管理员）
func ListAllResults() ([]Result, error) {
	var results []Result
	if err := db.GetDB().Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// DeleteResultsByTask 删除某任务下的所有结果
func DeleteResultsByTask(taskID int) error {
	return db.GetDB().Where("task_id = ?", taskID).Delete(&Result{}).Error
}
