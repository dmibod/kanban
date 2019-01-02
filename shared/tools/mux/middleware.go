package mux

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

// SessionMiddleWare provides mongo session
type SessionMiddleWare struct {
	sp mongo.SessionProvider
}

// SessionProvider puts session in context and calls next handler
func (m *SessionMiddleWare) SessionProvider(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := m.sp.WithSession(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
