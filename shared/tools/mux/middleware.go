package mux

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
	"github.com/go-chi/cors"
)

// CreateSessionProvider middleware
func CreateSessionProvider(f mongo.ContextFactory) func(http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		return sessionProvider(f, next)
	}
	return fn
}

func sessionProvider(f mongo.ContextFactory, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, err := f.Context(r.Context())
		if err != nil {
			RenderError(w, http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// CreateCorsEnabler middleware
func CreateCorsEnabler() func(http.Handler) http.Handler {
	mw := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	return mw.Handler
}
