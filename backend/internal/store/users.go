package store

import (
	"fmt"
	"time"

	"github.com/mesahub-db/mesahub-pkg-go/mesahub"
)

// Users handles user persistence in MesaHub.
type Users struct {
	db *mesahub.DatabaseHandle
}

// NewUsers returns a Users store backed by the given database.
func NewUsers(db *mesahub.DatabaseHandle) *Users {
	return &Users{db: db}
}

// GetOrCreate finds the user by nube-auth ID or inserts a new row.
// Returns the User and whether it was newly created.
func (u *Users) GetOrCreate(id, email, name string) (*User, bool, error) {
	rows, err := u.db.Query(
		`SELECT id, email, name, plan, created_at FROM users WHERE id = ? LIMIT 1`,
		[]any{id},
	)
	if err != nil {
		return nil, false, fmt.Errorf("users.GetOrCreate query: %w", err)
	}

	if len(rows.Rows) > 0 {
		return rowToUser(rows.Rows[0]), false, nil
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_, err = u.db.Exec(
		`INSERT INTO users (id, email, name, plan, created_at) VALUES (?, ?, ?, 'free', ?)`,
		[]any{id, email, name, now},
	)
	if err != nil {
		return nil, false, fmt.Errorf("users.GetOrCreate insert: %w", err)
	}

	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Plan:      "free",
		CreatedAt: now,
	}, true, nil
}

// Get returns a user by ID. Returns nil, nil when not found.
func (u *Users) Get(id string) (*User, error) {
	rows, err := u.db.Query(
		`SELECT id, email, name, plan, created_at FROM users WHERE id = ? LIMIT 1`,
		[]any{id},
	)
	if err != nil {
		return nil, fmt.Errorf("users.Get: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	return rowToUser(rows.Rows[0]), nil
}

func rowToUser(row map[string]any) *User {
	return &User{
		ID:        str(row["id"]),
		Email:     str(row["email"]),
		Name:      str(row["name"]),
		Plan:      str(row["plan"]),
		CreatedAt: str(row["created_at"]),
	}
}

func str(v any) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
