import api from '@/services/api'
import { setupMockHandlers } from './handlers'

export function enableMockMode() {
  setupMockHandlers(api)
}

export * from './data'
