import { useEffect } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import { Button, Space, Tabs, Typography } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { useActiveConfig, useTenant } from '../api/hooks'
import { applyBranding } from '../theme/applyBranding'
import { ConfigTab } from './tabs/ConfigTab'
import { PreviewTab } from './tabs/PreviewTab'
import { VersionsTab } from './tabs/VersionsTab'
import { CompareTab } from './tabs/CompareTab'

export function TenantDetail() {
  const { slug = '' } = useParams()
  const navigate = useNavigate()
  const [params] = useSearchParams()
  const tab = params.get('tab') ?? 'config' // deep-linked from the tenant list (Preview/History)
  const tenant = useTenant(slug)
  const config = useActiveConfig(slug)
  useEffect(() => {
    applyBranding(config.data?.branding)
    return () => applyBranding(undefined)
  }, [config.data])
  if (tenant.isError) return <Typography.Text type="danger">Tenant not found.</Typography.Text>
  return (
    <Space direction="vertical" size="large" style={{ width: '100%' }}>
      <Space direction="vertical" size={4} style={{ width: '100%' }}>
        <Button type="link" icon={<ArrowLeftOutlined />} style={{ padding: 0 }} onClick={() => navigate('/')}>
          Tenants
        </Button>
        <Typography.Title level={3} style={{ margin: 0, color: 'var(--brand-primary)' }}>
          {tenant.data?.name ?? slug}
        </Typography.Title>
      </Space>
      <Tabs
        key={`${slug}:${tab}`}
        defaultActiveKey={tab}
        items={[
          { key: 'config', label: 'Config', children: <ConfigTab slug={slug} /> },
          { key: 'preview', label: 'Preview', children: <PreviewTab slug={slug} /> },
          { key: 'versions', label: 'Versions', children: <VersionsTab slug={slug} /> },
          { key: 'compare', label: 'Compare', children: <CompareTab slug={slug} /> },
        ]}
      />
    </Space>
  )
}
