export type ClaimType = 'OUTPATIENT' | 'INPATIENT' | 'DENTAL' | 'MATERNITY' | 'OPTICAL'

export interface FieldError { field: string; message: string }

export interface BrandingConfig { displayName: string; logoUrl?: string; primaryColor?: string; secondaryColor?: string; supportEmail?: string }
export interface ClaimTypeConfig { enabled: boolean; requiredDocuments: string[] }
export interface ApprovalTier { label: string; maxAmount: number | null; approverRole: string }
export interface Committee { name: string; requiredApprovals: number }
export interface ApprovalConfig { autoApproveThreshold: number; model: 'tiered' | 'committee'; tiers?: ApprovalTier[]; committee?: Committee | null }
export interface NotificationsConfig { channels: string[]; events: Record<string, string[]>; webhookUrl?: string }
export interface Escalation { warnBeforeDays: number; notifyRole: string }
export interface SLAConfig { defaultDays: number; perClaimType?: Record<string, number>; escalation: Escalation }
export interface FieldValidation { pattern?: string; min?: number; max?: number }
export interface CustomFieldConfig { key: string; label: string; type: string; required: boolean; options?: string[]; validation?: FieldValidation | null }

export interface ConfigDocument {
  branding: BrandingConfig
  claimTypes: Record<string, ClaimTypeConfig>
  approval: ApprovalConfig
  notifications: NotificationsConfig
  sla: SLAConfig
  customFields: CustomFieldConfig[]
}

export interface Tenant { slug: string; name: string; status: string; activeVersionNumber?: number; createdAt: string; updatedAt: string }
export interface ConfigVersion { tenantSlug: string; versionNumber: number; status: string; note?: string; createdBy?: string; config: ConfigDocument; createdAt: string }

export interface Claim { type: ClaimType; amount: number; submittedAt: string; customFields: Record<string, unknown> }

export interface ApprovalRoute { model: string; tierLabel?: string; approverRole?: string; committeeName?: string; requiredApprovals?: number }
export interface ApprovalDecision { outcome: 'auto_approved' | 'routed'; route?: ApprovalRoute }
export interface NotificationFire { event: string; channels: string[] }
export interface CustomFieldValidation { valid: boolean; errors: FieldError[] }
export interface TraceEntry { dimension: string; explanation: string }
export interface ClaimDecision {
  accepted: boolean
  rejectionReasons: string[]
  requiredDocuments?: string[]
  approval?: ApprovalDecision
  notifications: NotificationFire[]
  slaDeadline?: string
  slaDays?: number
  escalation?: Escalation
  customFieldValidation?: CustomFieldValidation
  trace: TraceEntry[]
}

export interface Change { path: string; type: 'added' | 'removed' | 'changed'; left?: unknown; right?: unknown }

export interface FieldDescriptor { key: string; label: string; type: string; widget: string; required?: boolean; options?: string[] }
export interface DimensionSchema { key: string; jsonSchema: Record<string, unknown>; ui: FieldDescriptor[] }
export interface ConfigSchemaResponse { dimensions: DimensionSchema[] }
