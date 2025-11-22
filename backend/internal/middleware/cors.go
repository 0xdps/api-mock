package middleware

import (
	"github.com/go-chi/cors"
)

// SetupCORS configures CORS for the API
func SetupCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-No-Cache", "Cache-Control"},
		ExposedHeaders:   []string{"Link", "X-Cache"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
