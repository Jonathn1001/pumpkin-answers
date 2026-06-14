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

// AllClaimTypes trả về thứ tự cố định để render/seed (slice mới mỗi lần, tránh mutation từ caller).
func AllClaimTypes() []ClaimType {
	return []ClaimType{Outpatient, Inpatient, Dental, Maternity, Optical}
}

type Claim struct {
	Type         ClaimType      `json:"type"`
	Amount       int64          `json:"amount"`
	SubmittedAt  time.Time      `json:"submittedAt"`
	CustomFields map[string]any `json:"customFields"`
}
