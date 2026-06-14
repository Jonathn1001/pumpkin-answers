import { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Space, Tabs, Typography } from 'antd'
import { useActiveConfig, useTenant } from '../api/hooks'
import { applyBranding } from '../theme/applyBranding'
import { ConfigTab } from './tabs/ConfigTab'
import { PreviewTab } from './tabs/PreviewTab'
import { VersionsTab } from './tabs/VersionsTab'
import { CompareTab } from './tabs/CompareTab'

export function TenantDetail() {
  const { slug = '' } = useParams()
  const tenant = useTenant(slug)
  const config = useActiveConfig(slug)
  useEffect(() => {
    applyBranding(config.data?.branding)
    return () => applyBranding(undefined)
  }, [config.data])
  if (tenant.isError) return <Typography.Text type="danger">Tenant not found.</Typography.Text>
  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ margin: 0, color: 'var(--brand-primary)' }}>
        {tenant.data?.name ?? slug}
      </Typography.Title>
      <Tabs
        key={slug}
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
