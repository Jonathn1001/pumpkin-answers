import { useState } from 'react'
import type { WidgetProps } from './types'

export function FieldLabel({ p }: { p: WidgetProps }) {
  return (
    <div className="mb-1 text-sm font-medium">
      {p.descriptor.label}{p.descriptor.required && <span className="text-red-600"> *</span>}
      {p.errors.map((e, i) => <span key={i} className="ml-2 text-xs text-red-600">{e.message}</span>)}
    </div>
  )
}

export function TextInput(p: WidgetProps) {
  return (
    <label className="block">
      <FieldLabel p={p} />
      <input className="w-full rounded border px-2 py-1" value={(p.value as string) ?? ''} onChange={e => p.onChange(e.target.value)} />
    </label>
  )
}

export function NumberInput(p: WidgetProps) {
  return (
    <label className="block">
      <FieldLabel p={p} />
      <input type="number" className="w-full rounded border px-2 py-1" value={(p.value as number) ?? 0} onChange={e => p.onChange(Number(e.target.value))} />
    </label>
  )
}

export function Toggle(p: WidgetProps) {
  return (
    <label className="flex items-center gap-2">
      <input type="checkbox" checked={Boolean(p.value)} onChange={e => p.onChange(e.target.checked)} />
      <span className="text-sm font-medium">{p.descriptor.label}</span>
    </label>
  )
}

export function Select(p: WidgetProps) {
  return (
    <label className="block">
      <FieldLabel p={p} />
      <select className="w-full rounded border px-2 py-1" value={(p.value as string) ?? ''} onChange={e => p.onChange(e.target.value)}>
        <option value="" disabled>Select…</option>
        {(p.descriptor.options ?? []).map(o => <option key={o} value={o}>{o}</option>)}
      </select>
    </label>
  )
}

export function ColorInput(p: WidgetProps) {
  return (
    <label className="block">
      <FieldLabel p={p} />
      <input type="color" value={(p.value as string) || '#000000'} onChange={e => p.onChange(e.target.value)} />
    </label>
  )
}

export function LogoInput(p: WidgetProps) {
  const url = (p.value as string) ?? ''
  const [brokenUrl, setBrokenUrl] = useState('')
  const showImg = url !== '' && url !== brokenUrl
  return (
    <label className="block">
      <FieldLabel p={p} />
      <div className="flex items-center gap-3">
        {showImg ? (
          <img src={url} alt="logo preview" className="h-12 w-12 rounded border object-contain" onError={() => setBrokenUrl(url)} />
        ) : (
          <div className="flex h-12 w-12 items-center justify-center rounded border text-center text-[10px] text-gray-400">
            {url ? 'invalid' : 'no logo'}
          </div>
        )}
        <input className="flex-1 rounded border px-2 py-1" placeholder="https://…" value={url} onChange={e => p.onChange(e.target.value)} />
        {url && <button type="button" className="text-sm text-blue-700" onClick={() => p.onChange('')}>Remove</button>}
      </div>
    </label>
  )
}

export function FallbackWidget(p: WidgetProps) {
  return (
    <label className="block">
      <FieldLabel p={p} />
      <textarea key={JSON.stringify(p.value)} className="w-full rounded border px-2 py-1 font-mono text-xs" rows={4}
        defaultValue={JSON.stringify(p.value, null, 2)}
        onBlur={e => { try { p.onChange(JSON.parse(e.target.value)) } catch { /* ignore invalid json */ } }} />
    </label>
  )
}
