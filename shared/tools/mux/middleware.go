package mux

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
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
