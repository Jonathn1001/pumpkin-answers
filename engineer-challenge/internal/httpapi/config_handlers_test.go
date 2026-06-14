package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/httpapi"
	"claimsplatform/internal/seed"
	"claimsplatform/internal/usecase"
)

func TestSaveInvalidDraftReturns422(t *testing.T) {
	svc := usecase.New(memory.New())
	_, _ = svc.CreateTenant(context.Background(), "Co", "")
	r := httpapi.NewRouter(svc)
	bad := seed.SafeGuard()
	bad.SLA.DefaultDays = 0 // invalid
	w := doJSON(r, http.MethodPost, "/api/tenants/co/versions", map[string]any{"config": bad, "note": "x"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d %s", w.Code, w.Body)
	}
}

func TestRollbackEndpointRestoresConfig(t *testing.T) {
	svc := usecase.New(memory.New())
	_, _ = svc.CreateTenant(context.Background(), "Co", "") // v1 default (tiered)
	r := httpapi.NewRouter(svc)

	w := doJSON(r, http.MethodPost, "/api/tenants/co/versions", map[string]any{"config": seed.GovHealth(), "note": "gov"})
	if w.Code != http.StatusCreated {
		t.Fatalf("draft: %d %s", w.Code, w.Body)
	}
	w = doJSON(r, http.MethodPost, "/api/tenants/co/versions/2/publish", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("publish: %d %s", w.Code, w.Body)
	}
	w = doJSON(r, http.MethodPost, "/api/tenants/co/rollback", map[string]any{"targetVersion": 1})
	if w.Code != http.StatusOK {
		t.Fatalf("rollback: %d %s", w.Code, w.Body)
	}
	w = doJSON(r, http.MethodGet, "/api/tenants/co/config", nil)
	var cfg domain.ConfigDocument
	if err := json.Unmarshal(w.Body.Bytes(), &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Approval.Model != "tiered" {
		t.Fatalf("expected tiered model after rollback to v1, got %q", cfg.Approval.Model)
	}
}
