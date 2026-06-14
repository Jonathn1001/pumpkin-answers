package configschema_test

import (
	"testing"

	"claimsplatform/internal/configschema"
)

func TestSchemaCoversAllDimensions(t *testing.T) {
	resp := configschema.Get()
	keys := map[string]bool{}
	for _, d := range resp.Dimensions {
		keys[d.Key] = true
		if d.JSONSchema == nil {
			t.Fatalf("dimension %s missing JSON schema", d.Key)
		}
	}
	for _, want := range []string{"branding", "claimTypes", "approval", "notifications", "sla", "customFields"} {
		if !keys[want] {
			t.Fatalf("config-schema missing dimension %q", want)
		}
	}
}
