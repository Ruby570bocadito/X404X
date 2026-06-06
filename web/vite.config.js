import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8445',
      '/ws': {
        target: 'ws://localhost:8446',
        ws: true,
      },
    },
  },
})
