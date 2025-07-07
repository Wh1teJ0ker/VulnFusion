// internal/models/models/user.go
package models

import (
	"VulnFusion/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:user"` // user / admin
}

// CreateUser 创建新用户记录，写入用户名、密码哈希、角色等字段
func CreateUser(user *User) error {
	return db.GetDB().Create(user).Error
}

// GetUserByUsername 通过用户名查询单个用户信息，常用于登录验证或唯一性校验
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := db.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyPassword 验证用户输入的明文密码与数据库中的哈希值是否一致
func VerifyPassword(hashed string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

// ListAllUsers 返回数据库中所有用户的完整信息，仅管理员接口使用
func ListAllUsers() ([]User, error) {
	var users []User
	if err := db.GetDB().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID 根据用户 ID 查询详细信息，适用于管理员查看或编辑单个用户
func GetUserByID(id uint) (*User, error) {
	var user User
	if err := db.GetDB().First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserByID 根据用户 ID 更新用户信息（如角色、状态、密码等），用于管理员编辑
func UpdateUserByID(id uint, updates map[string]interface{}) error {
	return db.GetDB().Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteUserByID 根据用户 ID 删除指定用户，执行前应做好权限校验与数据保护
func DeleteUserByID(id uint) error {
	return db.GetDB().Delete(&User{}, id).Error
}
