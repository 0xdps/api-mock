package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/0xdps/api-mock/go/internal/cache"
	"github.com/0xdps/api-mock/go/internal/handlers"
	"github.com/0xdps/api-mock/go/internal/middleware"
	"github.com/0xdps/api-mock/go/internal/redisstore"
	"github.com/0xdps/api-mock/go/internal/schema"
	"github.com/0xdps/api-mock/go/internal/static"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load schemas from embedded files
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		log.Fatalf("Failed to load embedded schemas: %v", err)
	}

	log.Printf("Loaded %d schemas: %v", len(registry.Schemas), registry.GetAllResourceNames())

	// Initialize Redis (optional - gracefully handle connection failures)
	redisConfig := redisstore.Config{
		Host:     getEnvString("REDIS_HOST", "localhost"),
		Port:     getEnvInt("REDIS_PORT", 6379),
		Password: getEnvString("REDIS_PASSWORD", ""),
		DB:       getEnvInt("REDIS_DB", 0),
	}

	redisStore, err := redisstore.NewStore(redisConfig)
	if err != nil {
		log.Printf("⚠️  Failed to connect to Redis: %v", err)
		log.Printf("⚠️  Falling back to local-only cache mode")
		redisStore = nil
	} else {
		defer redisStore.Close()
	}

	// Initialize cache with configuration from environment
	cacheModeStr := getEnvString("CACHE_MODE", "all") // off, local, remote, or all
	cacheMode := cache.CacheMode(cacheModeStr)
	
	// If Redis failed to connect, force local-only mode
	if redisStore == nil && (cacheMode == cache.CacheModeRemote || cacheMode == cache.CacheModeAll) {
		log.Printf("⚠️  Redis unavailable, switching from '%s' to 'local' cache mode", cacheMode)
		cacheMode = cache.CacheModeLocal
	}
	
	// Validate cache mode
	validModes := map[string]bool{
		"off":    true,
		"local":  true,
		"remote": true,
		"all":    true,
	}
	if !validModes[cacheModeStr] {
		log.Printf("⚠️  Invalid CACHE_MODE '%s', defaulting to 'all'", cacheModeStr)
		cacheMode = cache.CacheModeAll
	}
	
	log.Printf("📦 Cache mode: %s", cacheMode)
	
	cacheConfig := cache.Config{
		ItemsPerResource:    getEnvInt("CACHE_ITEMS_PER_RESOURCE", 100),
		Seed:                getEnvInt64("CACHE_SEED", 42), // Fixed seed for reproducibility
		MaxItemsPerResource: getEnvInt("MAX_ITEMS_PER_RESOURCE", 1000),
		Mode:                cacheMode,
	}

	apiCache := cache.NewCache(registry, cacheConfig, redisStore)

	// Check if Redis already has data (only if Redis is available)
	resourceNames := registry.GetAllResourceNames()
	hasData := false
	if redisStore != nil {
		hasData, err = redisStore.HasResources(resourceNames)
		if err != nil {
			log.Printf("⚠️  Failed to check Redis: %v", err)
			hasData = false
		}
	}

	if hasData {
		// Load from Redis
		log.Printf("✅ Redis has existing data, loading into memory...")
		if err := apiCache.LoadFromRedis(); err != nil {
			log.Fatalf("Failed to load from Redis: %v", err)
		}
	} else {
		// Generate and populate both Redis and in-memory
		log.Printf("📝 Redis is empty, generating and populating cache...")
		if err := apiCache.Warmup(); err != nil {
			log.Fatalf("Failed to warmup cache: %v", err)
		}
	}

	r := chi.NewRouter()

	// Core middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.SetupCORS().Handler)

	// Global parameter middleware
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.TenantMiddleware)
	r.Use(middleware.RBACMiddleware)
	r.Use(middleware.DelayMiddleware)
	r.Use(middleware.FlakyMiddleware)
	r.Use(middleware.CacheBypassMiddleware)
	r.Use(middleware.FieldFilterMiddleware)

	// Query parameter middleware (pagination, sorting, search)
	r.Use(middleware.PaginationMiddleware)
	r.Use(middleware.SortingMiddleware)
	r.Use(middleware.SearchMiddleware)

	// Idempotency middleware (with simple in-memory cache)
	idempotencyCache := middleware.NewSimpleIdempotencyCache()
	r.Use(middleware.IdempotencyMiddleware(idempotencyCache))

	log.Printf("🔌 Global middleware registered:")
	log.Printf("   ✓ Request ID, Tenant, RBAC")
	log.Printf("   ✓ Delay, Flaky, Cache Bypass, Field Filtering")
	log.Printf("   ✓ Pagination (page, limit, offset)")
	log.Printf("   ✓ Sorting (sort, order)")
	log.Printf("   ✓ Search (q, search)")
	log.Printf("   ✓ Idempotency (Idempotency-Key)")

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		groups := registry.GetAllGroups()
		response := map[string]interface{}{
			"message":   "Mockly API",
			"version":   "1.0.0",
			"docs":      "https://mockly.codes/docs",
			"resources": resourceNames,
			"groups":    groups,
		}
		json.NewEncoder(w).Encode(response)
	})

	// Favicon routes (serve embedded assets)
	r.Get("/favicon.svg", static.ServeIconSVG)
	r.Get("/favicon.ico", static.ServeIconICO)

	// Dynamic handler with cache
	dynamicHandler := handlers.NewDynamicHandler(registry, apiCache)

	// Utility handlers for testing
	utilityHandler := handlers.NewUtilityHandlers()

	// Testing utility routes (Phase 1)
	r.Route("/test", func(r chi.Router) {
		// Echo & Inspect
		r.Post("/echo", utilityHandler.Echo)
		r.Get("/echo", utilityHandler.Echo)
		
		// Delay & Timeout Simulation
		r.Get("/delay/{ms}", utilityHandler.Delay)
		r.Get("/delay-random", utilityHandler.DelayRandom)
		
		// Status & Errors
		r.Get("/status/{code}", utilityHandler.Status)
		r.Get("/error/validation", utilityHandler.ErrorValidation)
		
		// Chaos Testing
		r.Get("/flaky", utilityHandler.Flaky)
		r.Get("/chaos", utilityHandler.Chaos)
	})

	log.Printf("🧪 Testing utility routes:")
	log.Printf("   POST/GET /test/echo            - Echo request details")
	log.Printf("   GET      /test/delay/:ms       - Delay response by N milliseconds")
	log.Printf("   GET      /test/delay-random    - Random delay (min/max params)")
	log.Printf("   GET      /test/status/:code    - Return specific HTTP status")
	log.Printf("   GET      /test/error/validation - Return 422 validation error")
	log.Printf("   GET      /test/flaky            - Randomly succeed/fail")
	log.Printf("   GET      /test/chaos            - Random status/response shape")

	// Group routes (must come before dynamic resource routes to avoid conflicts)
	groups := registry.GetAllGroups()
	groupNames := make([]string, 0, len(groups))
	for groupName := range groups {
		groupNames = append(groupNames, groupName)
	}

	// Register group-specific routes first (more specific paths)
	for _, groupName := range groupNames {
		// Capture groupName in closure for handlers
		gn := groupName // Capture for closure

		// Group info endpoint: /{group} - returns metadata only
		r.Get("/"+gn, func(w http.ResponseWriter, req *http.Request) {
			resourceNames := registry.GetResourceNamesByGroup(gn)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"group":     gn,
				"resources": resourceNames,
				"count":     len(resourceNames),
			})
		})

		// Group resource endpoints: /{group}/{resource}
		r.Get("/"+gn+"/{resource}", func(w http.ResponseWriter, req *http.Request) {
			resourceName := chi.URLParam(req, "resource")

			// Verify the resource belongs to this group
			groupResources := registry.GetResourceNamesByGroup(gn)
			found := false
			for _, res := range groupResources {
				if res == resourceName {
					found = true
					break
				}
			}

			if !found {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Resource not found in this group",
				})
				return
			}

			// Get count parameter
			countStr := req.URL.Query().Get("count")
			count := 10
			if countStr != "" {
				if c, err := strconv.Atoi(countStr); err == nil && c > 0 {
					count = c
					if count > 100 {
						count = 100
					}
				}
			}

			// Generate data using the schema
			data, err := registry.GenerateData(resourceName, count)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		})

		r.Get("/"+gn+"/{resource}/{id}", func(w http.ResponseWriter, req *http.Request) {
			resourceName := chi.URLParam(req, "resource")
			id := chi.URLParam(req, "id")

			// Verify the resource belongs to this group
			groupResources := registry.GetResourceNamesByGroup(gn)
			found := false
			for _, res := range groupResources {
				if res == resourceName {
					found = true
					break
				}
			}

			if !found {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Resource not found in this group",
				})
				return
			}

			// Generate a single item
			data, err := registry.GenerateData(resourceName, 1)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			if len(data) == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Resource not found",
				})
				return
			}

			// Set the ID to the requested ID
			item := data[0]
			if idNum, err := strconv.Atoi(id); err == nil {
				item["id"] = idNum
			} else {
				item["id"] = id
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
		})

		log.Printf("Registered group routes: /%s, /%s/{resource}, /%s/{resource}/{id}", gn, gn, gn)
	}

	// API routes - automatically generated from schemas with custom paths
	r.Route("/", func(r chi.Router) {
		// Sort resources: nested routes (with :params) first, then regular routes
		// This prevents router conflicts
		nestedResources := []string{}
		regularResources := []string{}

		// Create a map of group names for quick lookup
		groupNameMap := make(map[string]bool)
		for groupName := range groups {
			groupNameMap[groupName] = true
		}

		for _, name := range resourceNames {
			// Skip resources that conflict with group names
			// These are accessible via /{group}/{resource} instead
			routePath := registry.GetRoutePath(name)
			if groupNameMap[strings.TrimPrefix(routePath, "/")] {
				log.Printf("Skipping direct route for %s (conflicts with group name, use /{group}/%s instead)", name, name)
				continue
			}

			if strings.Contains(routePath, ":") {
				nestedResources = append(nestedResources, name)
			} else {
				regularResources = append(regularResources, name)
			}
		}

		// Register nested routes first (more specific paths)
		for _, resourceName := range nestedResources {
			routePath := registry.GetRoutePath(resourceName)

			// For nested routes, register without adding /:id suffix
			r.Get(routePath, dynamicHandler.GetCollection(resourceName))
			r.Post(routePath, dynamicHandler.PostCollection(resourceName))
			r.Get(routePath+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
			log.Printf("Registered nested route: %s -> %s", resourceName, routePath)

			// Register aliases
			for _, alias := range registry.GetRouteAliases(resourceName) {
				r.Get(alias, dynamicHandler.GetCollection(resourceName))
				r.Post(alias, dynamicHandler.PostCollection(resourceName))
				r.Get(alias+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
				log.Printf("  + alias: %s", alias)
			}
		}

		// Then register regular routes (less specific)
		for _, resourceName := range regularResources {
			routePath := registry.GetRoutePath(resourceName)

			// Standard routes with /:id
			r.Get(routePath, dynamicHandler.GetCollection(resourceName))
			r.Post(routePath, dynamicHandler.PostCollection(resourceName))
			r.Get(routePath+"/{id}", dynamicHandler.GetSingle(resourceName))
			r.Put(routePath+"/{id}", dynamicHandler.PutSingle(resourceName))
			r.Delete(routePath+"/{id}", dynamicHandler.DeleteSingle(resourceName))
			r.Get(routePath+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
			log.Printf("Registered routes: %s -> %s", resourceName, routePath)

			// Register aliases
			for _, alias := range registry.GetRouteAliases(resourceName) {
				r.Get(alias, dynamicHandler.GetCollection(resourceName))
				r.Post(alias, dynamicHandler.PostCollection(resourceName))
				r.Get(alias+"/{id}", dynamicHandler.GetSingle(resourceName))
				r.Put(alias+"/{id}", dynamicHandler.PutSingle(resourceName))
				r.Delete(alias+"/{id}", dynamicHandler.DeleteSingle(resourceName))
				r.Get(alias+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
				log.Printf("  + alias: %s", alias)
			}
		}
	})

	// Admin routes for cache management
	r.Route("/admin", func(r chi.Router) {
		// Cache stats
		r.Get("/cache/stats", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Cache statistics",
				"cache":   apiCache.GetStats(),
			})
		})

		// Cache refresh
		r.Post("/cache/refresh", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("🔄 Cache refresh requested")
			if err := apiCache.Refresh(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Cache refreshed successfully",
				"cache":   apiCache.GetStats(),
			})
		})
	})

	log.Printf("📊 Admin routes:")
	log.Printf("   GET  /admin/cache/stats   - View cache statistics")
	log.Printf("   POST /admin/cache/refresh - Refresh cache")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

// getEnvString reads a string from environment variable with default
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Invalid %s value '%s', using default %d", key, value, defaultValue)
		return defaultValue
	}
	return intValue
}

// getEnvInt64 reads an int64 from environment variable with default
func getEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Warning: Invalid %s value '%s', using default %d", key, value, defaultValue)
		return defaultValue
	}
	return intValue
}
