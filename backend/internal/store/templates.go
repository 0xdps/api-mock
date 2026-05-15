package store

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/mesahub-db/mesahub-pkg-go/mesahub"
)

var slugRe = regexp.MustCompile(`^[a-z0-9-]{1,60}$`)

// Templates handles template persistence.
type Templates struct {
	db *mesahub.DatabaseHandle
}

// NewTemplates returns a Templates store.
func NewTemplates(db *mesahub.DatabaseHandle) *Templates {
	return &Templates{db: db}
}

// ValidateSlug returns an error if the slug is not URL-safe.
func ValidateSlug(slug string) error {
	if !slugRe.MatchString(slug) {
		return fmt.Errorf("slug must be 1-60 lowercase alphanumeric characters or hyphens")
	}
	return nil
}

// Create inserts a new template. Slug uniqueness per user is enforced.
func (t *Templates) Create(tmpl *Template) (*Template, error) {
	if err := ValidateSlug(tmpl.Slug); err != nil {
		return nil, err
	}
	if len(tmpl.SchemaJSON) > 65536 {
		return nil, fmt.Errorf("schema_json exceeds 64 KB limit")
	}

	// Enforce slug uniqueness per user
	exists, err := t.db.Query(
		`SELECT id FROM templates WHERE user_id = ? AND slug = ? LIMIT 1`,
		[]any{tmpl.UserID, tmpl.Slug},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.Create slug check: %w", err)
	}
	if len(exists.Rows) > 0 {
		return nil, fmt.Errorf("a template with slug %q already exists", tmpl.Slug)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	visibility := tmpl.Visibility
	if visibility != "public" && visibility != "private" {
		visibility = "private"
	}
	ttype := tmpl.Type
	if ttype != "fork" && ttype != "custom" {
		ttype = "custom"
	}

	_, err = t.db.Exec(
		`INSERT INTO templates
		 (id, user_id, name, slug, description, visibility, type, base_schema, schema_json, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		[]any{id, tmpl.UserID, tmpl.Name, tmpl.Slug, tmpl.Description,
			visibility, ttype, nilIfEmpty(tmpl.BaseSchema), tmpl.SchemaJSON, now, now},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.Create insert: %w", err)
	}

	tmpl.ID = id
	tmpl.CreatedAt = now
	tmpl.UpdatedAt = now
	tmpl.Visibility = visibility
	tmpl.Type = ttype
	return tmpl, nil
}

// CreateFeatured inserts an official (featured) template, skipping if slug already exists globally.
func (t *Templates) CreateFeatured(tmpl *Template) (*Template, error) {
	exists, err := t.db.Query(
		`SELECT id FROM templates WHERE slug = ? LIMIT 1`,
		[]any{tmpl.Slug},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.CreateFeatured slug check: %w", err)
	}
	if len(exists.Rows) > 0 {
		return nil, nil // already seeded — skip
	}

	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.New().String()

	_, err = t.db.Exec(
		`INSERT INTO templates
		 (id, user_id, name, slug, description, visibility, type, base_schema, schema_json, is_featured, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 'public', 'custom', NULL, ?, 1, ?, ?)`,
		[]any{id, SystemUserID, tmpl.Name, tmpl.Slug, tmpl.Description, tmpl.SchemaJSON, now, now},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.CreateFeatured insert: %w", err)
	}

	tmpl.ID = id
	tmpl.UserID = SystemUserID
	tmpl.Visibility = "public"
	tmpl.Type = "custom"
	tmpl.IsFeatured = true
	tmpl.CreatedAt = now
	tmpl.UpdatedAt = now
	return tmpl, nil
}

// GetBySlug returns the featured template matching slug (for /v1/{slug} public access).
func (t *Templates) GetBySlug(slug string) (*Template, error) {
	rows, err := t.db.Query(
		`SELECT id, user_id, name, slug, description, visibility, type,
		        COALESCE(base_schema,'') as base_schema, schema_json, is_featured, created_at, updated_at
		 FROM templates WHERE slug = ? AND is_featured = 1 LIMIT 1`,
		[]any{slug},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.GetBySlug: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	tmpl := rowsToTemplates(rows.Rows)[0]
	return &tmpl, nil
}

// CountFeatured returns the number of featured (official) templates.
func (t *Templates) CountFeatured() (int, error) {
	rows, err := t.db.Query(
		`SELECT COUNT(*) as cnt FROM templates WHERE is_featured = 1`,
		nil,
	)
	if err != nil {
		return 0, err
	}
	if len(rows.Rows) == 0 {
		return 0, nil
	}
	cnt, _ := rows.Rows[0]["cnt"].(float64)
	return int(cnt), nil
}

// ListByUser returns all templates owned by a user.
func (t *Templates) ListByUser(userID string) ([]Template, error) {
	rows, err := t.db.Query(
		`SELECT id, user_id, name, slug, description, visibility, type,
		        COALESCE(base_schema,'') as base_schema, schema_json, is_featured, created_at, updated_at
		 FROM templates WHERE user_id = ? ORDER BY created_at DESC`,
		[]any{userID},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.ListByUser: %w", err)
	}
	return rowsToTemplates(rows.Rows), nil
}

// CountPublic returns the total number of public templates.
func (t *Templates) CountPublic() (int, error) {
	rows, err := t.db.Query(
		`SELECT COUNT(*) as total FROM templates WHERE visibility = 'public'`,
		[]any{},
	)
	if err != nil {
		return 0, fmt.Errorf("templates.CountPublic: %w", err)
	}
	if len(rows.Rows) == 0 {
		return 0, nil
	}
	switch v := rows.Rows[0]["total"].(type) {
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return 0, nil
	}
}

// ListPublic returns all public templates (for discovery).
func (t *Templates) ListPublic(limit, offset int) ([]Template, error) {
	rows, err := t.db.Query(
		`SELECT id, user_id, name, slug, description, visibility, type,
		        COALESCE(base_schema,'') as base_schema, schema_json, is_featured, created_at, updated_at
		 FROM templates WHERE visibility = 'public'
		 ORDER BY is_featured DESC, created_at DESC LIMIT ? OFFSET ?`,
		[]any{limit, offset},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.ListPublic: %w", err)
	}
	return rowsToTemplates(rows.Rows), nil
}

// GetByID returns a template by its UUID.
func (t *Templates) GetByID(id string) (*Template, error) {
	rows, err := t.db.Query(
		`SELECT id, user_id, name, slug, description, visibility, type,
		        COALESCE(base_schema,'') as base_schema, schema_json, is_featured, created_at, updated_at
		 FROM templates WHERE id = ? LIMIT 1`,
		[]any{id},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.GetByID: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	tmpl := rowsToTemplates(rows.Rows)[0]
	return &tmpl, nil
}

// GetPublicByID returns a template by ID only if it is publicly visible
// (visibility='public' or is_featured=1). Used by the unauthenticated detail endpoint.
func (t *Templates) GetPublicByID(id string) (*Template, error) {
	rows, err := t.db.Query(
		`SELECT id, user_id, name, slug, description, visibility, type,
		        COALESCE(base_schema,'') as base_schema, schema_json, is_featured, created_at, updated_at
		 FROM templates WHERE id = ? AND (visibility = 'public' OR is_featured = 1) LIMIT 1`,
		[]any{id},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.GetPublicByID: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	tmpl := rowsToTemplates(rows.Rows)[0]
	return &tmpl, nil
}

// GetByUserSlug returns a template by owner username-slug pair (for /t/{user}/{slug}).
func (t *Templates) GetByUserSlug(userID, slug string) (*Template, error) {
	rows, err := t.db.Query(
		`SELECT id, user_id, name, slug, description, visibility, type,
		        COALESCE(base_schema,'') as base_schema, schema_json, is_featured, created_at, updated_at
		 FROM templates WHERE user_id = ? AND slug = ? LIMIT 1`,
		[]any{userID, slug},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.GetByUserSlug: %w", err)
	}
	if len(rows.Rows) == 0 {
		return nil, nil
	}
	tmpl := rowsToTemplates(rows.Rows)[0]
	return &tmpl, nil
}

// Update replaces mutable fields of a template.
func (t *Templates) Update(id, userID string, patch map[string]any) (*Template, error) {
	if schemaJSON, ok := patch["schema_json"].(string); ok && len(schemaJSON) > 65536 {
		return nil, fmt.Errorf("schema_json exceeds 64 KB limit")
	}
	if slug, ok := patch["slug"].(string); ok {
		if err := ValidateSlug(slug); err != nil {
			return nil, err
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	_, err := t.db.Exec(
		`UPDATE templates SET
		   name        = COALESCE(?, name),
		   slug        = COALESCE(?, slug),
		   description = COALESCE(?, description),
		   visibility  = COALESCE(?, visibility),
		   schema_json = COALESCE(?, schema_json),
		   updated_at  = ?
		 WHERE id = ? AND user_id = ?`,
		[]any{
			patch["name"], patch["slug"], patch["description"],
			patch["visibility"], patch["schema_json"],
			now, id, userID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("templates.Update: %w", err)
	}
	return t.GetByID(id)
}

// Delete removes a template owned by userID. Featured templates cannot be deleted.
func (t *Templates) Delete(id, userID string) error {
	tmpl, err := t.GetByID(id)
	if err != nil {
		return err
	}
	if tmpl != nil && tmpl.IsFeatured {
		return fmt.Errorf("cannot delete a featured template")
	}
	_, err = t.db.Exec(
		`DELETE FROM templates WHERE id = ? AND user_id = ?`,
		[]any{id, userID},
	)
	return err
}

// CountByUser returns the number of non-featured templates owned by a user.
func (t *Templates) CountByUser(userID string) (int, error) {
	rows, err := t.db.Query(
		`SELECT COUNT(*) as cnt FROM templates WHERE user_id = ? AND is_featured = 0`,
		[]any{userID},
	)
	if err != nil {
		return 0, err
	}
	if len(rows.Rows) == 0 {
		return 0, nil
	}
	cnt, _ := rows.Rows[0]["cnt"].(float64)
	return int(cnt), nil
}

func rowsToTemplates(rows []map[string]any) []Template {
	out := make([]Template, 0, len(rows))
	for _, r := range rows {
		out = append(out, Template{
			ID:          str(r["id"]),
			UserID:      str(r["user_id"]),
			Name:        str(r["name"]),
			Slug:        str(r["slug"]),
			Description: str(r["description"]),
			Visibility:  str(r["visibility"]),
			Type:        str(r["type"]),
			BaseSchema:  str(r["base_schema"]),
			SchemaJSON:  str(r["schema_json"]),
			IsFeatured:  boolFromRow(r["is_featured"]),
			CreatedAt:   str(r["created_at"]),
			UpdatedAt:   str(r["updated_at"]),
		})
	}
	return out
}

func boolFromRow(v any) bool {
	switch vv := v.(type) {
	case float64:
		return vv != 0
	case int64:
		return vv != 0
	case bool:
		return vv
	}
	return false
}

func nilIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
