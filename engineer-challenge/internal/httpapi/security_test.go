package httpapi_test

import (
	"bytes"
	"net/http"
	"testing"
)

// M2: requests larger than the body cap are rejected (memory-exhaustion DoS guard).
// newTestServer and doJSON are defined in tenant_handlers_test.go (same test package).
func TestOversizedRequestBodyRejected(t *testing.T) {
	r := newTestServer(t)
	big := string(bytes.Repeat([]byte("a"), 2<<20)) // 2 MiB > 1 MiB cap
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"slug": "x", "name": big})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for oversized body, got %d", w.Code)
	}
}
