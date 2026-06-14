import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useTenants, useUpdateTenant } from '../api/hooks'
import { Badge } from '../components/Badge'
import { CreateTenantWizard } from '../components/CreateTenantWizard'
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
  const [wizardOpen, setWizardOpen] = useState(false)
  const [initialClone, setInitialClone] = useState('')
  if (isLoading) return <div>Loading…</div>
  function openCreate(clone: string) { setInitialClone(clone); setWizardOpen(true) }
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Tenants</h2>
        <button className="rounded bg-blue-600 px-3 py-1.5 text-white" onClick={() => openCreate('')}>+ Create Tenant</button>
      </div>
      <table className="w-full rounded-lg border bg-white text-sm">
        <thead><tr className="border-b text-left"><th className="p-2">Name</th><th>Slug</th><th>Status</th><th>Active v</th><th className="p-2">Actions</th></tr></thead>
        <tbody>{tenants?.map(t => <TenantRow key={t.slug} tenant={t} onClone={openCreate} />)}</tbody>
      </table>
      {wizardOpen && (
        <CreateTenantWizard
          tenants={tenants ?? []}
          initialClone={initialClone}
          onClose={() => setWizardOpen(false)}
        />
      )}
    </div>
  )
}
