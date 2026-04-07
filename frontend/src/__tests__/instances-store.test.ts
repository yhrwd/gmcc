import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { apiClient } from '@/api/client'
import { syncCoordinator } from '@/lib/sync'
import { getInstanceActions, useInstancesStore } from '@/stores/instances'

describe('instances store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('uses start action for pending status', () => {
    expect(getInstanceActions('pending')).toEqual(['start'])
  })

  it('uses start and restart for error status', () => {
    expect(getInstanceActions('error')).toEqual(['start', 'restart'])
  })

  it('uses cautious action set for unknown status', () => {
    expect(getInstanceActions('unknown')).toEqual(['restart'])
  })

  it('requests instances and overview sync after create succeeds', async () => {
    const store = useInstancesStore()
    const requestNow = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => {})
    vi.spyOn(apiClient, 'createInstance').mockResolvedValue({ success: true } as never)

    const result = await store.createInstance({
      id: 'instance-1',
      accountId: 'account-1',
      serverAddress: 'example.com',
      enabled: true,
      autoStart: false,
    })

    expect(result).toEqual({ valid: true, message: '' })
    expect(requestNow).toHaveBeenCalledWith(['instances', 'overview'], 'instance-created')
    expect(store.items[0]?.id).toBe('instance-1')
  })

  it('requests instances and overview sync after delete succeeds', async () => {
    const store = useInstancesStore()
    const requestNow = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => {})
    vi.spyOn(apiClient, 'deleteInstance').mockResolvedValue({ success: true } as never)
    store.items = [
      {
        id: 'instance-1',
        accountId: 'account-1',
        serverAddress: 'example.com',
        statusTone: 'pending',
        statusLabel: '待出勤',
        onlineDurationLabel: '0s',
        health: null,
        food: null,
        positionLabel: '',
      },
      {
        id: 'instance-2',
        accountId: 'account-2',
        serverAddress: 'example.org',
        statusTone: 'running',
        statusLabel: '执行中',
        onlineDurationLabel: '12s',
        health: null,
        food: null,
        positionLabel: '',
      },
    ]

    await store.deleteInstance('instance-1')

    expect(requestNow).toHaveBeenCalledWith(['instances', 'overview'], 'instance-deleted')
    expect(store.items).toHaveLength(1)
    expect(store.items[0]?.id).toBe('instance-2')
  })

  it('requests instances and overview sync with action-specific reason after runAction succeeds', async () => {
    const store = useInstancesStore()
    const requestNow = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => {})
    vi.spyOn(apiClient, 'restartInstance').mockResolvedValue({ success: true } as never)
    store.items = [
      {
        id: 'instance-1',
        accountId: 'account-1',
        serverAddress: 'example.com',
        statusTone: 'running',
        statusLabel: '执行中',
        onlineDurationLabel: '12s',
        health: null,
        food: null,
        positionLabel: '',
      },
    ]

    await store.runAction('instance-1', 'restart')

    expect(requestNow).toHaveBeenCalledWith(['instances', 'overview'], 'instance-restarted')
    expect(store.actionBusyMap['instance-1']).toBeUndefined()
    expect(store.items[0]?.statusTone).toBe('starting')
  })
})
