package mux

import (
	"net/http"

	"github.com/dmibod/kanban/shared/tools/db/mongo"
)

type MongoMiddleWare struct {
	mongo.SessionProvider
}

func (m *MongoMiddleWare) WithSession(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := m.SessionProvider.WithSession(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
