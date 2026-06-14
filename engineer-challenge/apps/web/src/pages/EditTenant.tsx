import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { App, Button, Card, Space, Spin, Steps, Typography } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { useActiveConfig, useConfigSchema, usePublish, useSaveDraft, useTenant } from '../api/hooks'
import { ApiError } from '../api/client'
import { SchemaForm } from '../schemaform/SchemaForm'
import { webhookRequired } from '../schemaform/conditional'
import type { ConfigDocument, FieldError } from '../api/types'

// Same dimension steps as the create page, minus Basics — the tenant already
// exists, so editing starts straight at the first dimension. There is no Review
// step: a Save button above the stepper publishes the whole config at any point.
const DIMENSION_STEPS = [
  { key: 'branding', title: 'Branding' },
  { key: 'claimTypes', title: 'Claim Types' },
  { key: 'approval', title: 'Approval Rules' },
  { key: 'notifications', title: 'Notifications' },
  { key: 'sla', title: 'SLA' },
  { key: 'customFields', title: 'Custom Fields' },
] as const

// dimensionOf("approval.tiers[0].maxAmount") -> "approval"
const dimensionOf = (field: string) => field.split(/[.[]/)[0]
// 0-based index of the step owning a dimension, or -1 if unknown.
const stepOfDimension = (key: string) => DIMENSION_STEPS.findIndex((d) => d.key === key)

// Dedicated edit page (/tenants/:slug/edit): the same dimension steps as Create,
// seeded from the tenant's active published config. Unlike Create, saving is not
// gated behind a final Review — the Save button above the stepper enables as soon
// as the config differs from what's published, and publishes the whole config.
export function EditTenant() {
  const { slug = '' } = useParams()
  const navigate = useNavigate()
  const { message } = App.useApp()
  const tenant = useTenant(slug)
  const schema = useConfigSchema()
  const active = useActiveConfig(slug)
  const save = useSaveDraft(slug)
  const publish = usePublish(slug)

  const [step, setStep] = useState(0)
  const [config, setConfig] = useState<ConfigDocument | null>(null)
  // Baseline snapshot of the last-published config; dirty = config differs from it.
  const [baseline, setBaseline] = useState<string | null>(null)
  const [errors, setErrors] = useState<FieldError[]>([])

  // Seed the editable buffer + baseline from the active config on first load
  // (adjust-state-during-render pattern; setByPath keeps edits immutable).
  if (config === null && active.data) {
    setConfig(active.data)
    setBaseline(JSON.stringify(active.data))
  }

  const dirty = config !== null && baseline !== null && JSON.stringify(config) !== baseline
  const publishing = save.isPending || publish.isPending

  function submit() {
    if (!config || !dirty) return
    // Cross-field rule the editor also enforces: webhook channel needs a URL.
    if (webhookRequired(config) && !config.notifications?.webhookUrl) {
      setErrors([{ field: 'notifications.webhookUrl', message: 'Webhook URL is required when the webhook channel is enabled' }])
      setStep(stepOfDimension('notifications'))
      return
    }
    save.mutate(
      { config, note: 'edited via admin' },
      {
        onSuccess: (v) => {
          setErrors([])
          publish.mutate(v.versionNumber, {
            onSuccess: () => {
              // New baseline → the Save button disables until the next change.
              setBaseline(JSON.stringify(config))
              message.success(`Published v${v.versionNumber}`)
            },
            onError: (e) => message.error(`Draft saved but publish failed: ${(e as Error).message}`),
          })
        },
        onError: (er) => {
          if (er instanceof ApiError && er.fields?.length) {
            setErrors(er.fields)
            const target = stepOfDimension(dimensionOf(er.fields[0].field))
            setStep(target >= 0 ? target : step)
            message.error('Please fix the highlighted fields')
          } else {
            message.error(er instanceof ApiError ? er.message : 'Failed to save configuration')
          }
        },
      },
    )
  }

  if (tenant.isError) return <Typography.Text type="danger">Tenant not found.</Typography.Text>

  const items = DIMENSION_STEPS.map((d) => ({ title: d.title }))
  const dim = schema.data?.dimensions.find((d) => d.key === DIMENSION_STEPS[step].key)

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', gap: 16 }}>
        <Space direction="vertical" size={4}>
          <Button type="link" icon={<ArrowLeftOutlined />} style={{ padding: 0 }} onClick={() => navigate(`/t/${slug}`)}>
            {tenant.data?.name ?? slug}
          </Button>
          <Typography.Title level={3} style={{ margin: 0 }}>
            Edit Configuration
          </Typography.Title>
        </Space>
        {/* Save lives above the stepper and enables only when there are changes. */}
        <Button type="primary" disabled={!dirty} loading={publishing} onClick={submit}>
          Save &amp; Publish
        </Button>
      </div>

      {/* Vertical label placement keeps all six step titles readable; steps are clickable. */}
      <Steps current={step} onChange={setStep} size="small" labelPlacement="vertical" items={items} />

      <Card>
        {!schema.data || !config || !dim ? (
          <Spin />
        ) : (
          <>
            <SchemaForm
              schema={{ dimensions: [dim] }}
              config={config}
              errors={errors.filter((e) => dimensionOf(e.field) === dim.key)}
              onChange={setConfig}
            />
            <div style={{ textAlign: 'right', marginTop: 24 }}>
              <Button disabled={step === 0} onClick={() => setStep(step - 1)} style={{ marginRight: 8 }}>
                Back
              </Button>
              <Button disabled={step === DIMENSION_STEPS.length - 1} onClick={() => setStep(step + 1)}>
                Next
              </Button>
            </div>
          </>
        )}
      </Card>
    </Space>
  )
}
