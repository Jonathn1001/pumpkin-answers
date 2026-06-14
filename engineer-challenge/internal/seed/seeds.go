package seed

import "claimsplatform/internal/domain"

func intp(v int64) *int64   { return &v }
func strp(s string) *string { return &s }

func SafeGuard() domain.ConfigDocument {
	return domain.ConfigDocument{
		Branding: domain.BrandingConfig{DisplayName: "SafeGuard Insurance", PrimaryColor: "#0A4D2C", SecondaryColor: "#082B19", SupportEmail: "support@safeguard.example"},
		ClaimTypes: map[domain.ClaimType]domain.ClaimTypeConfig{
			domain.Outpatient: {Enabled: true, RequiredDocuments: []string{"receipt", "prescription", "medical_report"}},
			domain.Inpatient:  {Enabled: true, RequiredDocuments: []string{"receipt", "discharge_summary", "itemized_bill"}},
			domain.Dental:     {Enabled: true, RequiredDocuments: []string{"receipt", "dental_chart"}},
			domain.Maternity:  {Enabled: false, RequiredDocuments: []string{}},
			domain.Optical:    {Enabled: false, RequiredDocuments: []string{}},
		},
		Approval: domain.ApprovalConfig{AutoApproveThreshold: 20000, Model: domain.ApprovalModelTiered, Tiers: []domain.ApprovalTier{
			{Label: "Manager", MaxAmount: intp(50000), ApproverRole: "claims_manager"},
			{Label: "Director", MaxAmount: intp(200000), ApproverRole: "claims_director"},
			{Label: "Board", MaxAmount: nil, ApproverRole: "board"},
		}},
		Notifications: domain.NotificationsConfig{Channels: []string{"email"}, Events: map[string][]string{
			"claim_submitted": {"email"}, "claim_auto_approved": {"email"}, "claim_routed": {"email"}, "sla_breach_warning": {"email"},
		}},
		SLA: domain.SLAConfig{DefaultDays: 7, PerClaimType: map[domain.ClaimType]int{domain.Outpatient: 5},
			Escalation: domain.Escalation{WarnBeforeDays: 2, NotifyRole: "claims_manager"}},
		CustomFields: []domain.CustomFieldConfig{
			{Key: "employeeId", Label: "Employee ID", Type: "string", Required: true, Validation: &domain.FieldValidation{Pattern: strp(`^EMP\d{4,6}$`)}},
			{Key: "department", Label: "Department", Type: "string", Required: false},
		},
	}
}

func HealthFirst() domain.ConfigDocument {
	return domain.ConfigDocument{
		Branding: domain.BrandingConfig{DisplayName: "HealthFirst", PrimaryColor: "#1D4ED8", SecondaryColor: "#1E3A8A", SupportEmail: "care@healthfirst.example"},
		ClaimTypes: map[domain.ClaimType]domain.ClaimTypeConfig{
			domain.Outpatient: {Enabled: true, RequiredDocuments: []string{"receipt", "prescription"}},
			domain.Inpatient:  {Enabled: true, RequiredDocuments: []string{"receipt", "discharge_summary"}},
			domain.Dental:     {Enabled: true, RequiredDocuments: []string{"receipt"}},
			domain.Maternity:  {Enabled: true, RequiredDocuments: []string{"receipt", "medical_report", "birth_certificate"}},
			domain.Optical:    {Enabled: true, RequiredDocuments: []string{"receipt", "optical_prescription"}},
		},
		Approval: domain.ApprovalConfig{AutoApproveThreshold: 5000, Model: domain.ApprovalModelTiered, Tiers: []domain.ApprovalTier{
			{Label: "Officer", MaxAmount: intp(50000), ApproverRole: "claims_officer"},
			{Label: "Manager", MaxAmount: nil, ApproverRole: "claims_manager"},
		}},
		Notifications: domain.NotificationsConfig{Channels: []string{"email", "sms"}, Events: map[string][]string{
			"claim_submitted": {"email", "sms"}, "claim_auto_approved": {"email", "sms"}, "claim_routed": {"email"}, "sla_breach_warning": {"email", "sms"},
		}},
		SLA: domain.SLAConfig{DefaultDays: 7, PerClaimType: map[domain.ClaimType]int{},
			Escalation: domain.Escalation{WarnBeforeDays: 2, NotifyRole: "claims_manager"}},
		CustomFields: []domain.CustomFieldConfig{
			{Key: "policyNumber", Label: "Policy Number", Type: "string", Required: true, Validation: &domain.FieldValidation{Pattern: strp(`^HF-\d{8}$`)}},
			{Key: "memberTier", Label: "Member Tier", Type: "select", Required: true, Options: []string{"Silver", "Gold", "Platinum"}},
		},
	}
}

func GovHealth() domain.ConfigDocument {
	return domain.ConfigDocument{
		Branding: domain.BrandingConfig{DisplayName: "GovHealth", PrimaryColor: "#7C2D12", SecondaryColor: "#431407", SupportEmail: "help@govhealth.example.gov"},
		ClaimTypes: map[domain.ClaimType]domain.ClaimTypeConfig{
			domain.Outpatient: {Enabled: true, RequiredDocuments: []string{"receipt", "national_id", "referral_letter"}},
			domain.Inpatient:  {Enabled: true, RequiredDocuments: []string{"receipt", "national_id", "discharge_summary", "admission_authorization"}},
			domain.Dental:     {Enabled: false, RequiredDocuments: []string{}},
			domain.Maternity:  {Enabled: false, RequiredDocuments: []string{}},
			domain.Optical:    {Enabled: false, RequiredDocuments: []string{}},
		},
		Approval: domain.ApprovalConfig{AutoApproveThreshold: 0, Model: domain.ApprovalModelCommittee,
			Committee: &domain.Committee{Name: "Government Claims Committee", RequiredApprovals: 3}},
		Notifications: domain.NotificationsConfig{Channels: []string{"email", "webhook"}, WebhookURL: "https://govhealth.example.gov/claim-events",
			Events: map[string][]string{
				"claim_submitted": {"email", "webhook"}, "claim_routed": {"email", "webhook"}, "sla_breach_warning": {"email", "webhook"},
			}},
		SLA: domain.SLAConfig{DefaultDays: 15, PerClaimType: map[domain.ClaimType]int{},
			Escalation: domain.Escalation{WarnBeforeDays: 3, NotifyRole: "committee_chair"}},
		CustomFields: []domain.CustomFieldConfig{
			{Key: "nationalId", Label: "National ID", Type: "string", Required: true, Validation: &domain.FieldValidation{Pattern: strp(`^\d{12}$`)}},
			{Key: "citizenCategory", Label: "Citizen Category", Type: "select", Required: true, Options: []string{"General", "LowIncome", "Veteran", "Disabled"}},
		},
	}
}
