package comparison_test

import (
	"testing"

	"claimsplatform/internal/comparison"
	"claimsplatform/internal/seed"
)

func TestIdenticalConfigsHaveNoChanges(t *testing.T) {
	changes, err := comparison.Diff(seed.SafeGuard(), seed.SafeGuard())
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %+v", changes)
	}
}

func TestChangedThresholdIsReportedAsChanged(t *testing.T) {
	right := seed.SafeGuard()
	right.Approval.AutoApproveThreshold = 99999
	changes, err := comparison.Diff(seed.SafeGuard(), right)
	if err != nil {
		t.Fatal(err)
	}
	sawChanged := false
	for _, c := range changes {
		if c.Type == comparison.Changed {
			sawChanged = true
		}
	}
	if !sawChanged {
		t.Fatalf("expected a 'changed' entry, got %+v", changes)
	}
}
