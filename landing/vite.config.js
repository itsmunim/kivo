import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  base: '/kivo/',
  build: {
    outDir: 'dist',
    rollupOptions: {
      input: {
        index: fileURLToPath(new URL('./index.html', import.meta.url)),
        'web-console': fileURLToPath(new URL('./web-console.html', import.meta.url)),
      },
    },
  },
})