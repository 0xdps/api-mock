package main

import (
	"log"
	"os"
	"strings"

	"github.com/0xdps/api-mock/lib/handlers"
	"github.com/0xdps/api-mock/lib/middleware"
	"github.com/0xdps/api-mock/lib/schema"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load schemas from embedded files
	registry := schema.NewRegistry()
	if err := registry.LoadEmbeddedSchemas(); err != nil {
		log.Fatalf("Failed to load embedded schemas: %v", err)
	}

	log.Printf("Loaded %d schemas: %v", len(registry.Schemas), registry.GetAllResourceNames())

	r := gin.Default()

	// Middleware
	r.Use(middleware.SetupCORS())

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "apimock.codes API",
			"version":   "1.0.0",
			"docs":      "https://apimock.codes/docs",
			"resources": registry.GetAllResourceNames(),
		})
	})

	// Dynamic handler
	dynamicHandler := handlers.NewDynamicHandler(registry)

	// API routes - automatically generated from schemas with custom paths
	api := r.Group("/api")
	{
		// Sort resources: nested routes (with :params) first, then regular routes
		// This prevents Gin router conflicts
		resourceNames := registry.GetAllResourceNames()
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
			api.GET(routePath, dynamicHandler.GetCollection(resourceName))
			api.GET(routePath+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
			log.Printf("Registered nested route: %s -> %s", resourceName, routePath)
			
			// Register aliases
			for _, alias := range registry.GetRouteAliases(resourceName) {
				api.GET(alias, dynamicHandler.GetCollection(resourceName))
				api.GET(alias+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
				log.Printf("  + alias: %s", alias)
			}
		}
		
		// Then register regular routes (less specific)
		for _, resourceName := range regularResources {
			routePath := registry.GetRoutePath(resourceName)
			
			// Standard routes with /:id
			api.GET(routePath, dynamicHandler.GetCollection(resourceName))
			api.GET(routePath+"/:id", dynamicHandler.GetSingle(resourceName))
			api.GET(routePath+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
			log.Printf("Registered routes: %s -> %s", resourceName, routePath)
			
			// Register aliases
			for _, alias := range registry.GetRouteAliases(resourceName) {
				api.GET(alias, dynamicHandler.GetCollection(resourceName))
				api.GET(alias+"/:id", dynamicHandler.GetSingle(resourceName))
				api.GET(alias+"/meta", dynamicHandler.GetResourceMetadata(resourceName))
				log.Printf("  + alias: %s", alias)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
