import type { Deployment } from '@/types'

export const formatDate = (date: string): string => {
  return new Date(date).toLocaleString('zh-CN')
}

export const formatDuration = (seconds: number): string => {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分${seconds % 60}秒`
  return `${Math.floor(seconds / 3600)}小时${Math.floor((seconds % 3600) / 60)}分`
}

export const getStatusColor = (status: Deployment['status']): string => {
  const colorMap: Record<Deployment['status'], string> = {
    pending: 'warning',
    running: 'primary',
    paused: 'info',
    completed: 'success',
  }
  return colorMap[status]
}

export const getStatusText = (status: Deployment['status']): string => {
  const textMap: Record<Deployment['status'], string> = {
    pending: '待开始',
    running: '运行中',
    paused: '暂停中',
    completed: '完成',
  }
  return textMap[status]
}

export const debounce = <T extends (...args: any[]) => any>(
  fn: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: NodeJS.Timeout
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

// Environment 辅助函数
export const getEnvironmentTypeDisplay = (type: string): string => {
  return type === 'kubernetes' ? 'Kubernetes' : '物理机'
}

export const getEnvironmentTypeBadgeColor = (type: string): string => {
  return type === 'kubernetes' ? 'primary' : 'success'
}

export const getEnvironmentStatusDisplay = (isActive: boolean): string => {
  return isActive ? '运行中' : '已停止'
}

export const getEnvironmentStatusBadgeColor = (isActive: boolean): string => {
  return isActive ? 'success' : 'secondary'
}

// Application 辅助函数
export const getApplicationTypeDisplay = (type: string): string => {
  return type === 'microservice' ? '微服务' : '单体应用'
}

export const getApplicationTypeBadgeColor = (type: string): string => {
  return type === 'microservice' ? 'blue' : 'purple'
}
