import { useState } from 'react'
import { useActiveConfig, useConfigSchema, useSaveDraft, usePublish, useVersions } from '../../api/hooks'
import { ApiError } from '../../api/client'
import { SchemaForm } from '../../schemaform/SchemaForm'
import { useAjv } from '../../schemaform/useAjv'
import type { ConfigDocument, FieldError } from '../../api/types'

export function ConfigTab({ slug }: { slug: string }) {
  const schema = useConfigSchema()
  const active = useActiveConfig(slug)
  const versions = useVersions(slug)
  const save = useSaveDraft(slug)
  const publish = usePublish(slug)
  const validate = useAjv(schema.data?.dimensions ?? [])
  const [config, setConfig] = useState<ConfigDocument | null>(null)
  const [errors, setErrors] = useState<FieldError[]>([])
  const [msg, setMsg] = useState('')
  // Seed local edit buffer on first load (store-and-update-during-render pattern)
  if (config === null && active.data) setConfig(active.data)
  if (!schema.data || !config) return <div>Loading…</div>
  return (
    <div className="space-y-4">
      <SchemaForm schema={schema.data} config={config} errors={errors} onChange={setConfig} />
      <div className="flex items-center gap-3">
        <button className="rounded bg-blue-600 px-3 py-1 text-white" onClick={() => {
          setMsg('')
          const clientErrs = validate(config as unknown as Record<string, unknown>)
          if (clientErrs.length) { setErrors(clientErrs); return }
          save.mutate({ config, note: 'edited via admin' }, {
            onSuccess: v => { setErrors([]); versions.refetch(); publish.mutate(v.versionNumber, { onSuccess: () => setMsg(`Published v${v.versionNumber}`) }) },
            onError: e => setErrors(e instanceof ApiError && e.fields ? e.fields : [{ field: '', message: (e as Error).message }]),
          })
        }}>Save draft + Publish</button>
        {msg && <span className="text-sm text-green-700">{msg}</span>}
      </div>
      {errors.filter(e => !e.field).length > 0 && <div className="text-sm text-red-600">{errors.filter(e => !e.field).map(e => e.message).join('; ')}</div>}
    </div>
  )
}
