package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"VulnFusion/internal/log"
	"VulnFusion/internal/storage/models"
)

var dbInstance *gorm.DB

// InitDatabase 初始化数据库并进行结构验证
func InitDatabase(dbPath string) (*gorm.DB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("无法解析数据库路径: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return nil, fmt.Errorf("无法创建数据库目录: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	if err := AutoRebuildModels(db); err != nil {
		return nil, fmt.Errorf("模型结构同步失败: %v", err)
	}

	dbInstance = db
	return db, nil
}

// GetDB 返回全局数据库实例
func GetDB() *gorm.DB {
	return dbInstance
}

// AutoRebuildModels 检查并重建模型表结构（字段缺失/冗余/类型不一致）
func AutoRebuildModels(db *gorm.DB) error {
	modelsToCheck := []interface{}{
		&models.User{},
		&models.Task{},
		&models.Result{},
	}

	for _, model := range modelsToCheck {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return fmt.Errorf("解析模型失败: %v", err)
		}

		tableName := stmt.Schema.Table
		tempTableName := tableName + "_backup"

		// 获取原始列信息
		var oldColumns []string
		db.Raw(fmt.Sprintf("PRAGMA table_info(`%s`)", tableName)).Scan(&oldColumns)
		if len(oldColumns) == 0 {
			// 表不存在，直接创建
			if err := db.Migrator().CreateTable(model); err != nil {
				return fmt.Errorf("创建表失败: %v", err)
			}
			continue
		}

		// 备份原数据
		if err := db.Migrator().RenameTable(model, tempTableName); err != nil {
			return fmt.Errorf("重命名表失败: %v", err)
		}
		if err := db.Migrator().CreateTable(model); err != nil {
			return fmt.Errorf("重建表失败: %v", err)
		}

		// 构造字段交集插入语句
		destFields := getFieldNames(stmt.Schema)
		insertSQL := fmt.Sprintf(
			"INSERT INTO `%s` (%s) SELECT %s FROM `%s`",
			tableName,
			joinFields(destFields),
			joinFields(destFields),
			tempTableName,
		)
		if err := db.Exec(insertSQL).Error; err != nil {
			return fmt.Errorf("数据迁移失败: %v", err)
		}

		// 删除临时表
		if err := db.Migrator().DropTable(tempTableName); err != nil {
			return fmt.Errorf("删除备份表失败: %v", err)
		}
	}
	return nil
}

// getFieldNames 获取模型定义的所有字段名称
func getFieldNames(sch *schema.Schema) []string {
	fields := make([]string, 0)
	for _, f := range sch.Fields {
		if !f.IgnoreMigration {
			fields = append(fields, f.DBName)
		}
	}
	return fields
}

// joinFields 拼接字段名为SQL字符串
func joinFields(fields []string) string {
	return "`" + join(fields, "`, `") + "`"
}

// join 手动拼接字符串（避免导入 strings）
func join(arr []string, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	out := arr[0]
	for i := 1; i < len(arr); i++ {
		out += sep + arr[i]
	}
	return out
}
