package db_test

import (
	"os"
	"testing"

	"VulnFusion/internal/db"
	"VulnFusion/internal/log"    // ✅ 添加日志模块初始化
	"VulnFusion/internal/models" // ✅ 正确的模型路径

	"github.com/stretchr/testify/assert"
)

const testDBPath = "./testdata/test.db"

func TestMain(m *testing.M) {
	log.InitLogger("dev", "debug") // ✅ 初始化日志系统（否则 logger 是 nil 会 panic）
	code := m.Run()
	os.Exit(code)
}

func cleanupTestDB() {
	_ = os.RemoveAll("./testdata")
}

func TestInitDatabase(t *testing.T) {
	defer cleanupTestDB()

	database, err := db.InitDatabase(testDBPath)
	assert.NoError(t, err)
	assert.NotNil(t, database)

	// 这行测试其实不太必要，因为 InitDatabase 本身已经执行了 AutoRebuildModels，
	// 你可以保留也可以删掉它
	err = database.AutoMigrate(&models.User{})
	assert.NoError(t, err)

	user := &models.User{
		Username: "testuser",
		Password: "testpass",
		Role:     "admin",
	}
	err = database.Create(user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestAutoRebuildModels(t *testing.T) {
	defer cleanupTestDB()

	database, err := db.InitDatabase(testDBPath)
	assert.NoError(t, err)

	err = db.AutoRebuildModels(database)
	assert.NoError(t, err)
}
