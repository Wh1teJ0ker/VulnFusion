package storage

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"VulnFusion/internal/log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB 初始化 SQLite 数据库，包含连接检测和建表逻辑
func InitDB() {
	dbPath := "data/vulnfusion.db"

	// 创建数据库目录
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatal("无法创建数据库目录: %v", err)
	}

	// 建立数据库连接
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("数据库连接失败: %v", err)
	}

	// 设置数据库连接池参数
	DB.SetConnMaxLifetime(time.Minute * 10)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	// 检查连接是否正常
	if err := checkConnection(); err != nil {
		log.Fatal("数据库不可用: %v", err)
	}

	// 创建数据表结构
	createTables()

	log.Info("[DB] 数据库初始化完成: %s", dbPath)
}

// checkConnection 尝试执行一次 ping 和简单查询，确保数据库读写正常
func checkConnection() error {
	if err := DB.Ping(); err != nil {
		return err
	}

	row := DB.QueryRow("SELECT 1")
	var val int
	if err := row.Scan(&val); err != nil {
		return errors.New("数据库读取测试失败")
	}
	if val != 1 {
		return errors.New("数据库读取返回异常")
	}
	return nil
}
