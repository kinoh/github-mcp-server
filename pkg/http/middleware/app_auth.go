package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/github/github-mcp-server/pkg/http/headers"
)

type jsonRPCErrorResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      any          `json:"id"`
	Error   jsonRPCError `json:"error"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RejectAuthorizationHeader rejects requests with Authorization headers.
// Used in server-side GitHub App auth mode where user-provided tokens are disallowed.
func RejectAuthorizationHeader() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(headers.AuthorizationHeader) != "" {
				w.Header().Set(headers.ContentTypeHeader, headers.ContentTypeJSON)
				w.WriteHeader(http.StatusBadRequest)
				resp := jsonRPCErrorResponse{
					JSONRPC: "2.0",
					ID:      nil,
					Error: jsonRPCError{
						Code:    -32600,
						Message: "Authorization header is not allowed in GitHub App auth mode",
					},
				}
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					http.Error(w, "failed to write JSON-RPC error response", http.StatusInternalServerError)
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
