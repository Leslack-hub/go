import { app, BrowserWindow, ipcMain, dialog } from 'electron'
import { spawn, ChildProcess } from 'child_process'
import * as path from 'path'
import * as fs from 'fs'

const isDev = process.env.NODE_ENV === 'development'
let mainWindow: BrowserWindow | null = null
let goProcess: ChildProcess | null = null

// 获取Go程序路径
function getGoExecutablePath(): string {
  const appPath = app.getAppPath()

  console.log('Debug - App path:', appPath)
  console.log('Debug - Is dev mode:', isDev)

  // 根据平台检查不同的可执行文件
  let executableName: string
  switch (process.platform) {
    case 'win32':
      executableName = 'fetch_and_order.exe'
      break
    case 'darwin':
      executableName = 'fetch_and_order_darwin_arm64'
      break
    case 'linux':
      executableName = 'fetch_and_order'
      break
    default:
      executableName = 'fetch_and_order'
  }

  // 使用内置的resources目录中的Go程序
  let goExePath: string
  if (isDev) {
    // 开发模式：使用项目根目录下的resources文件夹
    goExePath = path.join(appPath, 'resources', executableName)
  } else {
    // 生产模式：使用打包后的resources文件夹
    // 在打包后，resources会被复制到应用程序目录
    goExePath = path.join(process.resourcesPath, 'resources', executableName)
  }

  console.log('Debug - Go exe path:', goExePath)

  // 检查可执行文件是否存在
  if (fs.existsSync(goExePath)) {
    return goExePath
  }

  // 如果内置文件不存在，回退到外部路径（向后兼容）
  const fallbackPath = path.join(appPath, '..', executableName)
  console.log('Debug - Fallback path:', fallbackPath)
  return fallbackPath
}

function createWindow(): void {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 1000,
    minHeight: 600,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    },
    titleBarStyle: 'default',
    show: false,
    icon: path.join(__dirname, '../../assets/icon.svg')
  })

  // 加载应用
  if (isDev) {
    mainWindow.loadURL('http://localhost:5173')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(path.join(__dirname, '../renderer/index.html'))
  }

  mainWindow.once('ready-to-show', () => {
    mainWindow?.show()
  })

  mainWindow.on('closed', () => {
    mainWindow = null
    // 关闭Go进程
    if (goProcess) {
      goProcess.kill()
      goProcess = null
    }
  })
}

// IPC处理程序
ipcMain.handle('start-go-program', async (event, params) => {
  try {
    if (goProcess) {
      goProcess.kill()
      goProcess = null
    }

    const goPath = getGoExecutablePath()
    const args: string[] = []

    // 构建命令行参数
    if (params.day) {
      args.push('-day', params.day)
    }
    if (params.times) {
      args.push('-times', params.times)
    }
    if (params.startTime) {
      args.push('-start', params.startTime)
    }
    if (params.location) {
      args.push('-location', params.location)
    }

    // 判断是运行exe还是go文件
    let command: string
    let commandArgs: string[]

    if (goPath.endsWith('.exe')) {
      command = goPath
      commandArgs = args
    } else {
      command = 'go'
      commandArgs = ['run', goPath, ...args]
    }

    goProcess = spawn(command, commandArgs, {
      cwd: path.dirname(goPath),
      stdio: ['pipe', 'pipe', 'pipe']
    })

    // 监听输出
    goProcess.stdout?.on('data', (data) => {
      const output = data.toString()
      mainWindow?.webContents.send('go-output', {
        type: 'stdout',
        data: output
      })
    })

    goProcess.stderr?.on('data', (data) => {
      const output = data.toString()
      mainWindow?.webContents.send('go-output', {
        type: 'stderr',
        data: output
      })
    })

    goProcess.on('close', (code) => {
      mainWindow?.webContents.send('go-process-exit', code)
      goProcess = null
    })

    goProcess.on('error', (error) => {
      mainWindow?.webContents.send('go-output', {
        type: 'error',
        data: `进程错误: ${error.message}`
      })
      goProcess = null
    })

    return { success: true, message: '程序启动成功' }
  } catch (error) {
    return {
      success: false,
      message: `启动失败: ${error instanceof Error ? error.message : '未知错误'}`
    }
  }
})

ipcMain.handle('stop-go-program', async () => {
  try {
    if (goProcess) {
      goProcess.kill()
      goProcess = null
      return { success: true, message: '程序已停止' }
    } else {
      return { success: false, message: '没有正在运行的程序' }
    }
  } catch (error) {
    return {
      success: false,
      message: `停止失败: ${error instanceof Error ? error.message : '未知错误'}`
    }
  }
})

ipcMain.handle('get-go-path', async () => {
  return getGoExecutablePath()
})

ipcMain.handle('check-go-executable', async () => {
  const goPath = getGoExecutablePath()
  const exists = fs.existsSync(goPath)

  // 添加调试日志
  console.log('Debug - Go path:', goPath)
  console.log('Debug - File exists:', exists)
  console.log('Debug - Platform:', process.platform)

  // 判断是否为可执行文件
  let isExecutable = false
  if (exists) {
    if (process.platform === 'win32') {
      isExecutable = goPath.endsWith('.exe')
    } else {
      // 在Unix系统上，检查文件是否有执行权限
      try {
        const stats = fs.statSync(goPath)
        isExecutable = !goPath.endsWith('.go') && (stats.mode & parseInt('111', 8)) !== 0
      } catch (error) {
        isExecutable = !goPath.endsWith('.go')
      }
    }
  }

  console.log('Debug - Is executable:', isExecutable)

  const result = {
    path: goPath,
    exists,
    isExecutable
  }

  console.log('Debug - Final result:', result)

  return result
})

ipcMain.handle('show-error-dialog', async (event, title: string, content: string) => {
  if (mainWindow) {
    dialog.showErrorBox(title, content)
  }
})

// 应用事件
app.whenReady().then(() => {
  createWindow()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    // 确保关闭Go进程
    if (goProcess) {
      goProcess.kill()
      goProcess = null
    }
    app.quit()
  }
})

app.on('before-quit', () => {
  // 应用退出前关闭Go进程
  if (goProcess) {
    goProcess.kill()
    goProcess = null
  }
})