package httpapi_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/httpapi"
	"claimsplatform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	return httpapi.NewRouter(usecase.New(memory.New()))
}

func doJSON(r http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	var rd io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rd = bytes.NewReader(b)
	}
	req := httptest.NewRequest(method, path, rd)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateAndGetTenant(t *testing.T) {
	r := newTestServer(t)
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"slug": "acme", "name": "Acme"})
	if w.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", w.Code, w.Body)
	}
	w = doJSON(r, http.MethodGet, "/api/tenants/acme", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get: %d", w.Code)
	}
}

func TestGetUnknownTenantReturns404(t *testing.T) {
	w := doJSON(newTestServer(t), http.MethodGet, "/api/tenants/ghost", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateDuplicateSlugReturns409(t *testing.T) {
	r := newTestServer(t)
	_ = doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"slug": "dup", "name": "A"})
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"slug": "dup", "name": "B"})
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
}

func TestCreateTenantMissingFieldsReturns400(t *testing.T) {
	w := doJSON(newTestServer(t), http.MethodPost, "/api/tenants", map[string]string{"name": "NoSlug"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTenantWithInvalidStatusReturns422(t *testing.T) {
	r := newTestServer(t)
	_ = doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"slug": "co", "name": "Co"})
	w := doJSON(r, http.MethodPatch, "/api/tenants/co", map[string]string{"name": "X", "status": "banana"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d %s", w.Code, w.Body)
	}
}
