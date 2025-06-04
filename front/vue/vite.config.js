import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src') // 设置别名
    },
    extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue']
  },
  // 开发服务器配置
  server: {
    port: 7070,
    open: true,
    hmr: true, // 启用热模块替换
    host:  '0.0.0.0',// 允许外部访问
  },
  // 构建配置
  build: {
    outDir: 'dist',
    assetsInlineLimit: 4096, // 小于4kb的资源内联
    rollupOptions: {
      output: {
        manualChunks: {
          // 拆分第三方库
          vue: ['vue', 'vue-router'],
          lodash: ['lodash']
        }
      }
    }
  }
})