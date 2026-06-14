package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"claimsplatform/internal/configrepo/memory"
	"claimsplatform/internal/domain"
	"claimsplatform/internal/httpapi"
	"claimsplatform/internal/seed"
	"claimsplatform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func canonicalClaim() domain.Claim {
	return domain.Claim{
		Type: domain.Outpatient, Amount: 10000,
		SubmittedAt: time.Date(2026, 6, 14, 0, 0, 0, 0, time.UTC),
		CustomFields: map[string]any{
			"employeeId": "EMP1234", "policyNumber": "HF-12345678", "memberTier": "Gold",
			"nationalId": "123456789012", "citizenCategory": "General",
		},
	}
}

func newSeededServer(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	repo := memory.New()
	if err := seed.SeedAll(context.Background(), repo); err != nil {
		t.Fatal(err)
	}
	return httpapi.NewRouter(usecase.New(repo))
}

func TestProcessEndpointSafeGuardAutoApproves(t *testing.T) {
	w := doJSON(newSeededServer(t), http.MethodPost, "/api/tenants/safeguard/process", canonicalClaim())
	if w.Code != http.StatusOK {
		t.Fatalf("process: %d %s", w.Code, w.Body)
	}
	var dec domain.ClaimDecision
	if err := json.Unmarshal(w.Body.Bytes(), &dec); err != nil {
		t.Fatal(err)
	}
	if dec.Approval == nil || dec.Approval.Outcome != domain.AutoApproved {
		t.Fatalf("expected auto_approved, got %+v", dec.Approval)
	}
}

func TestPreviewEndpointWithInlineConfigRoutesToCommittee(t *testing.T) {
	body := map[string]any{"claim": canonicalClaim(), "config": seed.GovHealth()}
	w := doJSON(newSeededServer(t), http.MethodPost, "/api/tenants/safeguard/preview", body)
	if w.Code != http.StatusOK {
		t.Fatalf("preview: %d %s", w.Code, w.Body)
	}
	var dec domain.ClaimDecision
	if err := json.Unmarshal(w.Body.Bytes(), &dec); err != nil {
		t.Fatal(err)
	}
	if dec.Approval == nil || dec.Approval.Route == nil || dec.Approval.Route.CommitteeName == "" {
		t.Fatalf("inline GovHealth config should route to committee, got %+v", dec.Approval)
	}
}

func TestDiffEndpointReturnsChanges(t *testing.T) {
	w := doJSON(newSeededServer(t), http.MethodGet, "/api/diff?left=safeguard&right=govhealth", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("diff: %d %s", w.Code, w.Body)
	}
}

func TestConfigSchemaEndpoint(t *testing.T) {
	w := doJSON(newSeededServer(t), http.MethodGet, "/api/config-schema", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("config-schema: %d", w.Code)
	}
}
