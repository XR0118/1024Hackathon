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
    pending: 'secondary',
    running: 'primary',
    success: 'success',
    failed: 'danger',
    waiting_confirm: 'warning',
  }
  return colorMap[status]
}

export const getStatusText = (status: Deployment['status']): string => {
  const textMap: Record<Deployment['status'], string> = {
    pending: '待开始',
    running: '进行中',
    success: '成功',
    failed: '失败',
    waiting_confirm: '待确认',
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
