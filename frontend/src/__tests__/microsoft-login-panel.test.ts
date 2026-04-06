import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import { useUiStore } from '@/stores/ui'
import { useAccountsStore } from '@/stores/accounts'
import { apiClient } from '@/api/client'

describe('microsoft login task', () => {
  it('marks previous login task as replaced when a new login starts', () => {
    setActivePinia(createPinia())
    const ui = useUiStore()
    ui.openLoginPanel('acc-old')
    ui.setLoginTaskStatus('polling')
    ui.openLoginPanel('acc-new')
    expect(ui.loginTask.accountId).toBe('acc-new')
  })

  it('provides default countdown and verification fields', () => {
    setActivePinia(createPinia())
    const ui = useUiStore()
    expect(ui.loginTask.verificationUri).toBe('')
    expect(ui.loginTask.expiresInSeconds).toBe(0)
    expect(ui.loginTask.expiresAt).toBeNull()
  })

  it('clears stale login payload before a new login request resolves', async () => {
    setActivePinia(createPinia())
    const ui = useUiStore()
    const accounts = useAccountsStore()
    const originalInitMicrosoftLogin = apiClient.initMicrosoftLogin
    apiClient.initMicrosoftLogin = async () => ({ success: true })

    ui.setLoginTaskPayload({
      accountId: 'acc-old',
      userCode: 'ABCD-EFGH',
      verificationUri: 'https://example.test/device',
      verificationUriComplete: 'https://example.test/device?code=ABCD-EFGH',
      expiresInSeconds: 900,
      expiresAt: Date.now() + 900000,
      taskStatus: 'polling',
      panelOpen: true,
    })

    await accounts.startLogin('acc-new')

    expect(ui.loginTask.accountId).toBe('acc-new')
    expect(ui.loginTask.userCode).toBe('')
    expect(ui.loginTask.verificationUri).toBe('')
    expect(ui.loginTask.verificationUriComplete).toBe('')
    expect(ui.loginTask.expiresAt).toBeNull()
    apiClient.initMicrosoftLogin = originalInitMicrosoftLogin
  })

  it('keeps cancelled status distinct from failed status', async () => {
    setActivePinia(createPinia())
    const ui = useUiStore()
    const accounts = useAccountsStore()
    ui.setLoginTaskPayload({ accountId: 'acc-main', taskStatus: 'polling', panelOpen: true })
    const originalPollMicrosoftLogin = apiClient.pollMicrosoftLogin
    apiClient.pollMicrosoftLogin = async () => ({ success: false, status: 'cancelled', message: 'Device login cancelled' })

    await accounts.pollLogin()

    expect(ui.loginTask.taskStatus).toBe('cancelled')
    apiClient.pollMicrosoftLogin = originalPollMicrosoftLogin
  })
})
