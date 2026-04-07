import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import { mapLogItem } from '@/lib/mappers'
import type { LoadState, ViewLogItem } from '@/types/view'

function hasDateFilter(start: string, end: string) {
  return Boolean(start || end)
}

function toRFC3339(value: string) {
  if (!value) {
    return ''
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return date.toISOString()
}

export const useLogsStore = defineStore('logs', {
  state: () => ({
    state: 'idle' as LoadState,
    errorMessage: '',
    items: [] as ViewLogItem[],
    filteredState: 'idle' as LoadState,
    filteredErrorMessage: '',
    filteredItems: [] as ViewLogItem[],
    filterStart: '',
    filterEnd: '',
  }),
  getters: {
    hasActiveFilters(state) {
      return hasDateFilter(state.filterStart, state.filterEnd)
    },
    visibleItems(state) {
      return hasDateFilter(state.filterStart, state.filterEnd) ? state.filteredItems : state.items
    },
    currentState(state) {
      return hasDateFilter(state.filterStart, state.filterEnd) ? state.filteredState : state.state
    },
    currentErrorMessage(state) {
      return hasDateFilter(state.filterStart, state.filterEnd) ? state.filteredErrorMessage : state.errorMessage
    },
  },
  actions: {
    async loadLogs(force = false, silent = false) {
      if (!force && this.state === 'success') {
        return
      }
      if (!silent || this.state === 'idle') {
        this.state = 'loading'
      }
      try {
        const result = await apiClient.getOperationLogs()
        this.items = (result.logs ?? []).map(mapLogItem)
        this.state = 'success'
      } catch (error) {
        this.state = 'error'
        this.errorMessage = error instanceof Error ? error.message : '日志读取失败'
        throw error
      }
    },
    async loadFilteredLogs(force = false, silent = false) {
      if (!this.hasActiveFilters) {
        return
      }
      if (!force && this.filteredState === 'success') {
        return
      }
      if (!silent || this.filteredState === 'idle') {
        this.filteredState = 'loading'
      }
      try {
        const result = await apiClient.getOperationLogs({
          start: toRFC3339(this.filterStart),
          end: toRFC3339(this.filterEnd),
        })
        this.filteredItems = (result.logs ?? []).map(mapLogItem)
        this.filteredState = 'success'
        this.filteredErrorMessage = ''
      } catch (error) {
        this.filteredState = 'error'
        this.filteredErrorMessage = error instanceof Error ? error.message : '筛选日志读取失败'
        throw error
      }
    },
    async applyDateFilter(range: { start: string; end: string }) {
      this.filterStart = range.start
      this.filterEnd = range.end
      this.filteredItems = []
      this.filteredState = 'idle'
      this.filteredErrorMessage = ''
      if (!this.hasActiveFilters) {
        return
      }
      await this.loadFilteredLogs(true)
    },
    clearDateFilter() {
      this.filterStart = ''
      this.filterEnd = ''
      this.filteredItems = []
      this.filteredState = 'idle'
      this.filteredErrorMessage = ''
    },
    async refreshLogs() {
      await this.loadLogs(true, true)
      if (this.hasActiveFilters) {
        await this.loadFilteredLogs(true, true)
      }
    },
  },
})
