import { contextBridge } from 'electron'

// 安全暴露 API 给渲染进程
contextBridge.exposeInMainWorld('electronAPI', {
  platform: process.platform
})