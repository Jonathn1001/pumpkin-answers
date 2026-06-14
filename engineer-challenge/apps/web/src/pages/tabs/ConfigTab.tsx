import { useState } from 'react'
import { Alert, App, Button, Space } from 'antd'
import { useActiveConfig, useConfigSchema, useSaveDraft, usePublish } from '../../api/hooks'
import { ApiError } from '../../api/client'
import { SchemaForm } from '../../schemaform/SchemaForm'
import { useAjv } from '../../schemaform/useAjv'
import { webhookRequired } from '../../schemaform/conditional'
import type { ConfigDocument, FieldError } from '../../api/types'

export function ConfigTab({ slug }: { slug: string }) {
  const { message } = App.useApp()
  const schema = useConfigSchema()
  const active = useActiveConfig(slug)
  const save = useSaveDraft(slug)
  const publish = usePublish(slug)
  const validate = useAjv(schema.data?.dimensions ?? [])
  const [config, setConfig] = useState<ConfigDocument | null>(null)
  const [errors, setErrors] = useState<FieldError[]>([])
  // Seed local edit buffer on first load (adjust-state-during-render pattern).
  if (config === null && active.data) setConfig(active.data)
  if (!schema.data || !config) return <div>Loading…</div>

  function saveAndPublish() {
    if (!config) return
    const clientErrs = validate(config as unknown as Record<string, unknown>)
    if (webhookRequired(config) && !config.notifications?.webhookUrl) {
      clientErrs.push({
        field: 'notifications.webhookUrl',
        message: 'Webhook URL is required when the webhook channel is enabled',
      })
    }
    if (clientErrs.length) {
      setErrors(clientErrs)
      return
    }
    save.mutate(
      { config, note: 'edited via admin' },
      {
        onSuccess: (v) => {
          setErrors([])
          publish.mutate(v.versionNumber, {
            onSuccess: () => message.success(`Published v${v.versionNumber}`),
            onError: (e) =>
              setErrors([{ field: '', message: `Draft saved but publish failed: ${(e as Error).message}` }]),
          })
        },
        onError: (e) =>
          setErrors(e instanceof ApiError && e.fields ? e.fields : [{ field: '', message: (e as Error).message }]),
      },
    )
  }

  const globalErrs = errors.filter((e) => !e.field)
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <SchemaForm schema={schema.data} config={config} errors={errors} onChange={setConfig} />
      {globalErrs.length > 0 && <Alert type="error" showIcon message={globalErrs.map((e) => e.message).join('; ')} />}
      <Button type="primary" loading={save.isPending || publish.isPending} onClick={saveAndPublish}>
        Save draft + Publish
      </Button>
    </Space>
  )
}
