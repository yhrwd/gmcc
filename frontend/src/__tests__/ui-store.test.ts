import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import { useUiStore } from '@/stores/ui'
import { useLogsStore } from '@/stores/logs'

describe('ui store', () => {
  it('shares one login task across views', () => {
    setActivePinia(createPinia())
    const ui = useUiStore()
    ui.openLoginPanel('acc-main')
    expect(ui.loginTask.accountId).toBe('acc-main')
    expect(ui.loginTask.panelOpen).toBe(true)
  })

  it('keeps log filter state in the shared store', () => {
    setActivePinia(createPinia())
    const logs = useLogsStore()
    logs.filterStart = '2026-04-06T08:00'
    logs.filterEnd = '2026-04-06T10:00'
    expect(logs.hasActiveFilters).toBe(true)
  })
})
