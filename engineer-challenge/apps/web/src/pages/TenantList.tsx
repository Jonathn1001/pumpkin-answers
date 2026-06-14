import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTenants, useCreateTenant } from '../api/hooks'
import { ApiError } from '../api/client'
import { Badge } from '../components/Badge'

export function TenantList() {
  const { data: tenants, isLoading } = useTenants()
  const create = useCreateTenant()
  const [slug, setSlug] = useState('')
  const [name, setName] = useState('')
  const [cloneFrom, setCloneFrom] = useState('')
  const [err, setErr] = useState('')
  if (isLoading) return <div>Loading…</div>
  return (
    <div className="space-y-6">
      <table className="w-full rounded-lg border bg-white text-sm">
        <thead><tr className="border-b text-left"><th className="p-2">Name</th><th>Slug</th><th>Status</th><th>Active v</th></tr></thead>
        <tbody>{tenants?.map(t => (
          <tr key={t.slug} className="border-b hover:bg-gray-50">
            <td className="p-2"><Link className="text-blue-700" to={`/t/${t.slug}`}>{t.name}</Link></td>
            <td>{t.slug}</td><td><Badge>{t.status}</Badge></td><td>{t.activeVersionNumber ?? '—'}</td>
          </tr>
        ))}</tbody>
      </table>
      <form className="space-y-2 rounded-lg border bg-white p-4" onSubmit={e => {
        e.preventDefault(); setErr('')
        create.mutate({ slug, name, cloneFrom: cloneFrom || undefined }, {
          onSuccess: () => { setSlug(''); setName(''); setCloneFrom('') },
          onError: er => setErr(er instanceof ApiError ? (er.fields?.map(f => `${f.field}: ${f.message}`).join('; ') || er.message) : 'error'),
        })
      }}>
        <h3 className="font-semibold">Create tenant</h3>
        <div className="flex flex-wrap gap-2">
          <input className="rounded border px-2 py-1" placeholder="slug" value={slug} onChange={e => setSlug(e.target.value)} />
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
