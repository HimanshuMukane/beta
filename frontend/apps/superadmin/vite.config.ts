import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import { compression } from 'vite-plugin-compression2'
import path from 'path'

export default defineConfig({
  plugins: [
    tailwindcss(),
    react(),
    // Pre-compress with Brotli (.br) - highest compression ratio
    compression({
      algorithm: 'brotliCompress',
      exclude: [/\.(br)$/, /\.(gz)$/],
      threshold: 1024,
    }),
    // Fallback Gzip (.gz)
    compression({
      algorithm: 'gzip',
      exclude: [/\.(br)$/, /\.(gz)$/],
      threshold: 1024,
    }),
  ],
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: [
      { find: '@', replacement: path.resolve(import.meta.dirname, './src') },
      { find: /^@repo\/ui\/(.*)$/, replacement: path.resolve(import.meta.dirname, '../../packages/ui/src/$1') },
      { find: '@repo/ui', replacement: path.resolve(import.meta.dirname, '../../packages/ui/src/index.ts') },
    ],
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules/react') || id.includes('node_modules/react-dom')) {
            return 'vendor-react';
          }
          if (id.includes('node_modules/@radix-ui')) {
            return 'vendor-radix';
          }
          if (id.includes('node_modules/@tabler/icons-react')) {
            return 'vendor-icons';
          }
          if (id.includes('node_modules/@tanstack')) {
            return 'vendor-tanstack';
          }
          if (id.includes('node_modules/echarts') || id.includes('node_modules/zrender')) {
            return 'vendor-echarts';
          }
          if (id.includes('node_modules/@tiptap') || id.includes('node_modules/prosemirror')) {
            return 'vendor-tiptap';
          }
          if (id.includes('node_modules/recharts')) {
            return 'vendor-recharts';
          }
        },
      },
    },
  },
})
