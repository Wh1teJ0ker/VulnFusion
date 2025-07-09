package db

import (
	"os"
	"testing"
	"time"

	"VulnFusion/internal/db"
	"VulnFusion/internal/models"
	"VulnFusion/internal/utils"

	"github.com/stretchr/testify/assert"
)

const testDBPath = "./testdata/test.db"

func cleanupTestDB() {
	_ = os.RemoveAll("./testdata")
}

func setupTestDB(t *testing.T) {
	dbConn, err := db.InitDatabase(testDBPath)
	assert.NoError(t, err)

	err = db.AutoRebuildModels(dbConn)
	assert.NoError(t, err)
}

func TestUserModel(t *testing.T) {
	defer cleanupTestDB()
	setupTestDB(t)

	hashed, err := utils.HashPassword("testpass")
	assert.NoError(t, err)

	user := &models.User{
		Username: "testuser",
		Password: hashed,
		Role:     "admin",
	}
	err = models.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	got, err := models.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, got.Username)

	match := utils.CheckPassword("testpass", got.Password)
	assert.True(t, match)

	users, err := models.ListAllUsers()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)

	userGot, err := models.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", userGot.Username)

	err = models.UpdateUserByID(user.ID, map[string]interface{}{"Role": "user"})
	assert.NoError(t, err)

	err = models.DeleteUserByID(user.ID)
	assert.NoError(t, err)
}

func TestTaskModel(t *testing.T) {
	defer cleanupTestDB()
	setupTestDB(t)

	hashed, _ := utils.HashPassword("pass")
	user := &models.User{Username: "taskuser", Password: hashed, Role: "user"}
	err := models.CreateUser(user)
	assert.NoError(t, err)

	task := &models.Task{
		UserID:    user.ID,
		Target:    "http://example.com",
		Template:  "cves/example.yaml",
		CreatedAt: time.Now(),
	}
	err = models.CreateTask(task)
	assert.NoError(t, err)
	assert.NotZero(t, task.ID)

	got, err := models.GetTaskByID(task.ID)
	assert.NoError(t, err)
	assert.Equal(t, task.Target, got.Target)

	allTasks, err := models.ListAllTasks()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(allTasks), 1)

	userTasks, err := models.ListTasksByUserID(user.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(userTasks), 1)

	err = models.DeleteTaskByID(task.ID)
	assert.NoError(t, err)
}

func TestResultModel(t *testing.T) {
	defer cleanupTestDB()
	setupTestDB(t)

	hashed, _ := utils.HashPassword("x")
	user := &models.User{Username: "resultuser", Password: hashed, Role: "user"}
	err := models.CreateUser(user)
	assert.NoError(t, err)

	task := &models.Task{
		UserID:    user.ID,
		Target:    "http://x",
		Template:  "x",
		CreatedAt: time.Now(),
	}
	err = models.CreateTask(task)
	assert.NoError(t, err)

	result := &models.Result{
		TaskID:        task.ID,
		Target:        "http://x",
		Vulnerability: "Example Vulnerability",
		Severity:      "high",
		Timestamp:     time.Now(),
	}
	err = models.SaveScanResult(result)
	assert.NoError(t, err)
	assert.NotZero(t, result.ID)

	results, err := models.ListResultsByTaskID(task.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)

	allResults, err := models.ListAllResults()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(allResults), 1)

	err = models.DeleteResultsByTaskID(task.ID)
	assert.NoError(t, err)
}
