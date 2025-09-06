import { contextBridge, ipcRenderer } from 'electron'

// 定义API接口
export interface ElectronAPI {
  startGoProgram: (params: {
    day: string
    times?: string
    startTime?: string
    location?: string
  }) => Promise<{ success: boolean; message: string }>
  
  stopGoProgram: () => Promise<{ success: boolean; message: string }>
  
  getGoPath: () => Promise<string>
  
  checkGoExecutable: () => Promise<{
    path: string
    exists: boolean
    isExecutable: boolean
  }>
  
  showErrorDialog: (title: string, content: string) => Promise<void>
  
  onGoOutput: (callback: (data: { type: string; data: string }) => void) => void
  
  onGoProcessExit: (callback: (code: number | null) => void) => void
  
  removeAllListeners: () => void
}

// 暴露API到渲染进程
contextBridge.exposeInMainWorld('electronAPI', {
  startGoProgram: (params: any) => ipcRenderer.invoke('start-go-program', params),
  
  stopGoProgram: () => ipcRenderer.invoke('stop-go-program'),
  
  getGoPath: () => ipcRenderer.invoke('get-go-path'),
  
  checkGoExecutable: () => ipcRenderer.invoke('check-go-executable'),
  
  showErrorDialog: (title: string, content: string) => 
    ipcRenderer.invoke('show-error-dialog', title, content),
  
  onGoOutput: (callback: (data: { type: string; data: string }) => void) => {
    ipcRenderer.on('go-output', (event, data) => callback(data))
  },
  
  onGoProcessExit: (callback: (code: number | null) => void) => {
    ipcRenderer.on('go-process-exit', (event, code) => callback(code))
  },
  
  removeAllListeners: () => {
    ipcRenderer.removeAllListeners('go-output')
    ipcRenderer.removeAllListeners('go-process-exit')
  }
})

// 类型声明
declare global {
  interface Window {
    electronAPI: ElectronAPI
  }
}