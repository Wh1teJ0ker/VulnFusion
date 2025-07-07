// internal/db/db.go
package db

import (
	"VulnFusion/internal/log"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"path/filepath"
)

var dbInstance *gorm.DB

// SetDB 设置全局数据库实例
func SetDB(db *gorm.DB) {
	dbInstance = db
}

// GetDB 返回全局数据库实例
func GetDB() *gorm.DB {
	return dbInstance
}

// InitDatabase 初始化数据库连接并执行模型结构同步（重建机制）
func InitDatabase(dbPath string) (*gorm.DB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Error("无法解析数据库路径: %v", err)
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		log.Error("无法创建数据库目录: %v", err)
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Error("数据库连接失败: %v", err)
		return nil, err
	}

	if err := AutoRebuildModels(db); err != nil {
		log.Error("模型结构同步失败: %v", err)
		return nil, err
	}

	SetDB(db)
	log.Info("数据库初始化完成，路径：%s", absPath)
	return db, nil
}

// AutoRebuildModels 检查并重建模型表结构（字段缺失/冗余/类型不一致）
func AutoRebuildModels(db *gorm.DB) error {
	modelsToCheck := []interface{}{
		&User{},
		&Task{},
		&Result{},
	}

	for _, model := range modelsToCheck {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			log.Error("模型解析失败: %v", err)
			return err
		}

		tableName := stmt.Schema.Table
		tempTableName := tableName + "_backup"

		type ColumnInfo struct {
			Name string `gorm:"column:name"`
		}
		var oldColumns []ColumnInfo
		if err := db.Raw(fmt.Sprintf("PRAGMA table_info(`%s`)", tableName)).Scan(&oldColumns).Error; err != nil {
			log.Error("无法读取表字段信息: %v", err)
			return err
		}

		if len(oldColumns) == 0 {
			log.Info("表 %s 不存在，正在创建...", tableName)
			if err := db.Migrator().CreateTable(model); err != nil {
				log.Error("创建表 %s 失败: %v", tableName, err)
				return err
			}
			continue
		}

		oldColsMap := map[string]struct{}{}
		for _, col := range oldColumns {
			oldColsMap[col.Name] = struct{}{}
		}

		destFields := getFieldNames(stmt.Schema)
		finalFields := make([]string, 0)
		for _, f := range destFields {
			if _, ok := oldColsMap[f]; ok {
				finalFields = append(finalFields, f)
			}
		}

		log.Info("检测到模型 [%s] 字段变化，重建中...", tableName)

		if err := db.Migrator().RenameTable(model, tempTableName); err != nil {
			log.Error("重命名表 %s 失败: %v", tableName, err)
			return err
		}
		if err := db.Migrator().CreateTable(model); err != nil {
			log.Error("重建表 %s 失败: %v", tableName, err)
			return err
		}

		insertSQL := fmt.Sprintf(
			"INSERT INTO `%s` (%s) SELECT %s FROM `%s`",
			tableName,
			joinFields(finalFields),
			joinFields(finalFields),
			tempTableName,
		)
		if err := db.Exec(insertSQL).Error; err != nil {
			log.Error("数据迁移失败（%s → %s）: %v", tempTableName, tableName, err)
			return err
		}

		if err := db.Migrator().DropTable(tempTableName); err != nil {
			log.Error("删除旧表 %s 失败: %v", tempTableName, err)
			return err
		}

		log.Info("模型 [%s] 重建完成", tableName)
	}
	return nil
}

func getFieldNames(sch *schema.Schema) []string {
	fields := make([]string, 0)
	for _, f := range sch.Fields {
		if !f.IgnoreMigration {
			fields = append(fields, f.DBName)
		}
	}
	return fields
}

func joinFields(fields []string) string {
	return "`" + join(fields, "`, `") + "`"
}

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
