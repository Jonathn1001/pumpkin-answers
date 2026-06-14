import type { WidgetProps } from './types'
import type { ApprovalTier, ClaimTypeConfig, Committee, CustomFieldConfig, Escalation } from '../../api/types'
import { FieldLabel } from './base'

const CLAIM_TYPES = ['OUTPATIENT', 'INPATIENT', 'DENTAL', 'MATERNITY', 'OPTICAL']
const CHANNELS = ['email', 'sms', 'webhook']
const EVENTS = ['claim_submitted', 'claim_auto_approved', 'claim_routed', 'sla_breach_warning']

export function ClaimTypeGrid(p: WidgetProps) {
  const v = (p.value as Record<string, ClaimTypeConfig>) ?? {}
  const set = (t: string, c: ClaimTypeConfig) => p.onChange({ ...v, [t]: c })
  return (
    <div><FieldLabel p={p} />
      <table className="w-full text-sm"><tbody>
        {CLAIM_TYPES.map(t => {
          const c = v[t] ?? { enabled: false, requiredDocuments: [] }
          return (
            <tr key={t} className="border-b">
              <td className="py-1"><label className="flex items-center gap-2"><input type="checkbox" checked={c.enabled} onChange={e => set(t, { ...c, enabled: e.target.checked })} />{t}</label></td>
              <td className="py-1"><input className="w-full rounded border px-2 py-1" placeholder="docs, comma-separated" disabled={!c.enabled}
                value={c.requiredDocuments.join(', ')} onChange={e => set(t, { ...c, requiredDocuments: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })} /></td>
            </tr>
          )
        })}
      </tbody></table>
    </div>
  )
}

export function TierList(p: WidgetProps) {
  const tiers = (p.value as ApprovalTier[]) ?? []
  const upd = (i: number, t: ApprovalTier) => p.onChange(tiers.map((x, j) => (j === i ? t : x)))
  return (
    <div><FieldLabel p={p} />
      {tiers.map((t, i) => (
        <div key={i} className="mb-2 flex gap-2">
          <input className="rounded border px-2 py-1" placeholder="label" value={t.label} onChange={e => upd(i, { ...t, label: e.target.value })} />
          <input className="rounded border px-2 py-1" placeholder="maxAmount (blank = open)" value={t.maxAmount ?? ''} onChange={e => upd(i, { ...t, maxAmount: e.target.value === '' ? null : Number(e.target.value) })} />
          <input className="rounded border px-2 py-1" placeholder="approverRole" value={t.approverRole} onChange={e => upd(i, { ...t, approverRole: e.target.value })} />
          <button type="button" className="text-red-600" onClick={() => p.onChange(tiers.filter((_, j) => j !== i))}>✕</button>
        </div>
      ))}
      <button type="button" className="text-sm text-blue-700" onClick={() => p.onChange([...tiers, { label: '', maxAmount: null, approverRole: '' }])}>+ tier</button>
    </div>
  )
}

export function CommitteeForm(p: WidgetProps) {
  const c = (p.value as Committee | null) ?? { name: '', requiredApprovals: 1 }
  return (
    <div><FieldLabel p={p} />
      <div className="flex gap-2">
        <input className="rounded border px-2 py-1" placeholder="name" value={c.name} onChange={e => p.onChange({ ...c, name: e.target.value })} />
        <input type="number" className="rounded border px-2 py-1" placeholder="requiredApprovals" value={c.requiredApprovals} onChange={e => p.onChange({ ...c, requiredApprovals: Number(e.target.value) })} />
      </div>
    </div>
  )
}

export function ChannelMultiSelect(p: WidgetProps) {
  const sel = (p.value as string[]) ?? []
  const toggle = (ch: string) => p.onChange(sel.includes(ch) ? sel.filter(x => x !== ch) : [...sel, ch])
  return (
    <div><FieldLabel p={p} />
      <div className="flex gap-4">{CHANNELS.map(ch => <label key={ch} className="flex items-center gap-1"><input type="checkbox" checked={sel.includes(ch)} onChange={() => toggle(ch)} />{ch}</label>)}</div>
    </div>
  )
}

export function EventsGrid(p: WidgetProps) {
  const v = (p.value as Record<string, string[]>) ?? {}
  const toggle = (ev: string, ch: string) => {
    const cur = v[ev] ?? []
    p.onChange({ ...v, [ev]: cur.includes(ch) ? cur.filter(x => x !== ch) : [...cur, ch] })
  }
  return (
    <div><FieldLabel p={p} />
      <table className="text-sm"><thead><tr><th /> {CHANNELS.map(ch => <th key={ch} className="px-2">{ch}</th>)}</tr></thead>
        <tbody>{EVENTS.map(ev => (
          <tr key={ev}><td className="pr-2">{ev}</td>{CHANNELS.map(ch => <td key={ch} className="text-center"><input type="checkbox" checked={(v[ev] ?? []).includes(ch)} onChange={() => toggle(ev, ch)} /></td>)}</tr>
        ))}</tbody></table>
    </div>
  )
}

export function PerTypeNumberMap(p: WidgetProps) {
  const v = (p.value as Record<string, number>) ?? {}
  return (
    <div><FieldLabel p={p} />
      {CLAIM_TYPES.map(t => (
        <div key={t} className="mb-1 flex items-center gap-2 text-sm">
          <span className="w-28">{t}</span>
          <input type="number" className="rounded border px-2 py-1" placeholder="(default)" value={v[t] ?? ''} onChange={e => {
            const next = { ...v }; if (e.target.value === '') delete next[t]; else next[t] = Number(e.target.value); p.onChange(next)
          }} />
        </div>
      ))}
    </div>
  )
}

export function EscalationForm(p: WidgetProps) {
  const e = (p.value as Escalation) ?? { warnBeforeDays: 0, notifyRole: '' }
  return (
    <div><FieldLabel p={p} />
      <div className="flex gap-2">
        <input type="number" className="rounded border px-2 py-1" placeholder="warnBeforeDays" value={e.warnBeforeDays} onChange={ev => p.onChange({ ...e, warnBeforeDays: Number(ev.target.value) })} />
        <input className="rounded border px-2 py-1" placeholder="notifyRole" value={e.notifyRole} onChange={ev => p.onChange({ ...e, notifyRole: ev.target.value })} />
      </div>
    </div>
  )
}

export function CustomFieldsEditor(p: WidgetProps) {
  const fields = (p.value as CustomFieldConfig[]) ?? []
  const upd = (i: number, f: CustomFieldConfig) => p.onChange(fields.map((x, j) => (j === i ? f : x)))
  return (
    <div><FieldLabel p={p} />
      {fields.map((f, i) => (
        <div key={i} className="mb-2 flex flex-wrap gap-2">
          <input className="rounded border px-2 py-1" placeholder="key" value={f.key} onChange={e => upd(i, { ...f, key: e.target.value })} />
          <input className="rounded border px-2 py-1" placeholder="label" value={f.label} onChange={e => upd(i, { ...f, label: e.target.value })} />
          <select className="rounded border px-2 py-1" value={f.type} onChange={e => upd(i, { ...f, type: e.target.value })}>
            {['string', 'number', 'date', 'select', 'boolean'].map(t => <option key={t} value={t}>{t}</option>)}
          </select>
          <label className="flex items-center gap-1 text-sm"><input type="checkbox" checked={f.required} onChange={e => upd(i, { ...f, required: e.target.checked })} />required</label>
          <input className="rounded border px-2 py-1" placeholder="pattern (optional)" value={f.validation?.pattern ?? ''} onChange={e => upd(i, { ...f, validation: e.target.value ? { pattern: e.target.value } : null })} />
          <button type="button" className="text-red-600" onClick={() => p.onChange(fields.filter((_, j) => j !== i))}>✕</button>
        </div>
      ))}
      <button type="button" className="text-sm text-blue-700" onClick={() => p.onChange([...fields, { key: '', label: '', type: 'string', required: false }])}>+ field</button>
    </div>
  )
}
