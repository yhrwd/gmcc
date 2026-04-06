import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import { mapResourceSnapshot } from '@/lib/mappers'
import type { LoadState, ViewResourceSnapshot } from '@/types/view'

export const useResourcesStore = defineStore('resources', {
  state: () => ({
    state: 'idle' as LoadState,
    errorMessage: '',
    snapshot: null as ViewResourceSnapshot | null,
  }),
  actions: {
    async loadResources(clusterStatus = 'running') {
      this.state = 'loading'
      this.errorMessage = ''
      try {
        const result = await apiClient.getResources()
        this.snapshot = mapResourceSnapshot(result, clusterStatus)
        this.state = 'success'
      } catch (error) {
        this.state = 'error'
        this.errorMessage = error instanceof Error ? error.message : '资源读取失败'
      }
    },
  },
})
