package claimtypes_test

import (
	"testing"

	"claimsplatform/internal/dimensions/claimtypes"
	"claimsplatform/internal/domain"
)

func cfg() domain.ConfigDocument {
	return domain.ConfigDocument{
		ClaimTypes: map[domain.ClaimType]domain.ClaimTypeConfig{
			domain.Outpatient: {Enabled: true, RequiredDocuments: []string{"receipt", "prescription"}},
			domain.Dental:     {Enabled: false},
		},
	}
}

func TestEnabledTypeSetsRequiredDocs(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	claimtypes.New().Evaluate(cfg(), domain.Claim{Type: domain.Outpatient}, dec)
	if !dec.Accepted {
		t.Fatal("expected accepted")
	}
	if len(dec.RequiredDocuments) != 2 || dec.RequiredDocuments[0] != "receipt" {
		t.Fatalf("expected required docs, got %v", dec.RequiredDocuments)
	}
}

func TestDisabledTypeRejects(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	claimtypes.New().Evaluate(cfg(), domain.Claim{Type: domain.Dental}, dec)
	if dec.Accepted {
		t.Fatal("expected rejected for disabled type")
	}
	if len(dec.RejectionReasons) == 0 {
		t.Fatal("expected a rejection reason")
	}
}

func TestUnknownTypeRejects(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	claimtypes.New().Evaluate(cfg(), domain.Claim{Type: domain.Inpatient}, dec)
	if dec.Accepted {
		t.Fatal("expected rejected for a claim type absent from config")
	}
}

func TestValidateRejectsDocsOnDisabledType(t *testing.T) {
	c := domain.ConfigDocument{ClaimTypes: map[domain.ClaimType]domain.ClaimTypeConfig{
		domain.Dental: {Enabled: false, RequiredDocuments: []string{"receipt"}},
	}}
	if errs := claimtypes.New().Validate(c); len(errs) == 0 {
		t.Fatal("expected error: required documents on a disabled claim type")
	}
}

func TestZeroEnabledTypesFails(t *testing.T) {
	cfg := domain.ConfigDocument{ClaimTypes: map[domain.ClaimType]domain.ClaimTypeConfig{
		domain.Outpatient: {Enabled: false},
	}}
	errs := claimtypes.New().Validate(cfg)
	found := false
	for _, e := range errs {
		if e.Field == "claimTypes" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a claimTypes error when zero types enabled, got %v", errs)
	}
}

func TestDefaultConfigEnablesOutpatient(t *testing.T) {
	def, ok := claimtypes.New().DefaultConfig().(map[domain.ClaimType]domain.ClaimTypeConfig)
	if !ok {
		t.Fatalf("DefaultConfig wrong type: %T", claimtypes.New().DefaultConfig())
	}
	if !def[domain.Outpatient].Enabled {
		t.Fatal("expected OUTPATIENT enabled by default so a new tenant is processable")
	}
}
