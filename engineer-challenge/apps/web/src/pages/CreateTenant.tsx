import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Alert, App, Button, Card, Descriptions, Form, Input, Result, Select, Space, Spin, Steps, Tag, Typography } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { useConfigSchema, useCreateTenant, useTenants } from '../api/hooks'
import { ApiError, request } from '../api/client'
import { SchemaForm } from '../schemaform/SchemaForm'
import { webhookRequired } from '../schemaform/conditional'
import { claimTypeColor, claimTypeLabel } from '../labels'
import type { ConfigDocument, FieldError, Tenant } from '../api/types'

// One wizard step per config dimension, in display order. Each renders through
// the same SchemaForm the editor uses, scoped to its single dimension.
const DIMENSION_STEPS = [
  { key: 'branding', title: 'Branding' },
  { key: 'claimTypes', title: 'Claim Types' },
  { key: 'approval', title: 'Approval Rules' },
  { key: 'notifications', title: 'Notifications' },
  { key: 'sla', title: 'SLA' },
  { key: 'customFields', title: 'Custom Fields' },
] as const

const BASICS = 0
const REVIEW = DIMENSION_STEPS.length + 1

// dimensionOf("approval.tiers[0].maxAmount") -> "approval"
const dimensionOf = (field: string) => field.split(/[.[]/)[0]
const stepOfDimension = (key: string) => DIMENSION_STEPS.findIndex((d) => d.key === key) + 1

// Full "Add New Tenant" page: collect name + starting config, configure every
// dimension, then create. It lives on its own route (/tenants/new) because the
// flow carries a lot of information. The client sends the full config; the server
// validates it and reports field errors, which we route back to the owning step.
export function CreateTenant() {
  const navigate = useNavigate()
  const { message } = App.useApp()
  const create = useCreateTenant()
  const schema = useConfigSchema()
  const tenants = useTenants()

  const [step, setStep] = useState(BASICS)
  const [name, setName] = useState('')
  const [source, setSource] = useState('default') // 'default' | <tenant slug>
  const [config, setConfig] = useState<ConfigDocument | null>(null)
  const [seededFrom, setSeededFrom] = useState<string | null>(null)
  const [seeding, setSeeding] = useState(false)
  const [errors, setErrors] = useState<FieldError[]>([])
  const [created, setCreated] = useState<Tenant | null>(null)

  // Seed the editable config buffer from the chosen source, then enter the
  // dimension steps. Re-fetch only when the source changed (keeps edits on Back).
  async function startConfig() {
    if (!name.trim()) return
    if (seededFrom === source && config) {
      setStep(1)
      return
    }
    setSeeding(true)
    try {
      const seed =
        source === 'default'
          ? await request<ConfigDocument>('GET', '/config-default')
          : await request<ConfigDocument>('GET', `/tenants/${source}/config`)
      setConfig(seed)
      setSeededFrom(source)
      setErrors([])
      setStep(1)
    } catch {
      message.error('Failed to load the starting configuration')
    } finally {
      setSeeding(false)
    }
  }

  function submit() {
    if (!config) return
    // Cross-field rule the editor also enforces: webhook channel needs a URL.
    if (webhookRequired(config) && !config.notifications?.webhookUrl) {
      setErrors([{ field: 'notifications.webhookUrl', message: 'Webhook URL is required when the webhook channel is enabled' }])
      setStep(stepOfDimension('notifications'))
      return
    }
    create.mutate(
      { name, config },
      {
        onSuccess: setCreated,
        onError: (er) => {
          if (er instanceof ApiError && er.fields?.length) {
            setErrors(er.fields)
            setStep(stepOfDimension(dimensionOf(er.fields[0].field)) || REVIEW)
            message.error('Please fix the highlighted fields')
          } else {
            message.error(er instanceof ApiError ? er.message : 'Failed to create tenant')
          }
        },
      },
    )
  }

  const items = [{ title: 'Basics' }, ...DIMENSION_STEPS.map((d) => ({ title: d.title })), { title: 'Review & Create' }]
  const dim =
    step >= 1 && step <= DIMENSION_STEPS.length && schema.data
      ? schema.data.dimensions.find((d) => d.key === DIMENSION_STEPS[step - 1].key)
      : undefined
  const enabledTypes = config ? Object.entries(config.claimTypes).filter(([, v]) => v.enabled).map(([k]) => k) : []

  if (created) {
    return (
      <Card style={{ maxWidth: 880, margin: '0 auto' }}>
        <Result
          status="success"
          title="Tenant Created Successfully!"
          subTitle={`${created.name} · ${created.slug}`}
          extra={[
            <Button key="go" type="primary" onClick={() => navigate(`/t/${created.slug}`)}>
              Go to Tenant
            </Button>,
            <Button key="list" onClick={() => navigate('/')}>
              Back to Tenants
            </Button>,
          ]}
        />
      </Card>
    )
  }

  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space direction="vertical" size={4} style={{ width: '100%' }}>
        <Button type="link" icon={<ArrowLeftOutlined />} style={{ padding: 0 }} onClick={() => navigate('/')}>
          Tenants
        </Button>
        <Typography.Title level={3} style={{ margin: 0 }}>
          Add New Tenant
        </Typography.Title>
      </Space>

      {/* Vertical label placement keeps all eight step titles fully readable across the row. */}
      <Steps current={step} size="small" labelPlacement="vertical" items={items} />

      <Card>
        {step === BASICS && (
          <Form layout="vertical">
            <Form.Item label="Tenant name" required help="The URL slug is generated from the name by the server.">
              <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. SafeGuard Insurance" autoFocus />
            </Form.Item>
            <Form.Item label="Start from">
              <Select
                value={source}
                onChange={(v) => setSource(v)}
                options={[
                  { value: 'default', label: 'Default configuration' },
                  ...(tenants.data ?? [])
                    .filter((t) => t.status === 'active')
                    .map((t) => ({ value: t.slug, label: `Clone from: ${t.name}` })),
                ]}
              />
            </Form.Item>
            <div style={{ textAlign: 'right' }}>
              <Button onClick={() => navigate('/')} style={{ marginRight: 8 }}>
                Cancel
              </Button>
              <Button type="primary" disabled={!name.trim()} loading={seeding} onClick={startConfig}>
                Next
              </Button>
            </div>
          </Form>
        )}

        {step >= 1 && step <= DIMENSION_STEPS.length &&
          (!config || !dim ? (
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
                <Button onClick={() => setStep(step - 1)} style={{ marginRight: 8 }}>
                  Back
                </Button>
                <Button type="primary" onClick={() => setStep(step + 1)}>
                  Next
                </Button>
              </div>
            </>
          ))}

        {step === REVIEW && config && (
          <>
            <Descriptions column={1} bordered size="small" title="Tenant Summary">
              <Descriptions.Item label="Name">{name}</Descriptions.Item>
              <Descriptions.Item label="URL slug">generated from name</Descriptions.Item>
              <Descriptions.Item label="Source">
                {source === 'default' ? 'Default configuration' : `Clone of ${source}`}
              </Descriptions.Item>
              <Descriptions.Item label="Enabled claim types">
                {enabledTypes.length ? (
                  <Space size={[4, 4]} wrap>
                    {enabledTypes.map((k) => (
                      <Tag key={k} color={claimTypeColor(k)} style={{ marginInlineEnd: 0 }}>
                        {claimTypeLabel(k)}
                      </Tag>
                    ))}
                  </Space>
                ) : (
                  '—'
                )}
              </Descriptions.Item>
              <Descriptions.Item label="Approval">
                {config.approval.model} · auto-approve ≤ {config.approval.autoApproveThreshold}
              </Descriptions.Item>
              <Descriptions.Item label="Notifications">{config.notifications.channels.join(', ') || '—'}</Descriptions.Item>
              <Descriptions.Item label="SLA">{config.sla.defaultDays} days</Descriptions.Item>
              <Descriptions.Item label="Custom fields">{config.customFields.length}</Descriptions.Item>
            </Descriptions>
            {errors.length > 0 && (
              <Alert
                style={{ marginTop: 12 }}
                type="error"
                showIcon
                message="Some fields are invalid — go back to the highlighted steps to fix them."
              />
            )}
            <div style={{ textAlign: 'right', marginTop: 24 }}>
              <Button onClick={() => setStep(step - 1)} style={{ marginRight: 8 }}>
                Back
              </Button>
              <Button type="primary" loading={create.isPending} onClick={submit}>
                Create Tenant
              </Button>
            </div>
          </>
        )}
      </Card>
    </Space>
  )
}
