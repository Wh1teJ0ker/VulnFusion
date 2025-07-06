package storage

import (
	"time"
)

func SaveScanResult(userID int, target string, template string, resultRaw string, status string) error {
	_, err := DB.Exec(`
		INSERT INTO scan_tasks (user_id, target, template, result, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, userID, target, template, resultRaw, status, time.Now())
	return err
}
