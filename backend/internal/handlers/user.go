package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/0xdps/api-mock/go/internal/auth"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/0xdps/api-mock/go/internal/store"
	"github.com/go-chi/chi/v5"
)

// UserHandlers implements all user-facing management routes.
type UserHandlers struct {
	users     *store.Users
	templates *store.Templates
	apikeys   *store.APIKeys
	registry  *schema.Registry
}

// NewUserHandlers creates a UserHandlers.
func NewUserHandlers(
	users *store.Users,
	templates *store.Templates,
	apikeys *store.APIKeys,
	registry *schema.Registry,
) *UserHandlers {
	return &UserHandlers{
		users:     users,
		templates: templates,
		apikeys:   apikeys,
		registry:  registry,
	}
}

// ── Auth sync ─────────────────────────────────────────────────────────────────

// SyncUser upserts the nube-auth user into MesaHub and auto-creates their
// access key on first login.
//
//	POST /auth/sync
func (h *UserHandlers) SyncUser(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	if nubeUser == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
		return
	}

	user, created, err := h.users.GetOrCreate(nubeUser.ID, nubeUser.Email, nubeUser.Name)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to sync user"})
		return
	}

	resp := map[string]any{"user": user, "created": created}

	// Auto-create access key for brand-new users
	if created {
		_, rawKey, err := h.apikeys.CreateAccessKey(nubeUser.ID)
		if err == nil {
			resp["access_key"] = rawKey // shown once
		}
	}

	respondJSON(w, http.StatusOK, resp)
}

// ── Me ────────────────────────────────────────────────────────────────────────

// GetMe returns the current user's profile.
//
//	GET /me
func (h *UserHandlers) GetMe(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	if nubeUser == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
		return
	}

	user, err := h.users.Get(nubeUser.ID)
	if err != nil || user == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found — call POST /auth/sync first"})
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// ── Templates ─────────────────────────────────────────────────────────────────

// ListTemplates returns all templates owned by the authenticated user.
//
//	GET /me/templates
func (h *UserHandlers) ListTemplates(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	tmpls, err := h.templates.ListByUser(nubeUser.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list templates"})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"templates": tmpls, "count": len(tmpls)})
}

// CreateTemplate creates a new template.
//
//	POST /me/templates
func (h *UserHandlers) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())

	// Enforce per-user template limit
	count, err := h.templates.CountByUser(nubeUser.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to count templates"})
		return
	}
	if count >= store.FreeTierMaxTemplates {
		respondJSON(w, http.StatusForbidden, map[string]any{
			"error": "template limit reached",
			"limit": store.FreeTierMaxTemplates,
		})
		return
	}

	var body struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
		Type        string `json:"type"`
		BaseSchema  string `json:"base_schema"`
		SchemaJSON  string `json:"schema_json"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	if strings.TrimSpace(body.Name) == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}
	if strings.TrimSpace(body.Slug) == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "slug is required"})
		return
	}
	if strings.TrimSpace(body.SchemaJSON) == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "schema_json is required"})
		return
	}

	// Validate that schema_json is valid JSON
	if !json.Valid([]byte(body.SchemaJSON)) {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "schema_json is not valid JSON"})
		return
	}

	// If forking a built-in, verify the base schema exists
	if body.Type == "fork" && body.BaseSchema != "" {
		if _, found := h.registry.Schemas[body.BaseSchema]; !found {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "base_schema not found in built-in schemas"})
			return
		}
	}

	tmpl, err := h.templates.Create(&store.Template{
		UserID:      nubeUser.ID,
		Name:        body.Name,
		Slug:        body.Slug,
		Description: body.Description,
		Visibility:  body.Visibility,
		Type:        body.Type,
		BaseSchema:  body.BaseSchema,
		SchemaJSON:  body.SchemaJSON,
	})
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, tmpl)
}

// GetTemplate returns a single template owned by the authenticated user.
//
//	GET /me/templates/{id}
func (h *UserHandlers) GetTemplate(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	id := chi.URLParam(r, "id")

	tmpl, err := h.templates.GetByID(id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}
	if tmpl == nil || tmpl.UserID != nubeUser.ID {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}
	respondJSON(w, http.StatusOK, tmpl)
}

// UpdateTemplate partially updates a template.
//
//	PUT /me/templates/{id}
func (h *UserHandlers) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	id := chi.URLParam(r, "id")

	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	// Validate schema_json if provided
	if raw, ok := patch["schema_json"].(string); ok {
		if !json.Valid([]byte(raw)) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "schema_json is not valid JSON"})
			return
		}
	}

	updated, err := h.templates.Update(id, nubeUser.ID, patch)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if updated == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}
	respondJSON(w, http.StatusOK, updated)
}

// DeleteTemplate deletes a template owned by the authenticated user.
//
//	DELETE /me/templates/{id}
func (h *UserHandlers) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	id := chi.URLParam(r, "id")

	// Confirm ownership before delete
	tmpl, err := h.templates.GetByID(id)
	if err != nil || tmpl == nil || tmpl.UserID != nubeUser.ID {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}

	if err := h.templates.Delete(id, nubeUser.ID); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete failed"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── API Keys ──────────────────────────────────────────────────────────────────

// ListAPIKeys returns all API keys for the authenticated user (hashes omitted).
//
//	GET /me/api-keys
func (h *UserHandlers) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	keys, err := h.apikeys.ListByUser(nubeUser.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list keys"})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"keys": keys, "count": len(keys)})
}

// CreateAPIKey generates a new personal API key.
//
//	POST /me/api-keys
func (h *UserHandlers) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Name) == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "name is required"})
		return
	}

	key, rawKey, err := h.apikeys.CreatePersonalKey(nubeUser.ID, body.Name)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create key"})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"key":     key,
		"raw_key": rawKey, // shown once only
	})
}

// RevokeAPIKey revokes a personal API key.
//
//	DELETE /me/api-keys/{id}
func (h *UserHandlers) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.apikeys.Revoke(id, nubeUser.ID); err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetAccessKeyInfo returns the current user access key metadata (no raw key).
//
//	GET /me/access-key
func (h *UserHandlers) GetAccessKeyInfo(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())

	k, err := h.apikeys.GetAccessKeyInfo(nubeUser.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}
	if k == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "no active access key — call POST /auth/sync"})
		return
	}
	respondJSON(w, http.StatusOK, k)
}

// RegenerateAccessKey revokes the current access key and issues a new one.
//
//	POST /me/access-key/regenerate
func (h *UserHandlers) RegenerateAccessKey(w http.ResponseWriter, r *http.Request) {
	nubeUser := auth.GetNubeUser(r.Context())

	key, rawKey, err := h.apikeys.CreateAccessKey(nubeUser.ID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to regenerate key"})
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"key":     key,
		"raw_key": rawKey, // shown once only
	})
}

// ── Public template discovery ─────────────────────────────────────────────────

// ListPublicTemplates returns paginated public templates from all users.
//
//	GET /templates
func (h *UserHandlers) ListPublicTemplates(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r, 20, 100)

	tmpls, err := h.templates.ListPublic(limit, offset)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list templates"})
		return
	}
	total, err := h.templates.CountPublic()
	if err != nil {
		total = 0
	}
	respondJSON(w, http.StatusOK, map[string]any{"templates": tmpls, "limit": limit, "offset": offset, "total": total})
}

// GetPublicTemplate returns a single public template's metadata (no auth required).
//
//	GET /templates/{id}
func (h *UserHandlers) GetPublicTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "missing id"})
		return
	}
	tmpl, err := h.templates.GetPublicByID(id)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}
	if tmpl == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"template": tmpl})
}

// ── helpers ───────────────────────────────────────────────────────────────────

func parseLimitOffset(r *http.Request, defaultLimit, maxLimit int) (int, int) {
	limit := defaultLimit
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n := parseInt(v); n > 0 {
			limit = n
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n := parseInt(v); n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}
