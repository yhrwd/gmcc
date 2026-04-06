import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import { mapAccount } from '@/lib/mappers'
import { useUiStore } from '@/stores/ui'
import type { CreateAccountPayload } from '@/types/api'
import type { LoadState, ViewAccount } from '@/types/view'

export function validateCreateAccount(payload: CreateAccountPayload) {
  if (!payload.id.trim()) {
    return { valid: false, message: '账号 ID 不能为空' }
  }
  return { valid: true, message: '' }
}

export const useAccountsStore = defineStore('accounts', {
  state: () => ({
    state: 'idle' as LoadState,
    errorMessage: '',
    submitting: false,
    deletingId: '',
    loginStartingId: '',
    items: [] as ViewAccount[],
  }),
  getters: {
    readyAccounts(state) {
      return state.items.filter((item) => item.authStatusTone === 'ready')
    },
  },
  actions: {
    async loadAccounts(force = false, silent = false) {
      if (!force && this.state === 'success') {
        return
      }
      if (!silent || this.state === 'idle') {
        this.state = 'loading'
      }
      try {
        const result = await apiClient.getAccounts()
        this.items = (result.accounts ?? []).map(mapAccount)
        this.state = 'success'
      } catch (error) {
        this.state = 'error'
        this.errorMessage = error instanceof Error ? error.message : '账号列表加载失败'
      }
    },
    async createAccount(payload: CreateAccountPayload) {
      const ui = useUiStore()
      const validation = validateCreateAccount(payload)
      if (!validation.valid) {
        return validation
      }

      this.submitting = true
      try {
        await apiClient.createAccount({
          id: payload.id.trim(),
          label: payload.label?.trim() || '',
          note: payload.note?.trim() || '',
        })
        this.items = [
          {
            id: payload.id.trim(),
            label: payload.label?.trim() || payload.id.trim(),
            note: payload.note?.trim() || '',
            enabled: true,
            authStatus: 'not_logged_in',
            authStatusTone: 'pending',
            authStatusLabel: '待唤醒',
            hasToken: false,
          },
          ...this.items.filter((item) => item.id !== payload.id.trim()),
        ]
        this.state = 'success'
        void this.loadAccounts(true, true)
        ui.notify('success', '新的冒险者已经入驻基地')
        return { valid: true, message: '' }
      } catch (error) {
        return {
          valid: false,
          message: error instanceof Error ? error.message : '创建账号失败',
        }
      } finally {
        this.submitting = false
      }
    },
    async startLogin(accountId: string) {
      const ui = useUiStore()
      this.loginStartingId = accountId
      ui.openLoginPanel(accountId)
      ui.setLoginTaskPayload({
        accountId,
        panelOpen: true,
        taskStatus: 'initializing',
        userCode: '',
        verificationUri: '',
        verificationUriComplete: '',
        expiresInSeconds: 0,
        expiresAt: null,
        intervalSeconds: 5,
        lastMessage: '正在初始化登录流程...',
        transientErrorCount: 0,
      })

      try {
        const result = await apiClient.initMicrosoftLogin(accountId)
        const expiresInSeconds = result.expires_in || 0
        ui.setLoginTaskPayload({
          taskStatus: 'polling',
          userCode: result.user_code || '',
          verificationUri: result.verification_uri || '',
          verificationUriComplete: result.verification_uri_complete || '',
          expiresInSeconds,
          expiresAt: expiresInSeconds > 0 ? Date.now() + expiresInSeconds * 1000 : null,
          intervalSeconds: result.interval || 5,
          lastMessage: result.message || '请在新窗口完成设备码验证',
        })
      } catch (error) {
        ui.setLoginTaskPayload({
          taskStatus: 'failed',
          lastMessage: error instanceof Error ? error.message : '登录初始化失败',
        })
      } finally {
        this.loginStartingId = ''
      }
    },
    async pollLogin() {
      const ui = useUiStore()
      const accountId = ui.loginTask.accountId
      if (!accountId || ui.loginTask.taskStatus !== 'polling') {
        return
      }

      try {
        const result = await apiClient.pollMicrosoftLogin(accountId)
        const status = result.status || 'pending'
        if (status === 'pending') {
          ui.setLoginTaskPayload({
            lastMessage: result.message || '等待验证完成中...',
          })
          return
        }

        if (status === 'succeeded') {
          ui.setLoginTaskPayload({
            taskStatus: 'succeeded',
            lastMessage: result.message || '登录成功',
          })
          ui.notify('success', '账号已经成功醒来')
          void this.loadAccounts(true, true)
          return
        }

        if (status === 'expired') {
          ui.setLoginTaskPayload({ taskStatus: 'expired', lastMessage: result.message || '设备码已过期' })
          return
        }

        if (status === 'cancelled') {
          ui.setLoginTaskPayload({ taskStatus: 'failed', lastMessage: result.message || '登录流程已取消' })
          return
        }

        if (status === 'error') {
          ui.setLoginTaskPayload({ taskStatus: 'failed', lastMessage: result.message || '登录状态同步失败' })
          return
        }

        ui.setLoginTaskPayload({ taskStatus: 'failed', lastMessage: result.message || '登录失败' })
      } catch (error) {
        const nextCount = ui.loginTask.transientErrorCount + 1
        if (nextCount <= 3) {
          ui.setLoginTaskPayload({
            transientErrorCount: nextCount,
            lastMessage: '连接有点不稳，正在继续尝试',
          })
          return
        }

        ui.setLoginTaskPayload({
          taskStatus: 'failed',
          transientErrorCount: nextCount,
          lastMessage: error instanceof Error ? error.message : '轮询失败',
        })
      }
    },
    async deleteAccount(id: string) {
      const ui = useUiStore()
      this.deletingId = id
      try {
        await apiClient.deleteAccount(id)
        this.items = this.items.filter((item) => item.id !== id)
        this.state = 'success'
        void this.loadAccounts(true, true)
        ui.notify('success', `账号 ${id} 已从基地名册移除`)
      } catch (error) {
        ui.notify('error', error instanceof Error ? error.message : '删除账号失败')
      } finally {
        this.deletingId = ''
      }
    },
  },
})
