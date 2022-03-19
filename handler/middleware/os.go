package middleware

import (
	"context"
	"net/http"

	ua "github.com/mileusna/useragent"
)

type ctxKey struct{}

var osKey = ctxKey{}

func OS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := ua.Parse(r.UserAgent())
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), osKey, ua.OS)))
	}
	return http.HandlerFunc(fn)
}
