import { requestJson } from '@/lib/http'
import type {
  ApiAccountsResponse,
  ApiInitMicrosoftLoginResponse,
  ApiInstancesResponse,
  ApiOperationLogsResponse,
  ApiPollMicrosoftLoginResponse,
  ApiResourcesResponse,
  ApiResponse,
  ApiStatusResponse,
  CreateAccountPayload,
  CreateInstancePayload,
} from '@/types/api'

export const apiClient = {
  getStatus() {
    return requestJson<ApiStatusResponse>('/api/status')
  },
  getResources() {
    return requestJson<ApiResourcesResponse>('/api/resources')
  },
  getAccounts() {
    return requestJson<ApiAccountsResponse>('/api/accounts')
  },
  createAccount(payload: CreateAccountPayload) {
    return requestJson<ApiResponse>('/api/accounts', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },
  deleteAccount(id: string) {
    return requestJson<ApiResponse>(`/api/accounts/${id}`, {
      method: 'DELETE',
    })
  },
  initMicrosoftLogin(accountId: string) {
    return requestJson<ApiInitMicrosoftLoginResponse>('/api/auth/microsoft/init', {
      method: 'POST',
      body: JSON.stringify({ account_id: accountId }),
    })
  },
  pollMicrosoftLogin(accountId: string) {
    return requestJson<ApiPollMicrosoftLoginResponse>('/api/auth/microsoft/poll', {
      method: 'POST',
      body: JSON.stringify({ account_id: accountId }),
    })
  },
  getInstances() {
    return requestJson<ApiInstancesResponse>('/api/instances')
  },
  createInstance(payload: CreateInstancePayload) {
    return requestJson<ApiResponse>('/api/instances', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },
  deleteInstance(id: string) {
    return requestJson<ApiResponse>(`/api/instances/${id}`, { method: 'DELETE' })
  },
  startInstance(id: string) {
    return requestJson<ApiResponse>(`/api/instances/${id}/start`, { method: 'POST' })
  },
  stopInstance(id: string) {
    return requestJson<ApiResponse>(`/api/instances/${id}/stop`, { method: 'POST' })
  },
  restartInstance(id: string) {
    return requestJson<ApiResponse>(`/api/instances/${id}/restart`, { method: 'POST' })
  },
  getOperationLogs(params?: { start?: string; end?: string }) {
    const search = new URLSearchParams()
    if (params?.start) {
      search.set('start', params.start)
    }
    if (params?.end) {
      search.set('end', params.end)
    }
    const query = search.toString()
    return requestJson<ApiOperationLogsResponse>(query ? `/api/logs/operations?${query}` : '/api/logs/operations')
  },
}
