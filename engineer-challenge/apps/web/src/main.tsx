import '@ant-design/v5-patch-for-react-19'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter } from 'react-router-dom'
import { StyleProvider } from '@ant-design/cssinjs'
import { ConfigProvider, App as AntApp } from 'antd'
import App from './App'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: { queries: { retry: 1, refetchOnWindowFocus: false } },
})

// Theme tokens tuned to the design mockup (blue primary, rounded surfaces).
const theme = {
  token: { colorPrimary: '#2f54eb', borderRadius: 8 },
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    {/* `layer` keeps antd's cssinjs styles in a CSS @layer so Tailwind 4 utilities resolve predictably. */}
    <StyleProvider layer>
      <ConfigProvider theme={theme}>
        <AntApp>
          <QueryClientProvider client={queryClient}>
            <BrowserRouter>
              <App />
            </BrowserRouter>
          </QueryClientProvider>
        </AntApp>
      </ConfigProvider>
    </StyleProvider>
  </React.StrictMode>,
)
