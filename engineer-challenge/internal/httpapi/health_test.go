package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestHealthz checks the liveness probe used by the Railway healthcheck. It needs
// no service because the handler returns a static 200, so NewRouter(nil) is safe
// (svc is only dereferenced inside data handlers, never at route registration).
func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter(nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "ok") {
		t.Fatalf("body = %q, want it to report ok", w.Body.String())
	}
}
