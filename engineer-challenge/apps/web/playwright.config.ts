import { defineConfig } from '@playwright/test'

// Assumes the Go API (:8080) + Postgres are already running; Vite proxies /api → :8080.
export default defineConfig({
  testDir: './tests',
  use: { baseURL: 'http://localhost:5173' },
  webServer: { command: 'npm run dev', url: 'http://localhost:5173', reuseExistingServer: true, timeout: 60_000 },
})
