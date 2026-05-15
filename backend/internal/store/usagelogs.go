package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/mesahub-db/mesahub-pkg-go/mesahub"
)

// UsageLogs handles async write of usage records.
type UsageLogs struct {
	db *mesahub.DatabaseHandle
}

// NewUsageLogs returns a UsageLogs store.
func NewUsageLogs(db *mesahub.DatabaseHandle) *UsageLogs {
	return &UsageLogs{db: db}
}

// Append writes a usage record. Intended to be called in a goroutine.
func (u *UsageLogs) Append(apiKeyID, templateID string, statusCode int) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, _ = u.db.Exec(
		`INSERT INTO usage_logs (id, api_key_id, template_id, status_code, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		[]any{uuid.New().String(), apiKeyID, templateID, statusCode, now},
	)
}
