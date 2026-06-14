package registry_test

import (
	"testing"

	"claimsplatform/internal/domain"
	"claimsplatform/internal/registry"
)

type fakeDim struct{ key string }

func (f fakeDim) Key() string                                                               { return f.key }
func (f fakeDim) Evaluate(_ domain.ConfigDocument, _ domain.Claim, _ *domain.ClaimDecision) {}
func (f fakeDim) Validate(_ domain.ConfigDocument) []domain.FieldError                      { return nil }
func (f fakeDim) DefaultConfig() any                                                        { return nil }
func (f fakeDim) UISchema() []registry.FieldDescriptor                                      { return nil }

func TestRegisterAndAll(t *testing.T) {
	registry.Reset()
	registry.Register(fakeDim{key: "a"})
	registry.Register(fakeDim{key: "b"})
	got := registry.All()
	if len(got) != 2 || got[0].Key() != "a" || got[1].Key() != "b" {
		t.Fatalf("expected [a b] in order, got %v", got)
	}
}
