// Package engine: pure functions, no I/O. Fully driven by ConfigDocument.
package engine

import (
	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

func ProcessClaim(cfg domain.ConfigDocument, claim domain.Claim) domain.ClaimDecision {
	dec := &domain.ClaimDecision{
		Accepted:         true,
		RejectionReasons: []string{},
		Notifications:    []domain.NotificationFire{},
		Trace:            []domain.TraceEntry{},
	}

	// Gate: claimTypes runs first; if it rejects, stop.
	if gate := registry.Get(registry.KeyClaimTypes); gate != nil {
		gate.Evaluate(cfg, claim, dec)
	}
	if !dec.Accepted {
		return *dec
	}

	// Core dimensions next - independent of each other (approval, sla, customFields, branding).
	// notifications is excluded here on purpose (see below).
	for _, d := range registry.All() {
		if d.Key() == registry.KeyClaimTypes || d.Key() == registry.KeyNotifications {
			continue
		}
		d.Evaluate(cfg, claim, dec)
	}

	// notifications runs LAST: it reacts to the built decision (reads the approval
	// outcome to choose which events fire). Ordering is made explicit here so the
	// engine does not depend on dimension registration/import order.
	if n := registry.Get(registry.KeyNotifications); n != nil {
		n.Evaluate(cfg, claim, dec)
	}
	return *dec
}
