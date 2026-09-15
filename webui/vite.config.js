import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// The built assets are embedded into the kivo binary from
// internal/webui/dist, so Vite outputs there directly.
export default defineConfig({
  plugins: [react()],
  base: '/',
  build: {
    outDir: '../internal/webui/dist',
    emptyOutDir: true,
  },
  server: {
    // Proxy API calls to the running kivo server while developing.
    proxy: {
      '/api': 'http://localhost:3001',
    },
  },
})