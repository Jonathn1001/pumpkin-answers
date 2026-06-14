package domain

import "time"

type ClaimType string

const (
	Outpatient ClaimType = "OUTPATIENT"
	Inpatient  ClaimType = "INPATIENT"
	Dental     ClaimType = "DENTAL"
	Maternity  ClaimType = "MATERNITY"
	Optical    ClaimType = "OPTICAL"
)

// AllClaimTypes là thứ tự cố định để render/seed.
var AllClaimTypes = []ClaimType{Outpatient, Inpatient, Dental, Maternity, Optical}

type Claim struct {
	Type         ClaimType      `json:"type"`
	Amount       int64          `json:"amount"`
	SubmittedAt  time.Time      `json:"submittedAt"`
	CustomFields map[string]any `json:"customFields"`
}
