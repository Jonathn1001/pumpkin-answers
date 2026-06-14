import { Routes, Route, Link, Navigate } from 'react-router-dom'
import { TenantList } from './pages/TenantList'
import { TenantDetail } from './pages/TenantDetail'
import { Compare } from './pages/Compare'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      <header className="border-b bg-white px-6 py-3">
        <Link to="/" className="font-semibold">Claims Config Admin</Link>
        <Link to="/compare" className="ml-6 text-sm text-blue-700">Compare</Link>
      </header>
      <main className="mx-auto max-w-6xl p-6">
        <Routes>
          <Route path="/" element={<TenantList />} />
          <Route path="/t/:slug" element={<TenantDetail />} />
          <Route path="/compare" element={<Compare />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  )
}
