package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/github/github-mcp-server/pkg/http/headers"
)

func TestRejectAuthorizationHeader(t *testing.T) {
	handler := RejectAuthorizationHeader()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	if rr.Header().Get(headers.ContentTypeHeader) != headers.ContentTypeJSON {
		t.Fatalf("expected JSON content type, got %q", rr.Header().Get(headers.ContentTypeHeader))
	}

	var resp jsonRPCErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("expected JSON-RPC error response: %v", err)
	}
	if resp.JSONRPC != "2.0" {
		t.Fatalf("expected JSON-RPC 2.0 response, got %q", resp.JSONRPC)
	}
	if resp.Error.Code != -32600 {
		t.Fatalf("expected invalid request error code, got %d", resp.Error.Code)
	}
	if resp.Error.Message != "Authorization header is not allowed in GitHub App auth mode" {
		t.Fatalf("unexpected error message: %q", resp.Error.Message)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr2.Code)
	}
}
