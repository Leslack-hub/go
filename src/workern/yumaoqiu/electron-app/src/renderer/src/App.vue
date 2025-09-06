<template>
  <div class="main-container">
    <!-- 头部 -->
    <div class="header">
      <h1>羽毛球场地预订系统</h1>
      <div class="status-indicator">
        <span class="status-dot" :class="{ running: isRunning, stopped: !isRunning }"></span>
        <span>{{ isRunning ? '运行中' : '已停止' }}</span>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="content">
      <!-- 左侧配置面板 -->
      <div class="left-panel">
        <div class="form-section">
          <h3>基本配置</h3>
          <el-form :model="formData" label-width="100px" size="default">
                        
            <el-form-item label="用户账号" required>
              <el-input
                v-model="formData.netUserId"
                placeholder="请输入用户账号"
                style="width: 100%"
                clearable
              />
            </el-form-item>

            <el-form-item label="预订日期" required>
              <el-date-picker
                v-model="selectedDate"
                type="date"
                placeholder="选择日期"
                format="YYYY-MM-DD"
                value-format="YYYYMMDD"
                style="width: 100%"
                @change="onDateChange"
              />
            </el-form-item>
            
            <el-form-item label="执行次数">
              <el-input-number
                v-model="formData.times"
                :min="1"
                :max="100"
                placeholder="默认5次"
                style="width: 100%"
              />
            </el-form-item>
            
            <el-form-item label="开始时间">
              <el-date-picker
                v-model="startDateTime"
                type="datetime"
                placeholder="选择开始时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 100%"
                @change="onStartTimeChange"
              />
            </el-form-item>
            
            <el-form-item label="时间段选择">
              <el-select
                v-model="selectedTimeSegments"
                multiple
                :multiple-limit="2"
                placeholder="选择时间段（最多2项）"
                style="width: 100%"
                @change="onTimeSegmentChange"
              >
                <el-option
                  v-for="(timeSlot,index) in currentTimeSlots"
                  :key="index"
                  :label="timeSlot.label"
                  :value="timeSlot.segment"
                />
              </el-select>
            </el-form-item>

          </el-form>
        </div>

        <!-- 程序状态 - 仅在debug模式下显示 -->
        <div v-if="isDebugMode" class="form-section">
          <h3>程序状态</h3>
          <el-descriptions :column="1" size="small">
            <el-descriptions-item label="Go程序路径">
              <el-text size="small" type="info">{{ goExecutableInfo.path || '检测中...' }}</el-text>
            </el-descriptions-item>
            <el-descriptions-item label="程序状态">
              <el-tag :type="goExecutableInfo.exists ? 'success' : 'danger'" size="small">
                {{ goExecutableInfo.exists ? '已找到' : '未找到' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="程序类型">
              <el-tag :type="goExecutableInfo.isExecutable ? 'primary' : 'warning'" size="small">
                {{ goExecutableInfo.isExecutable ? '可执行文件' : 'Go源码' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 控制按钮 -->
        <div class="control-buttons">
          <el-button
            type="primary"
            :loading="isStarting"
            :disabled="!canStart"
            @click="startProgram"
            style="flex: 1"
          >
            <el-icon><VideoPlay /></el-icon>
            {{ isStarting ? '启动中...' : '开始预订' }}
          </el-button>
          
          <el-button
            type="danger"
            :disabled="!isRunning"
            @click="stopProgram"
            style="flex: 1"
          >
            <el-icon><VideoPause /></el-icon>
            停止程序
          </el-button>
        </div>
        
        <div class="control-buttons">
          <el-button
            @click="clearLogs"
            style="flex: 1"
          >
            <el-icon><Delete /></el-icon>
            清空日志
          </el-button>
          
          <el-button
            @click="refreshGoPath"
            style="flex: 1"
          >
            <el-icon><Refresh /></el-icon>
            刷新状态
          </el-button>
        </div>
      </div>

      <!-- 右侧日志面板 -->
      <div class="right-panel">
        <div class="log-container">
          <div class="log-header">
            <span>程序输出日志</span>
            <el-tag v-if="logs.length > 0" size="small" style="margin-left: 12px">
              {{ logs.length }} 条记录
            </el-tag>
          </div>
          <div class="log-content" ref="logContainer">
            <div v-if="logs.length === 0" class="log-line" style="color: #888; font-style: italic;">
              等待程序输出...
            </div>
            <div
              v-for="(log, index) in logs"
              :key="index"
              class="log-line"
              :class="getLogClass(log.type)"
            >
              <span class="log-time">[{{ log.time }}]</span>
              <span class="log-content-text">{{ log.content }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import {
  VideoPlay,
  VideoPause,
  Delete,
  Refresh
} from '@element-plus/icons-vue'

// 导入类型定义
import type { ElectronAPI } from './types/electron'

// 声明全局 window 接口扩展
declare global {
  interface Window {
    electronAPI: ElectronAPI
  }
}

// 响应式数据
const isRunning = ref(false)
const isStarting = ref(false)
const selectedDate = ref('')
const startDateTime = ref('')
const selectedTimeSegments = ref<number[]>([])
const logContainer = ref<HTMLElement>()

// Debug模式控制 - 默认关闭，可通过快捷键切换
const isDebugMode = ref(true)

// 时间段配置
const timeSlots = [
  { segment: 1, label: '12:00-13:00' },
  { segment: 2, label: '13:00-14:00' },
  { segment: 3, label: '14:00-15:00' },
  { segment: 4, label: '15:00-16:00' },
  { segment: 5, label: '16:00-17:00' },
  { segment: 6, label: '17:00-18:00' },
  { segment: 7, label: '18:00-19:00' },
  { segment: 8, label: '19:00-20:00' },
  { segment: 9, label: '20:00-21:00' },
  { segment: 10, label: '21:00-22:00' }
]

const timeSlots2 = [
  { segment: 1, label: '08:00-09:00' },
  { segment: 2, label: '09:00-10:00' },
  { segment: 3, label: '10:00-11:00' },
  { segment: 4, label: '11:00-12:00' },
  { segment: 5, label: '12:00-13:00' },
  { segment: 6, label: '12:00-13:00' },
  { segment: 7, label: '13:00-14:00' },
  { segment: 8, label: '14:00-15:00' },
  { segment: 9, label: '15:00-16:00' },
  { segment: 10, label: '16:00-17:00' },
  { segment: 11, label: '17:00-18:00' },
  { segment: 12, label: '18:00-19:00' },
  { segment: 13, label: '19:00-20:00' },
  { segment: 14, label: '20:00-21:00' },
  { segment: 15, label: '21:00-22:00' }
]

const formData = reactive({
  day: '',
  times: 5,
  startTime: '',
  location: '4,5',
  netUserId: ''
})

const goExecutableInfo = reactive({
  path: '',
  exists: false,
  isExecutable: false
})

interface LogEntry {
  time: string
  type: 'stdout' | 'stderr' | 'error' | 'info'
  content: string
}

const logs = ref<LogEntry[]>([])

// 计算属性
const canStart = computed(() => {
  return formData.day && formData.netUserId && goExecutableInfo.exists && !isRunning.value
})

// 根据选择的日期返回对应的时间段配置
const currentTimeSlots = computed(() => {
  if (!formData.day) {
    return timeSlots
  }
  
  // formData.day格式是YYYYMMDD，需要转换为YYYY-MM-DD格式
  const dateStr = formData.day
  const formattedDate = `${dateStr.slice(0,4)}-${dateStr.slice(4,6)}-${dateStr.slice(6,8)}`
  const date = new Date(formattedDate)
  const dayOfWeek = date.getDay() // 0=周日, 1=周一, ..., 6=周六
  
  // 周四(4)和周日(0)使用timeSlots2，其他使用timeSlots
  if (dayOfWeek === 4 || dayOfWeek === 0) {
    return timeSlots2
  } else {
    return timeSlots
  }
})

// 方法
const addLog = (type: LogEntry['type'], content: string) => {
  const now = new Date()
  const time = now.toLocaleTimeString()
  
  logs.value.push({
    time,
    type,
    content: content.trim()
  })
  
  // 自动滚动到底部
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

const getLogClass = (type: string) => {
  switch (type) {
    case 'stderr':
    case 'error':
      return 'error'
    case 'stdout':
      return 'success'
    default:
      return ''
  }
}

const onDateChange = (value: string) => {
  formData.day = value
  // 清空之前选择的时间段
  selectedTimeSegments.value = []
  formData.location = ''
}

const onTimeSegmentChange = (segments: number[]) => {
  // 根据选择的时间段计算对应的location值
  if (segments && segments.length > 0) {
    // formData.day格式是YYYYMMDD，需要转换为YYYY-MM-DD格式
    const dateStr = formData.day
    const formattedDate = `${dateStr.slice(0,4)}-${dateStr.slice(4,6)}-${dateStr.slice(6,8)}`
    const date = new Date(formattedDate)
    const dayOfWeek = date.getDay()
    
    let locations: number[]
    
    if (dayOfWeek === 4 || dayOfWeek === 0) {
      locations = segments.map(segment => segment)
    } else {

      locations = segments.map(segment => {
        return segment
      })
    }
    
    formData.location = locations.join(',')
  } else {
    formData.location = ''
  }
}

const onStartTimeChange = (value: string) => {
  formData.startTime = value
}

const startProgram = async () => {
  if (!formData.day) {
    ElMessage.error('请选择预订日期')
    return
  }
  
  if (!formData.netUserId) {
    ElMessage.error('请输入用户账号')
    return
  }
  
  if (typeof window === 'undefined' || !window.electronAPI) {
    ElMessage.error('当前运行在浏览器环境中，无法启动程序')
    return
  }
  
  isStarting.value = true
  
  try {
    const params = {
      day: formData.day,
      times: formData.times.toString(),
      startTime: formData.startTime,
      location: formData.location,
      netUserId: formData.netUserId
    }
    
    addLog('info', `开始启动程序，参数: ${JSON.stringify(params)}`)
    
    const result = await window.electronAPI.startGoProgram(params)
    
    if (result.success) {
      isRunning.value = true
      addLog('info', result.message)
      ElMessage.success('程序启动成功')
    } else {
      addLog('error', result.message)
      ElMessage.error(result.message)
    }
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : '启动失败'
    addLog('error', errorMsg)
    ElMessage.error(errorMsg)
  } finally {
    isStarting.value = false
  }
}

const stopProgram = async () => {
  if (typeof window === 'undefined' || !window.electronAPI) {
    ElMessage.error('当前运行在浏览器环境中，无法停止程序')
    return
  }
  
  try {
    const result = await window.electronAPI.stopGoProgram()
    
    if (result.success) {
      isRunning.value = false
      addLog('info', result.message)
      ElMessage.success('程序已停止')
    } else {
      addLog('error', result.message)
      ElMessage.error(result.message)
    }
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : '停止失败'
    addLog('error', errorMsg)
    ElMessage.error(errorMsg)
  }
}

const clearLogs = () => {
  logs.value = []
  ElMessage.success('日志已清空')
}

const refreshGoPath = async () => {
  if (typeof window === 'undefined' || !window.electronAPI) {
    addLog('info', '当前运行在浏览器环境中，无法检查Go程序状态')
    return
  }
  
  try {
    const info = await window.electronAPI.checkGoExecutable()
    Object.assign(goExecutableInfo, info)
    
    if (info.exists) {
      ElMessage.success('程序状态刷新成功')
    } else {
      ElMessage.warning('未找到Go程序，请检查文件路径')
    }
  } catch (error) {
    ElMessage.error('刷新状态失败')
  }
}

// 键盘事件处理 - 切换debug模式
const handleKeyDown = (event: KeyboardEvent) => {
  // Ctrl + Shift + D 切换debug模式
  if (event.ctrlKey && event.shiftKey && event.key === 'D') {
    event.preventDefault()
    isDebugMode.value = !isDebugMode.value
    ElMessage.info(`Debug模式: ${isDebugMode.value ? '开启' : '关闭'}`)
  }
}

// 生命周期
onMounted(async () => {
  // 检查Go程序状态
  await refreshGoPath()
  
  // 设置默认日期为明天
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const year = tomorrow.getFullYear()
  const month = String(tomorrow.getMonth() + 1).padStart(2, '0')
  const day = String(tomorrow.getDate()).padStart(2, '0')
  selectedDate.value = `${year}${month}${day}`
  formData.day = selectedDate.value
  
  // 添加键盘事件监听器
  document.addEventListener('keydown', handleKeyDown)
  
  // 检查是否在Electron环境中
  if (typeof window !== 'undefined' && window.electronAPI) {
    // 监听Go程序输出
    window.electronAPI.onGoOutput((data) => {
      addLog(data.type as LogEntry['type'], data.data)
    })
    
    // 监听Go程序退出
    window.electronAPI.onGoProcessExit((code) => {
      isRunning.value = false
      addLog('info', `程序退出，退出码: ${code}`)
      
      if (code === 0) {
        ElMessage.success('程序正常退出')
      } else {
        ElMessage.warning(`程序异常退出，退出码: ${code}`)
      }
    })
  } else {
    addLog('info', '当前运行在浏览器环境中，部分功能不可用')
  }
})

onUnmounted(() => {
  // 清理事件监听器
  document.removeEventListener('keydown', handleKeyDown)
  if (window.electronAPI) {
    window.electronAPI.removeAllListeners()
  }
})
</script>

<style scoped>
.log-time {
  color: #888;
  margin-right: 8px;
}

.log-content-text {
  word-break: break-all;
}

.el-descriptions {
  margin-top: 8px;
}

.el-descriptions :deep(.el-descriptions__body) {
  background: transparent;
}

.el-descriptions :deep(.el-descriptions__table) {
  border: none;
}

.el-descriptions :deep(.el-descriptions__cell) {
  border: none;
  padding: 4px 0;
}
</style>