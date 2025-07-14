import { defineConfig } from 'electron-vite'
import { resolve } from 'path'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  main: {
    build: {
      rollupOptions: {
        input: {
          index: resolve(__dirname, 'src/main/index.js')
        }
      }
    }
  },
  preload: {
    build: {
      rollupOptions: {
        input: {
          preload: resolve(__dirname, 'src/preload/index.js')
        },
        output: {
            // 修改preload的输出文件名
            entryFileNames: (chunkInfo) => {
            if (chunkInfo.name === 'preload') {
                return 'index.js';
            }
            return '[name].js';
            },
        },
      }
    }
  },
  renderer: {
    root: resolve(__dirname, 'src/renderer'),
    resolve: {
        alias: {
            '@': resolve(__dirname,'src/renderer')
        }
    },
    plugins: [vue()],
    build: {
      rollupOptions: {
        input: {
          index: resolve(__dirname, 'src/renderer/index.html')
        }
      }
    }
  }
})