// Package dimensions blank-imports every dimension so their init() Register() runs.
// Adding a new dimension = add one import line here (modular).
package dimensions

import (
	_ "claimsplatform/internal/dimensions/approval"
	_ "claimsplatform/internal/dimensions/branding"
	_ "claimsplatform/internal/dimensions/claimtypes"
	_ "claimsplatform/internal/dimensions/customfields"
	_ "claimsplatform/internal/dimensions/notifications"
	_ "claimsplatform/internal/dimensions/sla"
)
