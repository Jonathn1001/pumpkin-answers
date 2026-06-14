package dimensions_test

import (
	"testing"

	_ "claimsplatform/internal/dimensions" // trigger init() of all six
	"claimsplatform/internal/registry"
)

func TestAllSixDimensionsRegistered(t *testing.T) {
	keys := map[string]bool{}
	for _, d := range registry.All() {
		keys[d.Key()] = true
	}
	for _, want := range []string{"branding", "claimTypes", "approval", "notifications", "sla", "customFields"} {
		if !keys[want] {
			t.Fatalf("dimension %q not registered", want)
		}
	}
}
