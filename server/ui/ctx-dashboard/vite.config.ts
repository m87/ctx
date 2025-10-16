/// <reference types='vitest' />
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths'
import { viteStaticCopy } from 'vite-plugin-static-copy'

export default defineConfig(() => ({
  root: __dirname,
  cacheDir: './node_modules/.vite/ctx-dashboard',
  server: {
    port: 4200,
    host: 'localhost',
     proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      secure: false,
    },
  },
  },
  preview: {
    port: 4300,
    host: 'localhost',
  },
  plugins: [react(),     tsconfigPaths(),    viteStaticCopy({
      targets: [{ src: '*.md', dest: '' }],
    }),],
  resolve: {
    alias: {
  
    },
  },

  build: {
    outDir: './dist/ctx-dashboard',
    emptyOutDir: true,
    reportCompressedSize: true,
    commonjsOptions: {
      transformMixedEsModules: true,
    },
  },
  test: {
    watch: false,
    globals: true,
    environment: 'jsdom',
    include: ['{src,tests}/**/*.{test,spec}.{js,mjs,cjs,ts,mts,cts,jsx,tsx}'],
    reporters: ['default'],
    coverage: {
      reportsDirectory: './coverage/ctx-dashboard',
      provider: 'v8' as const,
    },
  },
}));
