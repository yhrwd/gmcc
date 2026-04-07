import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import { mapInstance } from '@/lib/mappers'
import { syncCoordinator } from '@/lib/sync'
import { useUiStore } from '@/stores/ui'
import type { CreateInstancePayload } from '@/types/api'
import type { LoadState, ViewInstance } from '@/types/view'

export type CreateInstanceForm = {
  id: string
  accountId: string
  serverAddress: string
  enabled: boolean
  autoStart: boolean
}

export function validateCreateInstance(payload: CreateInstanceForm) {
  if (!payload.id.trim()) return { valid: false, message: '实例 ID 不能为空' }
  if (!payload.accountId.trim()) return { valid: false, message: '请选择账号' }
  if (!payload.serverAddress.trim()) return { valid: false, message: '服务器地址不能为空' }
  if (!payload.enabled && payload.autoStart) return { valid: false, message: '禁用实例不能启用自动启动' }
  return { valid: true, message: '' }
}

export function getInstanceActions(status: 'pending' | 'starting' | 'running' | 'reconnecting' | 'stopped' | 'error' | 'unknown') {
  if (status === 'pending') return ['start']
  if (status === 'starting') return ['restart']
  if (status === 'running') return ['stop', 'restart']
  if (status === 'reconnecting') return ['stop', 'restart']
  if (status === 'stopped') return ['start']
  if (status === 'error') return ['start', 'restart']
  return ['restart']
}

export const useInstancesStore = defineStore('instances', {
  state: () => ({
    state: 'idle' as LoadState,
    errorMessage: '',
    submitting: false,
    actionBusyMap: {} as Record<string, 'start' | 'stop' | 'restart' | 'delete'>,
    deletingId: '',
    items: [] as ViewInstance[],
  }),
  actions: {
    async loadInstances(force = false, silent = false) {
      if (!force && this.state === 'success') {
        return
      }
      if (!silent || this.state === 'idle') {
        this.state = 'loading'
      }
      try {
        const result = await apiClient.getInstances()
        this.items = (result.instances ?? []).map(mapInstance)
        this.state = 'success'
      } catch (error) {
        this.state = 'error'
        this.errorMessage = error instanceof Error ? error.message : '实例列表加载失败'
      }
    },
    async createInstance(form: CreateInstanceForm) {
      const ui = useUiStore()
      const validation = validateCreateInstance(form)
      if (!validation.valid) {
        return validation
      }
      this.submitting = true
      try {
        const payload: CreateInstancePayload = {
          id: form.id.trim(),
          account_id: form.accountId,
          server_address: form.serverAddress.trim(),
          enabled: form.enabled,
          auto_start: form.autoStart,
        }
        await apiClient.createInstance(payload)
        this.items = [
          {
            id: payload.id,
            accountId: payload.account_id,
            serverAddress: payload.server_address,
            statusTone: form.autoStart ? 'starting' : 'pending',
            statusLabel: form.autoStart ? '启动中' : '待出勤',
            onlineDurationLabel: '0s',
            health: null,
            food: null,
            positionLabel: '',
          },
          ...this.items.filter((item) => item.id !== payload.id),
        ]
        this.state = 'success'
        syncCoordinator.requestNow(['instances', 'overview'], 'instance-created')
        ui.notify('success', '新的出勤小队已经加入基地')
        return { valid: true, message: '' }
      } catch (error) {
        return { valid: false, message: error instanceof Error ? error.message : '创建实例失败' }
      } finally {
        this.submitting = false
      }
    },
    async runAction(id: string, action: 'start' | 'stop' | 'restart') {
      const ui = useUiStore()
      this.actionBusyMap = {
        ...this.actionBusyMap,
        [id]: action,
      }
      this.items = this.items.map((item) => {
        if (item.id !== id) return item
        if (action === 'start' || action === 'restart') {
          return { ...item, statusTone: 'starting', statusLabel: '启动中' }
        }
        return { ...item, statusTone: 'stopped', statusLabel: '休息中' }
      })
      try {
        if (action === 'start') await apiClient.startInstance(id)
        if (action === 'stop') await apiClient.stopInstance(id)
        if (action === 'restart') await apiClient.restartInstance(id)
        ui.notify('success', `实例 ${id} 已${action === 'start' ? '启动' : action === 'stop' ? '停止' : '重启'}`)
        syncCoordinator.requestNow(['instances', 'overview'], `instance-${action === 'start' ? 'started' : action === 'stop' ? 'stopped' : 'restarted'}`)
      } catch (error) {
        ui.notify('error', error instanceof Error ? error.message : '实例操作失败')
        void this.loadInstances(true, true)
      } finally {
        const nextMap = { ...this.actionBusyMap }
        delete nextMap[id]
        this.actionBusyMap = nextMap
      }
    },
    async deleteInstance(id: string) {
      const ui = useUiStore()
      this.deletingId = id
      this.actionBusyMap = {
        ...this.actionBusyMap,
        [id]: 'delete',
      }
      try {
        await apiClient.deleteInstance(id)
        this.items = this.items.filter((item) => item.id !== id)
        this.state = 'success'
        syncCoordinator.requestNow(['instances', 'overview'], 'instance-deleted')
        ui.notify('success', `实例 ${id} 已移出基地`)
      } catch (error) {
        ui.notify('error', error instanceof Error ? error.message : '删除实例失败')
      } finally {
        this.deletingId = ''
        const nextMap = { ...this.actionBusyMap }
        delete nextMap[id]
        this.actionBusyMap = nextMap
      }
    },
  },
})
