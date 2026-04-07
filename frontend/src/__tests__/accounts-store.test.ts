import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { apiClient } from '@/api/client'
import { syncCoordinator } from '@/lib/sync'
import { useAccountsStore, validateCreateAccount } from '@/stores/accounts'
import { useUiStore } from '@/stores/ui'

describe('accounts store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('blocks empty account id before submit', () => {
    expect(validateCreateAccount({ id: '   ', label: '', note: '' })).toEqual({ valid: false, message: '账号 ID 不能为空' })
  })

  it('returns cached data without forcing reload', async () => {
    const store = useAccountsStore()
    store.state = 'success'
    store.items = [{ id: 'a', label: 'A', note: '', enabled: true, authStatus: 'logged_in', authStatusTone: 'ready', authStatusLabel: '在线冒险者', hasToken: true }]

    await store.loadAccounts()

    expect(store.items).toHaveLength(1)
  })

  it('requests accounts sync after create succeeds', async () => {
    const store = useAccountsStore()
    const requestNow = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => {})
    vi.spyOn(apiClient, 'createAccount').mockResolvedValue({ success: true } as never)

    const result = await store.createAccount({ id: ' account-1 ', label: ' Label ', note: ' Note ' })

    expect(result).toEqual({ valid: true, message: '' })
    expect(requestNow).toHaveBeenCalledWith(['accounts'], 'account-created')
    expect(store.items[0]?.id).toBe('account-1')
  })

  it('requests accounts, instances, and overview sync after delete succeeds', async () => {
    const store = useAccountsStore()
    const requestNow = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => {})
    vi.spyOn(apiClient, 'deleteAccount').mockResolvedValue({ success: true } as never)
    store.items = [
      { id: 'account-1', label: 'A', note: '', enabled: true, authStatus: 'logged_in', authStatusTone: 'ready', authStatusLabel: '在线冒险者', hasToken: true },
      { id: 'account-2', label: 'B', note: '', enabled: true, authStatus: 'logged_in', authStatusTone: 'ready', authStatusLabel: '在线冒险者', hasToken: true },
    ]

    await store.deleteAccount('account-1')

    expect(requestNow).toHaveBeenCalledWith(['accounts', 'instances', 'overview'], 'account-deleted')
    expect(store.items).toHaveLength(1)
    expect(store.items[0]?.id).toBe('account-2')
  })

  it('requests accounts and overview sync after login poll succeeds', async () => {
    const store = useAccountsStore()
    const ui = useUiStore()
    const requestNow = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => {})
    vi.spyOn(apiClient, 'pollMicrosoftLogin').mockResolvedValue({ status: 'succeeded', message: 'ok' } as never)
    ui.setLoginTaskPayload({ accountId: 'account-1', taskStatus: 'polling' })

    await store.pollLogin()

    expect(requestNow).toHaveBeenCalledWith(['accounts', 'overview'], 'account-login-succeeded')
    expect(ui.loginTask.taskStatus).toBe('succeeded')
  })
})
