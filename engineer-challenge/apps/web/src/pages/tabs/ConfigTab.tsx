import { useNavigate } from 'react-router-dom'
import { Button, Descriptions, Space, Tag } from 'antd'
import { EditOutlined } from '@ant-design/icons'
import { useActiveConfig } from '../../api/hooks'
import { claimTypeLabel } from '../../labels'

// Read-only overview of the tenant's active published config. Editing happens on
// the dedicated /tenants/:slug/edit page (the same stepped wizard as Create).
export function ConfigTab({ slug }: { slug: string }) {
  const navigate = useNavigate()
  const active = useActiveConfig(slug)
  if (!active.data) return <div>Loading…</div>

  const config = active.data
  const enabledTypes = Object.entries(config.claimTypes)
    .filter(([, v]) => v.enabled)
    .map(([k]) => k)

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <div style={{ textAlign: 'right' }}>
        <Button type="primary" icon={<EditOutlined />} onClick={() => navigate(`/tenants/${slug}/edit`)}>
          Edit configuration
        </Button>
      </div>
      <Descriptions column={1} bordered size="small" title="Active configuration">
        <Descriptions.Item label="Display name">{config.branding.displayName}</Descriptions.Item>
        <Descriptions.Item label="Primary color">
          {config.branding.primaryColor ? (
            // Real flex row so the swatch is centered against the text, not baseline-aligned
            // inside a Space line-box (which left it sitting a few px high).
            <span style={{ display: 'inline-flex', alignItems: 'center', gap: 8 }}>
              <span
                style={{
                  width: 16,
                  height: 16,
                  borderRadius: 4,
                  border: '1px solid #d9d9d9',
                  background: config.branding.primaryColor,
                  flexShrink: 0,
                }}
              />
              {config.branding.primaryColor}
            </span>
          ) : (
            '—'
          )}
        </Descriptions.Item>
        <Descriptions.Item label="Enabled claim types">
          {enabledTypes.length ? (
            <Space size={[4, 4]} wrap>
              {enabledTypes.map((k) => (
                <Tag key={k} style={{ marginInlineEnd: 0 }}>
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
    </Space>
  )
}
