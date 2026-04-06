import { defineStore } from 'pinia'
import type { ActiveView, LoginTaskStatus, LoginTaskView, ToastMessage, ToastTone } from '@/types/view'

function createDefaultLoginTask(): LoginTaskView {
  return {
    accountId: '',
    taskStatus: 'idle',
    panelOpen: false,
    userCode: '',
    verificationUri: '',
    verificationUriComplete: '',
    expiresInSeconds: 0,
    expiresAt: null,
    intervalSeconds: 5,
    lastMessage: '',
    transientErrorCount: 0,
  }
}

export const useUiStore = defineStore('ui', {
  state: () => ({
    activeView: 'home' as ActiveView,
    toasts: [] as ToastMessage[],
    loginTask: createDefaultLoginTask(),
  }),
  actions: {
    setActiveView(view: ActiveView) {
      this.activeView = view
    },
    openLoginPanel(accountId: string) {
      if (this.loginTask.accountId && this.loginTask.accountId !== accountId && this.loginTask.taskStatus === 'polling') {
        this.loginTask.taskStatus = 'replaced'
        this.loginTask.lastMessage = '上一轮登录任务已被新的会话替换'
      }
      this.loginTask.accountId = accountId
      this.loginTask.panelOpen = true
      if (this.loginTask.taskStatus === 'idle' || this.loginTask.taskStatus === 'replaced') {
        this.loginTask.lastMessage = ''
      }
    },
    closeLoginPanel() {
      this.loginTask.panelOpen = false
    },
    setLoginTaskStatus(status: LoginTaskStatus, message = '') {
      this.loginTask.taskStatus = status
      this.loginTask.lastMessage = message
    },
    setLoginTaskPayload(payload: Partial<LoginTaskView>) {
      this.loginTask = {
        ...this.loginTask,
        ...payload,
      }
    },
    resetLoginTask(accountId = '') {
      this.loginTask = {
        ...createDefaultLoginTask(),
        accountId,
        panelOpen: Boolean(accountId),
      }
    },
    notify(tone: ToastTone, message: string) {
      const id = `${Date.now()}-${Math.random().toString(16).slice(2)}`
      this.toasts.push({ id, tone, message })
      window.setTimeout(() => {
        this.dismissToast(id)
      }, 3600)
    },
    dismissToast(id: string) {
      this.toasts = this.toasts.filter((toast) => toast.id !== id)
    },
  },
})
