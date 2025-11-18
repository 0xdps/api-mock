package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/0xdps/api-mock/go/internal/handlers"
	"github.com/0xdps/api-mock/go/internal/middleware"
	"github.com/0xdps/api-mock/go/internal/schema"
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

	r := chi.NewRouter()
	resourceNames := registry.GetAllResourceNames()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.SetupCORS().Handler)

	// Health check
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"message":   "Mockly API",
			"version":   "1.0.0",
			"docs":      "https://mockly.codes/docs",
			"resources": resourceNames,
		}
		json.NewEncoder(w).Encode(response)
	})

	// Dynamic handler
	dynamicHandler := handlers.NewDynamicHandler(registry)

	// API routes - automatically generated from schemas with custom paths
	r.Route("/", func(r chi.Router) {
		// Sort resources: nested routes (with :params) first, then regular routes
		// This prevents router conflicts
		nestedResources := []string{}
		regularResources := []string{}

		for _, name := range resourceNames {
			routePath := registry.GetRoutePath(name)
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
			r.Get(routePath+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
			log.Printf("Registered nested route: %s -> %s", resourceName, routePath)

			// Register aliases
			for _, alias := range registry.GetRouteAliases(resourceName) {
				r.Get(alias, dynamicHandler.GetCollection(resourceName))
				r.Get(alias+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
				log.Printf("  + alias: %s", alias)
			}
		}

		// Then register regular routes (less specific)
		for _, resourceName := range regularResources {
			routePath := registry.GetRoutePath(resourceName)

			// Standard routes with /:id
			r.Get(routePath, dynamicHandler.GetCollection(resourceName))
			r.Get(routePath+"/{id}", dynamicHandler.GetSingle(resourceName))
			r.Get(routePath+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
			log.Printf("Registered routes: %s -> %s", resourceName, routePath)

			// Register aliases
			for _, alias := range registry.GetRouteAliases(resourceName) {
				r.Get(alias, dynamicHandler.GetCollection(resourceName))
				r.Get(alias+"/{id}", dynamicHandler.GetSingle(resourceName))
				r.Get(alias+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
				log.Printf("  + alias: %s", alias)
			}
		}
	})

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
