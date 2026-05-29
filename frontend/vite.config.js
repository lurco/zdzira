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
      '/api/v1/events': {
        target: backendUrl,
        changeOrigin: true,
        // selfHandleResponse bypasses http-proxy's default response handling,
        // which can buffer chunks. We copy headers and pipe directly so each
        // SSE frame is forwarded to the browser as soon as the backend writes it.
        selfHandleResponse: true,
        configure: proxy => {
          proxy.on('proxyReq', proxyReq => { proxyReq.setTimeout(0) })
          proxy.on('proxyRes', (proxyRes, _req, res) => {
            res.writeHead(proxyRes.statusCode, proxyRes.headers)
            proxyRes.pipe(res)
          })
          proxy.on('error', (_err, _req, res) => {
            if (!res.headersSent) res.writeHead(502)
            res.end()
          })
        },
      },
      '/api': { target: backendUrl, changeOrigin: true },
      '/mcp': { target: backendUrl, ws: true },
    },
  },
})
