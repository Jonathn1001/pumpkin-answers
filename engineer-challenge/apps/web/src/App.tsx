import { Routes, Route, Navigate, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Typography } from 'antd'
import { TeamOutlined, ThunderboltOutlined, DiffOutlined } from '@ant-design/icons'
import { TenantList } from './pages/TenantList'
import { TenantDetail } from './pages/TenantDetail'
import { CreateTenant } from './pages/CreateTenant'
import { EditTenant } from './pages/EditTenant'
import { Compare } from './pages/Compare'
import { Runtime } from './pages/Runtime'

const { Sider, Header, Content } = Layout

const NAV = [
  { key: '/', icon: <TeamOutlined />, label: 'Tenants' },
  { key: '/runtime', icon: <ThunderboltOutlined />, label: 'Runtime' },
  { key: '/compare', icon: <DiffOutlined />, label: 'Compare' },
]

export default function App() {
  const navigate = useNavigate()
  const { pathname } = useLocation()
  // /t/:slug detail and /tenants/* pages belong under the Tenants section.
  const selectedKey = pathname.startsWith('/t/') || pathname.startsWith('/tenants')
    ? '/'
    : NAV.map((n) => n.key).find((k) => (k === '/' ? pathname === '/' : pathname.startsWith(k))) ?? '/'

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg" collapsedWidth={0} theme="dark">
        <div
          style={{
            height: 56,
            margin: 16,
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            color: '#fff',
            fontWeight: 700,
            fontSize: 18,
          }}
        >
          <span style={{ fontSize: 22 }}>◆</span> ClaimConfig
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={NAV}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', borderBottom: '1px solid #f0f0f0', paddingInline: 24, display: 'flex', alignItems: 'center' }}>
          <Typography.Text strong style={{ fontSize: 16 }}>
            Multi-Tenant Claims Configuration
          </Typography.Text>
        </Header>
        <Content style={{ margin: 24 }}>
          <Routes>
            <Route path="/" element={<TenantList />} />
            <Route path="/tenants/new" element={<CreateTenant />} />
            <Route path="/tenants/:slug/edit" element={<EditTenant />} />
            <Route path="/t/:slug" element={<TenantDetail />} />
            <Route path="/runtime" element={<Runtime />} />
            <Route path="/compare" element={<Compare />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  )
}
