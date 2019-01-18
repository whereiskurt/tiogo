package middleware

import (
	"context"
	"net/http"
	"strings"
)

// Contexts extract all of the params related to their route
type contextMapKey string

func (c contextMapKey) String() string {
	return "pkg.server.context" + string(c)
}

// ContextMapKey is the key to the request context
var ContextMapKey = contextMapKey("ctxMap")

// ContextMap extract from request and type asserts it (helper function.)
func ContextMap(r *http.Request) map[string]string {
	return (r.Context().Value(ContextMapKey)).(map[string]string)
}

// InitialCtx runs for every route, sets the response to JSON for all responses and unpacks AccessKey&SecretKey
func InitialCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctxMap := make(map[string]string)

		xKeys := strings.Split(r.Header.Get("X-ApiKeys"), ";")
		for x := range xKeys {
			keys := strings.Split(xKeys[x], "=")
			switch {
			case strings.ToLower(keys[0]) == "accesskey":
				ctxMap["AccessKey"] = keys[1]

			case strings.ToLower(keys[0]) == "secretkey":
				ctxMap["SecretKey"] = keys[1]
			}
		}
		ctx := context.WithValue(r.Context(), ContextMapKey, ctxMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
