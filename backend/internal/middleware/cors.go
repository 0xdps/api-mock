package middleware

import (
	"github.com/go-chi/cors"
)

// SetupCORS configures CORS for the API
func SetupCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			"Link",
			"X-Cache",
			"X-Request-ID",
			"X-Tenant-ID",
			"X-Role",
		},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
