import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import { mapResourceSnapshot, mapStatusSummary } from '@/lib/mappers'
import type { LoadState, ViewResourceSnapshot, ViewStatusSummary } from '@/types/view'

export const useHomeStore = defineStore('home', {
  state: () => ({
    statusState: 'idle' as LoadState,
    resourcesState: 'idle' as LoadState,
    statusError: '',
    resourcesError: '',
    statusSummary: null as ViewStatusSummary | null,
    resourceSnapshot: null as ViewResourceSnapshot | null,
  }),
  actions: {
    async loadStatus(force = false, silent = false) {
      if (!force && this.statusState === 'success') {
        return
      }
      if (!silent || this.statusState === 'idle') {
        this.statusState = 'loading'
      }
      this.statusError = ''
      try {
        const status = await apiClient.getStatus()
        this.statusSummary = mapStatusSummary(status)
        this.statusState = 'success'
        if (this.resourceSnapshot) {
          this.resourceSnapshot = {
            ...this.resourceSnapshot,
            comfortLevel: mapResourceSnapshot({
              cpu_percent: this.resourceSnapshot.cpuPercent ?? undefined,
              memory: { used_percent: this.resourceSnapshot.memoryPercent ?? undefined },
            }, this.statusSummary.clusterStatus).comfortLevel,
          }
        }
      } catch (error) {
        this.statusState = 'error'
        this.statusError = error instanceof Error ? error.message : '状态获取失败'
        throw error
      }
    },
    async loadResources(force = false, silent = false) {
      if (!force && this.resourcesState === 'success') {
        return
      }
      if (!silent || this.resourcesState === 'idle') {
        this.resourcesState = 'loading'
      }
      this.resourcesError = ''
      try {
        const resources = await apiClient.getResources()
        this.resourceSnapshot = mapResourceSnapshot(resources, this.statusSummary?.clusterStatus)
        this.resourcesState = 'success'
      } catch (error) {
        this.resourcesState = 'error'
        this.resourcesError = error instanceof Error ? error.message : '资源获取失败'
        throw error
      }
    },
    async loadHome(force = false, silent = false) {
      const results = await Promise.allSettled([
        this.loadStatus(force, silent),
        this.loadResources(force, silent),
      ])
      const firstRejected = results.find((result) => result.status === 'rejected')
      if (firstRejected) {
        throw firstRejected.reason
      }
    },
    async retryModule(module: 'status' | 'resources') {
      if (module === 'status') return this.loadStatus()
      return this.loadResources()
    },
  },
})
