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
	if gate := registry.Get("claimTypes"); gate != nil {
		gate.Evaluate(cfg, claim, dec)
	}
	if !dec.Accepted {
		return *dec
	}

	// Remaining dimensions (independent fields, order-agnostic).
	for _, d := range registry.All() {
		if d.Key() == "claimTypes" {
			continue
		}
		d.Evaluate(cfg, claim, dec)
	}
	return *dec
}
