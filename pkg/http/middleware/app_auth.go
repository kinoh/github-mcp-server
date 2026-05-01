package middleware

import (
	"net/http"

	"github.com/github/github-mcp-server/pkg/http/headers"
)

// RejectAuthorizationHeader rejects requests with Authorization headers.
// Used in server-side GitHub App auth mode where user-provided tokens are disallowed.
func RejectAuthorizationHeader() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(headers.AuthorizationHeader) != "" {
				http.Error(w, "Authorization header is not allowed in GitHub App auth mode", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
