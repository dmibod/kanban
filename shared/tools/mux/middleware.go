package mux

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// CreateSessionProvider middleware
func CreateSessionProvider(sp mongo.SessionProvider) func(http.Handler) http.Handler {
	fn := func(next http.Handler) http.Handler {
		return sessionProvider(sp, next)
	}
	return fn
}

func sessionProvider(sp mongo.SessionProvider, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := sp.WithSession(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
