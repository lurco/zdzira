import { defineConfig } from 'vite'
import pug from 'vite-plugin-pug'
import { resolve, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))

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
    port: 5173,
  },
})
