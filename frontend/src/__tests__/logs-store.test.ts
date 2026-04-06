import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import { useLogsStore } from '@/stores/logs'
import { apiClient } from '@/api/client'

describe('logs store', () => {
  it('tracks active date filters', async () => {
    setActivePinia(createPinia())
    const store = useLogsStore()
    const originalGetOperationLogs = apiClient.getOperationLogs
    apiClient.getOperationLogs = async () => ({ logs: [] })
    await store.applyDateFilter({ start: '2026-04-06T08:00', end: '' })
    expect(store.hasActiveFilters).toBe(true)
    apiClient.getOperationLogs = originalGetOperationLogs
  })

  it('clears active date filters', () => {
    setActivePinia(createPinia())
    const store = useLogsStore()
    store.filterStart = '2026-04-06T08:00'
    store.filterEnd = '2026-04-06T10:00'
    store.clearDateFilter()
    expect(store.hasActiveFilters).toBe(false)
  })

  it('serializes datetime-local filters to RFC3339 when loading filtered logs', async () => {
    setActivePinia(createPinia())
    const store = useLogsStore()
    const originalGetOperationLogs = apiClient.getOperationLogs
    let captured: { start?: string; end?: string } | undefined
    apiClient.getOperationLogs = async (params) => {
      captured = params
      return { logs: [] }
    }

    await store.applyDateFilter({ start: '2026-04-06T08:00', end: '2026-04-06T10:30' })

    expect(captured?.start).toBe(new Date('2026-04-06T08:00').toISOString())
    expect(captured?.end).toBe(new Date('2026-04-06T10:30').toISOString())
    apiClient.getOperationLogs = originalGetOperationLogs
  })
})
