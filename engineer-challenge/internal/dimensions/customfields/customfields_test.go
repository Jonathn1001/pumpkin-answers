package customfields_test

import (
	"testing"

	"claimsplatform/internal/dimensions/customfields"
	"claimsplatform/internal/domain"
)

func sptr(s string) *string { return &s }

func cfg() domain.ConfigDocument {
	return domain.ConfigDocument{CustomFields: []domain.CustomFieldConfig{
		{Key: "employeeId", Label: "Emp", Type: "string", Required: true, Validation: &domain.FieldValidation{Pattern: sptr(`^EMP\d{4}$`)}},
	}}
}

func TestMissingRequiredFails(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	customfields.New().Evaluate(cfg(), domain.Claim{CustomFields: map[string]any{}}, dec)
	if dec.CustomFieldValidation.Valid {
		t.Fatal("expected invalid for missing required field")
	}
}

func TestPatternMismatchFails(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	customfields.New().Evaluate(cfg(), domain.Claim{CustomFields: map[string]any{"employeeId": "X1"}}, dec)
	if dec.CustomFieldValidation.Valid {
		t.Fatal("expected invalid for pattern mismatch")
	}
}

func TestValidPasses(t *testing.T) {
	dec := &domain.ClaimDecision{Accepted: true}
	customfields.New().Evaluate(cfg(), domain.Claim{CustomFields: map[string]any{"employeeId": "EMP1234"}}, dec)
	if !dec.CustomFieldValidation.Valid {
		t.Fatalf("expected valid, got %+v", dec.CustomFieldValidation.Errors)
	}
}
