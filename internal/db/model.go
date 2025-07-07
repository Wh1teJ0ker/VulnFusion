package db

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:user"` // user / admin
}

type Task struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"index"`
	Target    string `gorm:"not null"`
	Template  string
	CreatedAt int64
}

type Result struct {
	ID            uint      `gorm:"primaryKey"`
	TaskID        uint      `gorm:"not null"`       // 所属任务 ID
	Target        string    `gorm:"not null"`       // 受影响目标
	Vulnerability string    `gorm:"not null"`       // 漏洞名称或标识
	Severity      string    `gorm:"default:medium"` // 风险等级：low / medium / high / critical
	Detail        string    `gorm:"type:text"`      // 详细信息（原始输出或解析后的内容）
	Timestamp     time.Time `gorm:"autoCreateTime"` // 记录时间
}
