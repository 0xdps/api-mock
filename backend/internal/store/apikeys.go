package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mesahub-db/mesahub-pkg-go/mesahub"
)

const (
	prefixPersonal = "mk_"
	prefixAccess   = "mak_"
)

// APIKeys handles API key persistence.
type APIKeys struct {
	db *mesahub.DatabaseHandle
}

// NewAPIKeys returns an APIKeys store.
func NewAPIKeys(db *mesahub.DatabaseHandle) *APIKeys {
	return &APIKeys{db: db}
}

// GenerateKey creates a cryptographically random key with the given prefix.
// Returns the raw plaintext key (shown once) and its SHA-256 hash.
func GenerateKey(prefix string) (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate key: %w", err)
	}
	raw = prefix + hex.EncodeToString(b)
	h := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(h[:])
	return raw, hash, nil
}

// CreatePersonalKey generates a new personal API key (mk_) for a user.
// Returns the full raw key once — it is NOT stored.
func (a *APIKeys) CreatePersonalKey(userID, name string) (*APIKey, string, error) {
	raw, hash, err := GenerateKey(prefixPersonal)
	if err != nil {
		return nil, "", err
	}

	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	prefix := raw[:len(prefixPersonal)+8] // e.g. "mk_a1b2c3d4"

	_, err = a.db.Exec(
		`INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, type, status, created_at)
		 VALUES (?, ?, ?, ?, ?, 'personal', 'active', ?)`,
		[]any{id, userID, name, hash, prefix, now},
	)
	if err != nil {
		return nil, "", fmt.Errorf("apikeys.CreatePersonalKey: %w", err)
	}

	return &APIKey{
		ID:        id,
		UserID:    userID,
		Name:      name,
		KeyHash:   hash,
		KeyPrefix: prefix,
		Type:      "personal",
		Status:    "active",
		CreatedAt: now,
	}, raw, nil
}

// CreateAccessKey generates a new user access key (mak_) for a user.
// Any previous access key for this user is revoked first.
// Returns the full raw key once.
func (a *APIKeys) CreateAccessKey(userID string) (*APIKey, string, error) {
	// Revoke existing access keys for this user
	_, err := a.db.Exec(
		`UPDATE api_keys SET status = 'revoked' WHERE user_id = ? AND type = 'access' AND status = 'active'`,
		[]any{userID},
	)
	if err != nil {
		return nil, "", fmt.Errorf("apikeys.CreateAccessKey revoke old: %w", err)
	}

	raw, hash, err := GenerateKey(prefixAccess)
	if err != nil {
		return nil, "", err
	}

	id := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	prefix := raw[:len(prefixAccess)+8] // e.g. "mak_a1b2c3d4"

	_, err = a.db.Exec(
		`INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, type, status, created_at)
		 VALUES (?, ?, 'User Access Key', ?, ?, 'access', 'active', ?)`,
		[]any{id, userID, hash, prefix, now},
	)
	if err != nil {
		return nil, "", fmt.Errorf("apikeys.CreateAccessKey insert: %w", err)
	}

	return &APIKey{
		ID:        id,
		UserID:    userID,
		Name:      "User Access Key",
		KeyHash:   hash,
		KeyPrefix: prefix,
		Type:      "access",
		Status:    "active",
		CreatedAt: now,
	}, raw, nil
}

// GetByRawKey looks up an API key by its raw plaintext value.
// Returns nil, nil when not found or revoked.
func (a *APIKeys) GetByRawKey(raw string) (*APIKey, error) {
	h := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(h[:])

	rows, err := a.db.Query(
		`SELECT id, user_id, name, key_hash, key_prefix, type, status, created_at, last_used_at
		 FROM api_keys WHERE key_hash = ? AND status = 'active' LIMIT 1`,
		[]any{hash},
	)
	if err != nil {
		return nil, fmt.Errorf("apikeys.GetByRawKey: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	return rowToAPIKey(rows.Rows[0]), nil
}

// ListByUser returns all API keys for a user (hashes omitted in result).
func (a *APIKeys) ListByUser(userID string) ([]APIKey, error) {
	rows, err := a.db.Query(
		`SELECT id, user_id, name, key_hash, key_prefix, type, status, created_at, last_used_at
		 FROM api_keys WHERE user_id = ? ORDER BY created_at DESC`,
		[]any{userID},
	)
	if err != nil {
		return nil, fmt.Errorf("apikeys.ListByUser: %w", err)
	}
	out := make([]APIKey, 0, len(rows.Rows))
	for _, r := range rows.Rows {
		k := rowToAPIKey(r)
		k.KeyHash = "" // never expose hash
		out = append(out, *k)
	}
	return out, nil
}

// GetAccessKeyInfo returns the active access key metadata for a user (no hash).
func (a *APIKeys) GetAccessKeyInfo(userID string) (*APIKey, error) {
	rows, err := a.db.Query(
		`SELECT id, user_id, name, key_hash, key_prefix, type, status, created_at, last_used_at
		 FROM api_keys WHERE user_id = ? AND type = 'access' AND status = 'active' LIMIT 1`,
		[]any{userID},
	)
	if err != nil {
		return nil, fmt.Errorf("apikeys.GetAccessKeyInfo: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	k := rowToAPIKey(rows.Rows[0])
	k.KeyHash = ""
	return k, nil
}

// Revoke marks a key as revoked. Only the owning user can revoke their key.
func (a *APIKeys) Revoke(id, userID string) error {
	res, err := a.db.Exec(
		`UPDATE api_keys SET status = 'revoked' WHERE id = ? AND user_id = ? AND status = 'active'`,
		[]any{id, userID},
	)
	if err != nil {
		return fmt.Errorf("apikeys.Revoke: %w", err)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("key not found or already revoked")
	}
	return nil
}

// TouchLastUsed updates last_used_at asynchronously (called in a goroutine).
func (a *APIKeys) TouchLastUsed(id string) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, _ = a.db.Exec(
		`UPDATE api_keys SET last_used_at = ? WHERE id = ?`,
		[]any{now, id},
	)
}

func rowToAPIKey(r map[string]any) *APIKey {
	k := &APIKey{
		ID:        str(r["id"]),
		UserID:    str(r["user_id"]),
		Name:      str(r["name"]),
		KeyHash:   str(r["key_hash"]),
		KeyPrefix: str(r["key_prefix"]),
		Type:      str(r["type"]),
		Status:    str(r["status"]),
		CreatedAt: str(r["created_at"]),
	}
	if v, ok := r["last_used_at"].(string); ok && v != "" {
		k.LastUsedAt = &v
	}
	return k
}
