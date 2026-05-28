import { defineConfig } from 'vite'
import pug from 'vite-plugin-pug'
import { resolve, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))

// In Docker dev the backend is reachable as http://backend:8080.
// Outside Docker (bare npm run dev) it falls back to localhost:8080.
const backendUrl = process.env.VITE_BACKEND_URL || 'http://localhost:8080'

export default defineConfig({
  root: resolve(__dirname, 'src'),
  plugins: [pug({ pretty: true, localImports: true })],
  build: {
    outDir: resolve(__dirname, '../dist'),
    emptyOutDir: true,
    rollupOptions: {
      input: {
        index: resolve(__dirname, 'src/index.html'),
        board: resolve(__dirname, 'src/board.html'),
        issue: resolve(__dirname, 'src/issue.html'),
      },
    },
  },
  server: {
    host: true,   // bind 0.0.0.0 so the port is reachable outside the container
    port: 5173,
    proxy: {
      '/api': { target: backendUrl },
      '/mcp': { target: backendUrl, ws: true },
    },
  },
})
