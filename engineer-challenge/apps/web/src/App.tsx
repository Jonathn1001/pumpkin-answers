import { Routes, Route, Link, Navigate } from "react-router-dom";

export default function App() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      <header className="border-b bg-white px-6 py-3">
        <Link to="/" className="font-semibold">
          Claims Config Admin
        </Link>
        <Link to="/compare" className="ml-6 text-sm text-blue-700">
          Compare
        </Link>
      </header>
      <main className="mx-auto max-w-6xl p-6">
        <Routes>
          <Route path="/" element={<div>Tenants (Task 4)</div>} />
          <Route path="/t/:slug" element={<div>Tenant detail (Task 5)</div>} />
          <Route path="/compare" element={<div>Compare (Task 6)</div>} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  );
}
