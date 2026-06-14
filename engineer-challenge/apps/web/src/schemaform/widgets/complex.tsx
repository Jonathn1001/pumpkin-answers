import { Button, Checkbox, Input, InputNumber, Select, Space, Table, Typography } from 'antd'
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import type { WidgetProps } from './types'
import type { ApprovalTier, ClaimTypeConfig, Committee, CustomFieldConfig, Escalation } from '../../api/types'
import { FieldLabel } from './base'
import { eventLabel } from '../../labels'

const CLAIM_TYPES = ['OUTPATIENT', 'INPATIENT', 'DENTAL', 'MATERNITY', 'OPTICAL']
const CHANNELS = ['email', 'sms', 'webhook']
const EVENTS = ['claim_submitted', 'claim_auto_approved', 'claim_routed', 'sla_breach_warning']

export function ClaimTypeGrid(p: WidgetProps) {
  const v = (p.value as Record<string, ClaimTypeConfig>) ?? {}
  const set = (t: string, c: ClaimTypeConfig) => p.onChange({ ...v, [t]: c })
  return (
    <div>
      <FieldLabel p={p} />
      <Space direction="vertical" style={{ width: '100%' }}>
        {CLAIM_TYPES.map((t) => {
          const c = v[t] ?? { enabled: false, requiredDocuments: [] }
          return (
            <Space key={t} style={{ width: '100%' }} align="center">
              <Checkbox checked={c.enabled} onChange={(e) => set(t, { ...c, enabled: e.target.checked })} style={{ width: 130 }}>
                {t}
              </Checkbox>
              <Input
                style={{ width: 320 }}
                placeholder="required docs, comma-separated"
                disabled={!c.enabled}
                value={c.requiredDocuments.join(', ')}
                onChange={(e) =>
                  set(t, { ...c, requiredDocuments: e.target.value.split(',').map((s) => s.trim()).filter(Boolean) })
                }
              />
            </Space>
          )
        })}
      </Space>
    </div>
  )
}

export function TierList(p: WidgetProps) {
  const tiers = (p.value as ApprovalTier[]) ?? []
  const upd = (i: number, t: ApprovalTier) => p.onChange(tiers.map((x, j) => (j === i ? t : x)))
  return (
    <div>
      <FieldLabel p={p} />
      <Space direction="vertical" style={{ width: '100%' }}>
        {tiers.map((t, i) => (
          <Space key={i} wrap>
            <Input placeholder="label" value={t.label} onChange={(e) => upd(i, { ...t, label: e.target.value })} style={{ width: 140 }} />
            <InputNumber
              placeholder="maxAmount (blank = open)"
              value={t.maxAmount}
              onChange={(v) => upd(i, { ...t, maxAmount: v ?? null })}
              style={{ width: 200 }}
            />
            <Input placeholder="approverRole" value={t.approverRole} onChange={(e) => upd(i, { ...t, approverRole: e.target.value })} style={{ width: 160 }} />
            <Button danger type="text" icon={<DeleteOutlined />} onClick={() => p.onChange(tiers.filter((_, j) => j !== i))} />
          </Space>
        ))}
        <Button
          type="dashed"
          icon={<PlusOutlined />}
          onClick={() => p.onChange([...tiers, { label: '', maxAmount: null, approverRole: '' }])}
        >
          Add tier
        </Button>
      </Space>
    </div>
  )
}

export function CommitteeForm(p: WidgetProps) {
  const c = (p.value as Committee | null) ?? { name: '', requiredApprovals: 1 }
  return (
    <div>
      <FieldLabel p={p} />
      <Space wrap>
        <Input placeholder="name" value={c.name} onChange={(e) => p.onChange({ ...c, name: e.target.value })} />
        <InputNumber
          placeholder="requiredApprovals"
          min={1}
          value={c.requiredApprovals}
          onChange={(v) => p.onChange({ ...c, requiredApprovals: v ?? 1 })}
        />
      </Space>
    </div>
  )
}

export function ChannelMultiSelect(p: WidgetProps) {
  const sel = (p.value as string[]) ?? []
  return (
    <div>
      <FieldLabel p={p} />
      <Checkbox.Group options={CHANNELS} value={sel} onChange={(vals) => p.onChange(vals)} />
    </div>
  )
}

export function EventsGrid(p: WidgetProps) {
  const v = (p.value as Record<string, string[]>) ?? {}
  const toggle = (ev: string, ch: string) => {
    const cur = v[ev] ?? []
    p.onChange({ ...v, [ev]: cur.includes(ch) ? cur.filter((x) => x !== ch) : [...cur, ch] })
  }
  return (
    <div>
      <FieldLabel p={p} />
      <Table
        size="small"
        pagination={false}
        dataSource={EVENTS.map((ev) => ({ key: ev, ev }))}
        columns={[
          { title: 'Event', dataIndex: 'ev', render: (ev: string) => eventLabel(ev) },
          ...CHANNELS.map((ch) => ({
            title: ch,
            key: ch,
            align: 'center' as const,
            render: (_: unknown, row: { ev: string }) => (
              <Checkbox checked={(v[row.ev] ?? []).includes(ch)} onChange={() => toggle(row.ev, ch)} />
            ),
          })),
        ]}
      />
    </div>
  )
}

export function PerTypeNumberMap(p: WidgetProps) {
  const v = (p.value as Record<string, number>) ?? {}
  return (
    <div>
      <FieldLabel p={p} />
      <Space direction="vertical">
        {CLAIM_TYPES.map((t) => (
          <Space key={t}>
            <Typography.Text style={{ width: 110, display: 'inline-block' }}>{t}</Typography.Text>
            <InputNumber
              placeholder="(default)"
              value={v[t]}
              onChange={(val) => {
                const next = { ...v }
                if (val == null) delete next[t]
                else next[t] = val
                p.onChange(next)
              }}
            />
          </Space>
        ))}
      </Space>
    </div>
  )
}

export function EscalationForm(p: WidgetProps) {
  const e = (p.value as Escalation) ?? { warnBeforeDays: 0, notifyRole: '' }
  return (
    <div>
      <FieldLabel p={p} />
      <Space wrap>
        <InputNumber
          placeholder="warnBeforeDays"
          value={e.warnBeforeDays}
          onChange={(v) => p.onChange({ ...e, warnBeforeDays: v ?? 0 })}
        />
        <Input placeholder="notifyRole" value={e.notifyRole} onChange={(ev) => p.onChange({ ...e, notifyRole: ev.target.value })} />
      </Space>
    </div>
  )
}

// Derive a stable camelCase field key from a human label:
// "Employee ID" -> "employeeId", "Policy Number" -> "policyNumber".
function deriveKey(label: string): string {
  const words = label.split(/[^a-zA-Z0-9]+/).filter(Boolean)
  return words
    .map((w, i) => (i === 0 ? w.toLowerCase() : w.charAt(0).toUpperCase() + w.slice(1).toLowerCase()))
    .join('')
}

export function CustomFieldsEditor(p: WidgetProps) {
  const fields = (p.value as CustomFieldConfig[]) ?? []
  const upd = (i: number, f: CustomFieldConfig) => p.onChange(fields.map((x, j) => (j === i ? f : x)))
  return (
    <div>
      <FieldLabel p={p} />
      <Space direction="vertical" style={{ width: '100%' }}>
        {fields.map((f, i) => (
          <Space key={i} wrap align="start">
            <div style={{ display: 'flex', flexDirection: 'column', width: 220 }}>
              <Input
                placeholder="label (e.g. Employee ID)"
                value={f.label}
                onChange={(e) => upd(i, { ...f, label: e.target.value, key: deriveKey(e.target.value) })}
              />
              <Typography.Text type="secondary" style={{ fontSize: 12, marginTop: 2 }}>
                key: {f.key || '—'}
              </Typography.Text>
            </div>
            <Select
              value={f.type}
              onChange={(val) => upd(i, { ...f, type: val })}
              style={{ width: 110 }}
              options={['string', 'number', 'date', 'select', 'boolean'].map((t) => ({ value: t, label: t }))}
            />
            <Checkbox checked={f.required} onChange={(e) => upd(i, { ...f, required: e.target.checked })}>
              required
            </Checkbox>
            <Input
              placeholder="pattern (optional)"
              value={f.validation?.pattern ?? ''}
              onChange={(e) => upd(i, { ...f, validation: e.target.value ? { pattern: e.target.value } : null })}
              style={{ width: 160 }}
            />
            <Button danger type="text" icon={<DeleteOutlined />} onClick={() => p.onChange(fields.filter((_, j) => j !== i))} />
          </Space>
        ))}
        <Button
          type="dashed"
          icon={<PlusOutlined />}
          onClick={() => p.onChange([...fields, { key: '', label: '', type: 'string', required: false }])}
        >
          Add field
        </Button>
      </Space>
    </div>
  )
}
