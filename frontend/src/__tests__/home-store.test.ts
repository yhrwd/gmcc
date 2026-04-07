import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { apiClient } from '@/api/client'
import { useHomeStore } from '@/stores/home'

describe('home store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('keeps partial success visible when one module fails', () => {
    const store = useHomeStore()
    store.$patch({ statusState: 'success', resourcesState: 'error' })
    expect(store.statusState).toBe('success')
    expect(store.resourcesState).toBe('error')
  })

  it('rejects when one child request fails while preserving partial local state', async () => {
    vi.spyOn(apiClient, 'getStatus').mockResolvedValue({
      cluster_status: 'running',
      total_instances: 1,
      running_instances: 1,
      uptime: '1m',
    })
    vi.spyOn(apiClient, 'getResources').mockRejectedValue(new Error('resource down'))

    const store = useHomeStore()

    await expect(store.loadHome(true, true)).rejects.toThrow('resource down')
    expect(store.statusState).toBe('success')
    expect(store.resourcesState).toBe('error')
    expect(store.statusError).toBe('')
    expect(store.resourcesError).toBe('resource down')
    expect(store.statusSummary?.clusterStatus).toBe('running')
  })
})
