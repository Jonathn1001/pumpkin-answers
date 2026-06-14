package domain_test

import (
	"errors"
	"testing"

	"claimsplatform/internal/domain"
)

func TestValidationErrorImplementsError(t *testing.T) {
	var err error = domain.ValidationError{Fields: []domain.FieldError{{Field: "x", Message: "bad"}}}
	if err.Error() == "" {
		t.Fatal("ValidationError must produce a message")
	}
	var ve domain.ValidationError
	if !errors.As(err, &ve) || len(ve.Fields) != 1 {
		t.Fatalf("errors.As must recover the field list, got %+v", ve)
	}
}

func TestSentinelsAreDistinct(t *testing.T) {
	if errors.Is(domain.ErrTenantNotFound, domain.ErrVersionNotFound) {
		t.Fatal("sentinels must be distinct")
	}
}
