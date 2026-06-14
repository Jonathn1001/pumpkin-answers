package comparison_test

import (
	"testing"

	"claimsplatform/internal/comparison"
	"claimsplatform/internal/seed"
)

func TestDiffPathUsesJSONFieldNames(t *testing.T) {
	right := seed.SafeGuard()
	right.Branding.DisplayName = "Changed Name"
	changes, err := comparison.Diff(seed.SafeGuard(), right)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range changes {
		if c.Path == "branding.displayName" {
			return // found the expected camelCase JSON path
		}
	}
	t.Fatalf("expected a change with path 'branding.displayName', got %+v", changes)
}

func TestDiffDetectsSliceReordering(t *testing.T) {
	left := seed.SafeGuard() // 3 tiers: Manager, Director, Board
	right := seed.SafeGuard()
	// Swap tiers[0] and tiers[1]; same elements, different order.
	right.Approval.Tiers[0], right.Approval.Tiers[1] = right.Approval.Tiers[1], right.Approval.Tiers[0]
	changes, err := comparison.Diff(left, right)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("reordering approval tiers must produce at least one change")
	}
}

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
