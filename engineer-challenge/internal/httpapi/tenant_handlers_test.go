package httpapi_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/httpapi"
	"claimsplatform/internal/seed"
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
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"name": "Acme"})
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

func TestCreateSameNameTwiceAutoSuffixes(t *testing.T) {
	r := newTestServer(t)
	_ = doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"name": "Dup"})
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"name": "Dup"})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d %s", w.Code, w.Body)
	}
	var tn struct {
		Slug string `json:"slug"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &tn); err != nil {
		t.Fatal(err)
	}
	if tn.Slug != "dup-2" {
		t.Fatalf("second slug = %q, want %q", tn.Slug, "dup-2")
	}
}

func TestCreateTenantMissingNameReturns400(t *testing.T) {
	// name is the only field now; slug is derived from it server-side.
	w := doJSON(newTestServer(t), http.MethodPost, "/api/tenants", map[string]string{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateTenantNameOnlyDerivesSlug(t *testing.T) {
	w := doJSON(newTestServer(t), http.MethodPost, "/api/tenants", map[string]string{"name": "SafeGuard Insurance"})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d %s", w.Code, w.Body)
	}
	var tn struct {
		Slug string `json:"slug"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &tn); err != nil {
		t.Fatal(err)
	}
	if tn.Slug != "safeguard-insurance" {
		t.Fatalf("derived slug = %q, want %q", tn.Slug, "safeguard-insurance")
	}
}

func TestUpdateTenantWithInvalidStatusReturns422(t *testing.T) {
	r := newTestServer(t)
	_ = doJSON(r, http.MethodPost, "/api/tenants", map[string]string{"name": "Co"})
	w := doJSON(r, http.MethodPatch, "/api/tenants/co", map[string]string{"name": "X", "status": "banana"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d %s", w.Code, w.Body)
	}
}

// Create now accepts the full config from the client and persists it as v1.
func TestCreateTenantWithConfigPersistsIt(t *testing.T) {
	r := newTestServer(t)
	cfg := seed.GovHealth()
	cfg.Branding.DisplayName = "Custom Brand XYZ"
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]any{"name": "Custom Co", "config": cfg})
	if w.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", w.Code, w.Body)
	}
	w = doJSON(r, http.MethodGet, "/api/tenants/custom-co/config", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("get config: %d %s", w.Code, w.Body)
	}
	var got domain.ConfigDocument
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Branding.DisplayName != "Custom Brand XYZ" {
		t.Fatalf("client config not persisted: displayName = %q", got.Branding.DisplayName)
	}
}

// An invalid client-supplied config is rejected with field-level errors (422).
func TestCreateTenantWithInvalidConfigReturns422(t *testing.T) {
	r := newTestServer(t)
	cfg := seed.GovHealth()
	cfg.Branding.DisplayName = "" // branding.displayName is required
	w := doJSON(r, http.MethodPost, "/api/tenants", map[string]any{"name": "Bad Co", "config": cfg})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d %s", w.Code, w.Body)
	}
}

// The wizard seeds its "default" source from this endpoint.
func TestConfigDefaultReturnsValidDoc(t *testing.T) {
	w := doJSON(newTestServer(t), http.MethodGet, "/api/config-default", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d %s", w.Code, w.Body)
	}
	var got domain.ConfigDocument
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Branding.DisplayName == "" {
		t.Fatalf("expected a non-empty default branding displayName")
	}
}
