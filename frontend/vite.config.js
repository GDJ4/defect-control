import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

const buildProxy = (apiBase) => {
  try {
    const url = new URL(apiBase)
    const basePath = url.pathname.replace(/\/$/, '') || '/'
    return {
      '/api': {
        target: `${url.protocol}//${url.host}`,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, basePath),
      },
    }
  } catch {
    return undefined
  }
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const proxy = buildProxy(env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1')

  return {
    plugins: [vue()],
    server: {
      host: '0.0.0.0',
      port: 5173,
      proxy,
    },
    preview: {
      host: '0.0.0.0',
      port: 4173,
    },
  }
})
