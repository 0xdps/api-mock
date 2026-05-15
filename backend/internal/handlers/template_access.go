package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0xdps/api-mock/go/internal/auth"
	"github.com/0xdps/api-mock/go/internal/filters"
	"github.com/0xdps/api-mock/go/internal/middleware"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/0xdps/api-mock/go/internal/store"
	"github.com/go-chi/chi/v5"
)

// TemplateAccessHandler generates mock data from user-defined templates.
type TemplateAccessHandler struct {
	templates *store.Templates
	usage     *store.UsageLogs
	users     *store.Users
	registry  *schema.Registry
}

// NewTemplateAccessHandler creates a TemplateAccessHandler.
func NewTemplateAccessHandler(templates *store.Templates, usage *store.UsageLogs, users *store.Users, registry *schema.Registry) *TemplateAccessHandler {
	return &TemplateAccessHandler{templates: templates, usage: usage, users: users, registry: registry}
}

// GetTemplateData generates mock data for a user template.
// Route: GET /t/{templateId}
//
// Access rules:
//   - mak_ key → can only access public templates
//   - mk_  key → can access owner's own templates (any visibility)
func (h *TemplateAccessHandler) GetTemplateData(w http.ResponseWriter, r *http.Request) {
	apiKey := auth.GetAPIKey(r.Context())
	if apiKey == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
		return
	}

	templateID := chi.URLParam(r, "templateId")
	tmpl, err := h.templates.GetByID(templateID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}
	if tmpl == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}

	// Enforce access control
	if !canAccess(apiKey, tmpl) {
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	// Parse schema from stored JSON
	resourceSchema, err := parseTemplateSchema(tmpl)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid template schema"})
		return
	}

	// Determine how many items to generate
	ctx := r.Context()
	pagination, hasPagination := middleware.GetPagination(ctx)
	count := 10
	if hasPagination {
		count = pagination.Limit
	} else if v := r.URL.Query().Get("count"); v != "" {
		if n := parseInt(v); n > 0 && n <= 100 {
			count = n
		}
	}

	// Generate data using the parsed schema
	data, err := h.registry.GenerateFromSchema(resourceSchema, count)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "data generation failed"})
		return
	}

	// Apply search
	if search, hasSearch := middleware.GetSearch(ctx); hasSearch && search.Query != "" {
		data = middleware.ApplySearch(data, search)
	}

	// Apply filters
	if filterList := filters.ParseFilters(r.URL.Query()); len(filterList) > 0 {
		data = filters.ApplyFilters(data, filterList)
	}

	// Apply sorting
	if sorting, hasSorting := middleware.GetSorting(ctx); hasSorting && sorting.Field != "" {
		data = middleware.ApplySorting(data, sorting)
	}

	total := len(data)

	// Apply field projection
	if fields, hasFields := middleware.GetFields(ctx); hasFields {
		data = middleware.FilterFieldsSlice(data, fields)
	}

	// Apply pagination
	if hasPagination {
		start := pagination.Offset
		if start > total {
			start = total
		}
		end := start + pagination.Limit
		if end > total {
			end = total
		}
		data = data[start:end]
	}

	// Log usage asynchronously
	go h.usage.Append(apiKey.KeyID, tmpl.ID, http.StatusOK)

	respondJSON(w, http.StatusOK, map[string]any{
		"data": data,
		"pagination": map[string]any{
			"total": total,
			"limit": count,
			"offset": func() int {
				if hasPagination {
					return pagination.Offset
				}
				return 0
			}(),
		},
		"template": map[string]any{
			"id":         tmpl.ID,
			"name":       tmpl.Name,
			"slug":       tmpl.Slug,
			"visibility": tmpl.Visibility,
		},
	})
}

// GetTemplateMeta returns schema metadata for a user template.
// Route: GET /t/{templateId}/meta
func (h *TemplateAccessHandler) GetTemplateMeta(w http.ResponseWriter, r *http.Request) {
	apiKey := auth.GetAPIKey(r.Context())
	if apiKey == nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
		return
	}

	templateID := chi.URLParam(r, "templateId")
	tmpl, err := h.templates.GetByID(templateID)
	if err != nil || tmpl == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}

	if !canAccess(apiKey, tmpl) {
		respondJSON(w, http.StatusForbidden, map[string]string{"error": "access denied"})
		return
	}

	var rawSchema map[string]any
	if err := json.Unmarshal([]byte(tmpl.SchemaJSON), &rawSchema); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid template schema"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"template": map[string]any{
			"id":          tmpl.ID,
			"name":        tmpl.Name,
			"slug":        tmpl.Slug,
			"description": tmpl.Description,
			"visibility":  tmpl.Visibility,
			"type":        tmpl.Type,
			"base_schema": tmpl.BaseSchema,
			"created_at":  tmpl.CreatedAt,
			"updated_at":  tmpl.UpdatedAt,
		},
		"schema": rawSchema,
	})
}

// ── helpers ───────────────────────────────────────────────────────────────────

// canAccess enforces the two-key access model.
func canAccess(k *auth.APIKeyContext, tmpl *store.Template) bool {
	switch k.Type {
	case "access": // mak_ — read-only public templates from any user
		return tmpl.Visibility == "public"
	case "personal": // mk_ — owner's own templates only
		return tmpl.UserID == k.UserID
	}
	return false
}

// parseTemplateSchema unmarshals the stored schema JSON and returns a
// schema.Schema-compatible object that can generate data.
func parseTemplateSchema(tmpl *store.Template) (*schema.Schema, error) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(tmpl.SchemaJSON), &s); err != nil {
		return nil, err
	}
	// Ensure the resource metadata is set so GenerateItems works
	if s.Resource.Name == "" {
		s.Resource.Name = tmpl.Slug
	}
	return &s, nil
}

// GetBySlugPublic serves mock data for a featured template by slug.
// Route: GET /v1/{slug}
// No authentication required — featured templates are always public.
func (h *TemplateAccessHandler) GetBySlugPublic(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	tmpl, err := h.templates.GetBySlug(slug)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}
	if tmpl == nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "template not found"})
		return
	}

	resourceSchema, err := parseTemplateSchema(tmpl)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid template schema"})
		return
	}

	ctx := r.Context()
	pagination, hasPagination := middleware.GetPagination(ctx)
	count := 10
	if hasPagination {
		count = pagination.Limit
	} else if v := r.URL.Query().Get("count"); v != "" {
		if n := parseInt(v); n > 0 && n <= 100 {
			count = n
		}
	}

	data, err := h.registry.GenerateFromSchema(resourceSchema, count)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "data generation failed"})
		return
	}

	if search, hasSearch := middleware.GetSearch(ctx); hasSearch && search.Query != "" {
		data = middleware.ApplySearch(data, search)
	}

	if filterList := filters.ParseFilters(r.URL.Query()); len(filterList) > 0 {
		data = filters.ApplyFilters(data, filterList)
	}

	if sorting, hasSorting := middleware.GetSorting(ctx); hasSorting && sorting.Field != "" {
		data = middleware.ApplySorting(data, sorting)
	}

	total := len(data)

	if fields, hasFields := middleware.GetFields(ctx); hasFields {
		data = middleware.FilterFieldsSlice(data, fields)
	}

	if hasPagination {
		start := pagination.Offset
		if start > total {
			start = total
		}
		end := start + pagination.Limit
		if end > total {
			end = total
		}
		data = data[start:end]
	}

	// Log usage asynchronously (no API key for public access)
	go h.usage.Append("", tmpl.ID, http.StatusOK)

	respondJSON(w, http.StatusOK, map[string]any{
		"data": data,
		"pagination": map[string]any{
			"total": total,
			"limit": count,
			"offset": func() int {
				if hasPagination {
					return pagination.Offset
				}
				return 0
			}(),
		},
		"template": map[string]any{
			"id":          tmpl.ID,
			"name":        tmpl.Name,
			"slug":        tmpl.Slug,
			"visibility":  tmpl.Visibility,
			"is_featured": tmpl.IsFeatured,
		},
	})
}
