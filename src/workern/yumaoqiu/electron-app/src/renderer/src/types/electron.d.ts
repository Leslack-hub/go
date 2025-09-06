// Electron API 类型定义
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

// 全局类型声明
declare global {
  interface Window {
    electronAPI: ElectronAPI
  }
}

export {}