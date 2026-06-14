import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTenants, useCreateTenant, useUpdateTenant } from '../api/hooks'
import { ApiError } from '../api/client'
import { Badge } from '../components/Badge'
import type { Tenant } from '../api/types'

function TenantRow({ tenant: t, onClone }: { tenant: Tenant; onClone: (slug: string) => void }) {
  const update = useUpdateTenant(t.slug)
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState(t.name)
  const archived = t.status === 'archived'
  return (
    <tr className="border-b hover:bg-gray-50">
      <td className="p-2">
        {editing ? (
          <span className="flex items-center gap-1">
            <input className="rounded border px-1 py-0.5" value={name} onChange={e => setName(e.target.value)} />
            <button className="text-blue-700" disabled={update.isPending}
              onClick={() => update.mutate({ name, status: t.status }, { onSuccess: () => setEditing(false) })}>Save</button>
            <button className="text-gray-500" onClick={() => { setName(t.name); setEditing(false) }}>Cancel</button>
          </span>
        ) : (
          <Link className="text-blue-700" to={`/t/${t.slug}`}>{t.name}</Link>
        )}
      </td>
      <td>{t.slug}</td>
      <td><Badge>{t.status}</Badge></td>
      <td>{t.activeVersionNumber ?? '—'}</td>
      <td className="space-x-2 whitespace-nowrap p-2">
        {!editing && <button className="text-blue-700" onClick={() => setEditing(true)}>Edit</button>}
        <button className="text-blue-700" disabled={update.isPending}
          onClick={() => update.mutate({ name: t.name, status: archived ? 'active' : 'archived' })}>
          {archived ? 'Activate' : 'Archive'}
        </button>
        <button className="text-blue-700" onClick={() => onClone(t.slug)}>Clone</button>
      </td>
    </tr>
  )
}

export function TenantList() {
  const { data: tenants, isLoading } = useTenants()
  const create = useCreateTenant()
  const [name, setName] = useState('')
  const [cloneFrom, setCloneFrom] = useState('')
  const [err, setErr] = useState('')
  if (isLoading) return <div>Loading…</div>
  return (
    <div className="space-y-6">
      <table className="w-full rounded-lg border bg-white text-sm">
        <thead><tr className="border-b text-left"><th className="p-2">Name</th><th>Slug</th><th>Status</th><th>Active v</th><th className="p-2">Actions</th></tr></thead>
        <tbody>{tenants?.map(t => <TenantRow key={t.slug} tenant={t} onClone={setCloneFrom} />)}</tbody>
      </table>
      <form className="space-y-2 rounded-lg border bg-white p-4" onSubmit={e => {
        e.preventDefault(); setErr('')
        create.mutate({ name, cloneFrom: cloneFrom || undefined }, {
          onSuccess: () => { setName(''); setCloneFrom('') },
          onError: er => setErr(er instanceof ApiError ? (er.fields?.map(f => `${f.field}: ${f.message}`).join('; ') || er.message) : 'error'),
        })
      }}>
        <h3 className="font-semibold">Create tenant</h3>
        <p className="text-xs text-gray-500">The URL slug is generated from the name by the server.</p>
        <div className="flex flex-wrap gap-2">
          <input className="rounded border px-2 py-1" placeholder="name" value={name} onChange={e => setName(e.target.value)} />
          <select className="rounded border px-2 py-1" value={cloneFrom} onChange={e => setCloneFrom(e.target.value)}>
            <option value="">(default config)</option>{tenants?.map(t => <option key={t.slug} value={t.slug}>clone: {t.slug}</option>)}
          </select>
          <button className="rounded bg-blue-600 px-3 py-1 text-white disabled:opacity-50" type="submit" disabled={create.isPending}>Create</button>
        </div>
        {err && <div className="text-sm text-red-600">{err}</div>}
      </form>
    </div>
  )
}
