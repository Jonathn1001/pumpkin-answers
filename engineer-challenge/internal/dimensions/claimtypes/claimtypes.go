package claimtypes

import (
	"fmt"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

type dimension struct{}

func New() registry.Dimension { return dimension{} }

func init() { registry.Register(New()) }

func (dimension) Key() string { return "claimTypes" }

func (dimension) Evaluate(cfg domain.ConfigDocument, claim domain.Claim, dec *domain.ClaimDecision) {
	if !dec.Accepted {
		return
	}
	ct, ok := cfg.ClaimTypes[claim.Type]
	if !ok || !ct.Enabled {
		dec.Accepted = false
		dec.RejectionReasons = append(dec.RejectionReasons,
			fmt.Sprintf("claim type %s is not enabled for this tenant", claim.Type))
		dec.Trace = append(dec.Trace, domain.TraceEntry{
			Dimension: "claimTypes", Explanation: fmt.Sprintf("%s disabled -> rejected", claim.Type)})
		return
	}
	dec.RequiredDocuments = ct.RequiredDocuments
	dec.Trace = append(dec.Trace, domain.TraceEntry{
		Dimension: "claimTypes", Explanation: fmt.Sprintf("%s enabled; %d required documents", claim.Type, len(ct.RequiredDocuments))})
}

func (dimension) Validate(cfg domain.ConfigDocument) []domain.FieldError {
	var errs []domain.FieldError
	for ctype, c := range cfg.ClaimTypes {
		if !c.Enabled && len(c.RequiredDocuments) > 0 {
			errs = append(errs, domain.FieldError{
				Field:   fmt.Sprintf("claimTypes.%s.requiredDocuments", ctype),
				Message: "required documents set on a disabled claim type",
			})
		}
	}
	return errs
}

func (dimension) DefaultConfig() any {
	m := map[domain.ClaimType]domain.ClaimTypeConfig{}
	for _, t := range domain.AllClaimTypes() {
		m[t] = domain.ClaimTypeConfig{Enabled: false, RequiredDocuments: []string{}}
	}
	return m
}

func (dimension) UISchema() []registry.FieldDescriptor {
	return []registry.FieldDescriptor{{Key: "claimTypes", Label: "Claim Types & Documents", Type: "object", Widget: "claimtype-grid"}}
}
