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
        index:  resolve(__dirname, 'src/index.html'),
        board:  resolve(__dirname, 'src/board.html'),
      },
    },
  },
  server: {
    host: true,   // bind 0.0.0.0 so the port is reachable outside the container
    port: 5173,
    proxy: {
      // SSE stream: no timeout so the long-lived connection is not cut.
      // http-proxy streams the response as it arrives; the backend disables
      // buffering (X-Accel-Buffering: no) so frames reach the browser live.
      '/api/v1/events': {
        target: backendUrl,
        changeOrigin: true,
        configure: proxy => {
          proxy.on('proxyReq', proxyReq => { proxyReq.setTimeout(0) })
        },
      },
      '/api': { target: backendUrl, changeOrigin: true },
      '/mcp': { target: backendUrl, ws: true },
    },
  },
})
