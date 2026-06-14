package seed_test

import (
	"testing"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/seed"
)

func TestSeedTenantsDistinct(t *testing.T) {
	s := seed.SafeGuard()
	h := seed.HealthFirst()
	g := seed.GovHealth()
	if s.Approval.AutoApproveThreshold != 20000 || h.Approval.AutoApproveThreshold != 5000 || g.Approval.AutoApproveThreshold != 0 {
		t.Fatal("auto-approve thresholds wrong")
	}
	if g.Approval.Model != domain.ApprovalModelCommittee {
		t.Fatal("GovHealth should use committee model")
	}
	if c, ok := s.ClaimTypes[domain.Maternity]; ok && c.Enabled {
		t.Fatal("SafeGuard should not enable MATERNITY")
	}
}
