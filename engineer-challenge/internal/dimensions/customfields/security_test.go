package customfields_test

import (
	"testing"

	"claimsplatform/internal/dimensions/customfields"
	"claimsplatform/internal/domain"
)

// L3: an admin-supplied custom-field regex must be rejected at config-save time,
// not silently fail at claim-evaluation time.
func TestInvalidRegexPatternFailsValidation(t *testing.T) {
	pat := "[" // invalid regular expression
	cfg := domain.ConfigDocument{CustomFields: []domain.CustomFieldConfig{
		{Key: "x", Type: "string", Validation: &domain.FieldValidation{Pattern: &pat}},
	}}
	found := false
	for _, e := range customfields.New().Validate(cfg) {
		if e.Field == "customFields[0].validation.pattern" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected an invalid-pattern validation error")
	}
}

func TestValidRegexPatternPasses(t *testing.T) {
	pat := `^EMP\d{4}$`
	cfg := domain.ConfigDocument{CustomFields: []domain.CustomFieldConfig{
		{Key: "x", Type: "string", Validation: &domain.FieldValidation{Pattern: &pat}},
	}}
	for _, e := range customfields.New().Validate(cfg) {
		if e.Field == "customFields[0].validation.pattern" {
			t.Fatalf("valid pattern must not error, got %v", e)
		}
	}
}
