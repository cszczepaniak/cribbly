import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  // Dev server: same URL shape as behind Go (`/admin/games`, not `/app/admin/games`).
  // Production build + vite preview: assets stay under `/app/` for the Go static handler.
  const isViteDevServer = command === 'serve' && mode === 'development'

  return {
    base: isViteDevServer ? '/' : '/app/',
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    build: {
      outDir: '../internal/web/embed/dist',
      emptyOutDir: true,
    },
    server: {
      port: 5173,
      strictPort: true,
      proxy: {
        // Wire API / Connect RPC to the Go server when you add endpoints.
        '/api': { target: 'http://127.0.0.1:8080', changeOrigin: true },
      },
    },
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './src/test/setup.ts',
    },
  }
})
