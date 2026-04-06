# Vue Frontend Console Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Vue 3 single-page management console for gmcc that supports dashboard, account management, Microsoft device-code login, instance management, operation logs, responsive layouts, and embedded frontend delivery.

**Architecture:** Keep the frontend in `frontend/` as a Vite + Vue + TypeScript app with a thin app shell, feature panels, Pinia stores, a small API client, and a refresh coordinator composable. Reuse the existing Go embedding and packager flow so the built SPA can be served from the existing backend without changing `/api` semantics.

**Tech Stack:** Vue 3, TypeScript, Vite, Pinia, Tailwind CSS, Vitest, Vue Test Utils, Go, Gin

---

## File Map

- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tsconfig.app.json`
- Create: `frontend/tsconfig.node.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/postcss.config.js`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/style.css`
- Create: `frontend/src/env.d.ts`
- Create: `frontend/src/lib/http.ts`
- Create: `frontend/src/lib/errors.ts`
- Create: `frontend/src/lib/format.ts`
- Create: `frontend/src/types/api.ts`
- Create: `frontend/src/types/domain.ts`
- Create: `frontend/src/api/client.ts`
- Create: `frontend/src/stores/status.ts`
- Create: `frontend/src/stores/accounts.ts`
- Create: `frontend/src/stores/instances.ts`
- Create: `frontend/src/stores/logs.ts`
- Create: `frontend/src/stores/ui.ts`
- Create: `frontend/src/composables/useRefreshCoordinator.ts`
- Create: `frontend/src/components/layout/AppShell.vue`
- Create: `frontend/src/components/layout/AppHeader.vue`
- Create: `frontend/src/components/layout/DesktopSidebar.vue`
- Create: `frontend/src/components/layout/MobileTabbar.vue`
- Create: `frontend/src/components/panels/OverviewPanel.vue`
- Create: `frontend/src/components/panels/AccountsPanel.vue`
- Create: `frontend/src/components/panels/InstancesPanel.vue`
- Create: `frontend/src/components/panels/LogsPanel.vue`
- Create: `frontend/src/components/panels/ContextPanel.vue`
- Create: `frontend/src/components/accounts/AccountList.vue`
- Create: `frontend/src/components/accounts/AccountCard.vue`
- Create: `frontend/src/components/accounts/AccountCreateDialog.vue`
- Create: `frontend/src/components/accounts/MicrosoftLoginDialog.vue`
- Create: `frontend/src/components/instances/InstanceList.vue`
- Create: `frontend/src/components/instances/InstanceCard.vue`
- Create: `frontend/src/components/instances/InstanceCreateDialog.vue`
- Create: `frontend/src/components/logs/LogList.vue`
- Create: `frontend/src/components/shared/StatCard.vue`
- Create: `frontend/src/components/shared/ResourceCard.vue`
- Create: `frontend/src/components/shared/StatusBadge.vue`
- Create: `frontend/src/components/shared/ToastHost.vue`
- Create: `frontend/src/components/shared/EmptyState.vue`
- Create: `frontend/src/components/shared/LoadingBlock.vue`
- Create: `frontend/src/assets/fonts/.gitkeep`
- Create: `frontend/src/__tests__/api-client.test.ts`
- Create: `frontend/src/__tests__/accounts-store.test.ts`
- Create: `frontend/src/__tests__/instances-store.test.ts`
- Create: `frontend/src/__tests__/refresh-coordinator.test.ts`
- Create: `frontend/src/__tests__/instance-create-dialog.test.ts`
- Create: `frontend/src/__tests__/microsoft-login-dialog.test.ts`
- Modify: `frontend/.gitkeep`
- Modify: `tools/packager/main_test.go`
- Modify: `README.md`

## Implementation Notes

- Follow Vue 3 Composition API with `<script setup lang="ts">` in all SFCs.
- Keep source state minimal and derive display state with `computed`.
- Treat account selection and instance selection as mutually exclusive UI state.
- First version uses one active Microsoft device-code session at a time.
- First version uses list endpoints as the source for account and instance detail panels; do not depend on `/api/accounts/:id` or `/api/instances/:id`.
- API requests must use same-origin relative paths under `/api`.
- Vite build output must remain compatible with the current packager whitelist: root `index.html`, allowed root icons/manifests, and `assets/` files only.
- If dependency install or package download hits network failures, stop and report it to the user immediately.

## Component Map

- `App.vue`: compose app shell, initialize stores, and wire the refresh coordinator.
- `AppShell.vue`: responsive page frame with desktop and mobile layout slots.
- `OverviewPanel.vue`: show global stats, resources, cluster state, and recent log summary only when nothing is selected.
- `ContextPanel.vue`: show selected account or selected instance context and related actions.
- `AccountsPanel.vue`: account list, create account dialog, login dialog entry, and account selection.
- `InstancesPanel.vue`: instance list, create instance dialog, and instance lifecycle actions.
- `LogsPanel.vue`: full log list and manual refresh entry.
- `MicrosoftLoginDialog.vue`: one active device-code session UI and polling state.
- `InstanceCreateDialog.vue`: create-instance form with required account binding and guard rails.

### Task 1: Scaffold the Vue app and build toolchain

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tsconfig.app.json`
- Create: `frontend/tsconfig.node.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/postcss.config.js`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/style.css`
- Create: `frontend/src/env.d.ts`
- Modify: `frontend/.gitkeep`

- [ ] **Step 1: Replace the placeholder frontend marker with the real app skeleton**

Delete `frontend/.gitkeep` and create the base project files with these contents.

```json
{
  "name": "gmcc-frontend",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "preview": "vite preview",
    "test": "vitest run"
  },
  "dependencies": {
    "pinia": "^3.0.3",
    "vue": "^3.5.13"
  },
  "devDependencies": {
    "@tailwindcss/postcss": "^4.1.4",
    "@vitejs/plugin-vue": "^5.2.3",
    "@vue/test-utils": "^2.4.6",
    "autoprefixer": "^10.4.21",
    "jsdom": "^26.1.0",
    "postcss": "^8.5.3",
    "tailwindcss": "^4.1.4",
    "typescript": "^5.8.3",
    "vite": "^6.2.4",
    "vitest": "^3.1.1",
    "vue-tsc": "^2.2.8"
  }
}
```

```json
{
  "files": [],
  "references": [
    { "path": "./tsconfig.app.json" },
    { "path": "./tsconfig.node.json" }
  ]
}
```

```json
{
  "extends": "@vue/tsconfig/tsconfig.dom.json",
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.app.tsbuildinfo",
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src/**/*.ts", "src/**/*.tsx", "src/**/*.vue"]
}
```

```json
{
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.node.tsbuildinfo",
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "types": ["node"]
  },
  "include": ["vite.config.ts"]
}
```

```ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  base: '/',
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
  },
})
```

```js
export default {
  plugins: {
    '@tailwindcss/postcss': {},
    autoprefixer: {},
  },
}
```

```html
<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>gmcc Console</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
```

```ts
/// <reference types="vite/client" />
```

```ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

const app = createApp(App)

app.use(createPinia())
app.mount('#app')
```

```vue
<script setup lang="ts">
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-100">
    <div class="mx-auto flex min-h-screen max-w-[1600px] items-center justify-center px-6 py-10">
      <div class="rounded-3xl border border-slate-800 bg-slate-900/80 px-8 py-10 shadow-2xl shadow-cyan-950/20 backdrop-blur">
        <p class="text-sm uppercase tracking-[0.3em] text-cyan-300/80">gmcc</p>
        <h1 class="mt-3 text-3xl font-semibold text-white">Vue Console scaffolded</h1>
        <p class="mt-3 max-w-xl text-sm leading-7 text-slate-400">
          前端工程已初始化，下一步会接入 API client、stores 和控制台布局。
        </p>
      </div>
    </div>
  </div>
</template>
```

```css
@import 'tailwindcss';

:root {
  color-scheme: dark;
  font-family: "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
  background:
    radial-gradient(circle at top, rgba(34, 211, 238, 0.14), transparent 35%),
    linear-gradient(180deg, #020617 0%, #0f172a 55%, #020617 100%);
}

body {
  margin: 0;
  min-width: 320px;
}

#app {
  min-height: 100vh;
}
```

- [ ] **Step 2: Install dependencies**

Run: `npm install`
Expected: install completes successfully and creates `frontend/node_modules` plus `package-lock.json`

- [ ] **Step 3: Verify the empty scaffold builds**

Run: `npm run build`
Expected: PASS and Vite writes files under `frontend/dist`

- [ ] **Step 4: Commit the scaffold**

```bash
git add frontend
git commit -m "feat: scaffold vue frontend workspace"
```

### Task 2: Add typed domain models and the API client

**Files:**
- Create: `frontend/src/types/api.ts`
- Create: `frontend/src/types/domain.ts`
- Create: `frontend/src/lib/errors.ts`
- Create: `frontend/src/lib/http.ts`
- Create: `frontend/src/lib/format.ts`
- Create: `frontend/src/api/client.ts`
- Create: `frontend/src/__tests__/api-client.test.ts`

- [ ] **Step 1: Write failing API client tests for normalization and error extraction**

Create `frontend/src/__tests__/api-client.test.ts`.

```ts
import { describe, expect, it, vi } from 'vitest'
import { createApiClient } from '@/api/client'

describe('api client', () => {
  it('maps operation logs into the frontend model', async () => {
    const fetcher = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({
        logs: [
          {
            id: 'log-1',
            timestamp: '2026-04-05T10:00:00Z',
            action: 'instance_start',
            target_instance_id: 'bot-1',
            target_account_id: 'acc-main',
            success: true,
          },
        ],
      }), { status: 200 })
    )

    const client = createApiClient(fetcher)
    const result = await client.listOperationLogs()

    expect(result[0]).toEqual({
      id: 'log-1',
      timestamp: '2026-04-05T10:00:00Z',
      action: 'instance_start',
      targetInstanceId: 'bot-1',
      targetAccountId: 'acc-main',
      success: true,
      clientIp: undefined,
      userAgent: undefined,
    })
  })

  it('extracts message-based api errors', async () => {
    const fetcher = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ message: 'Instance already exists' }), { status: 400 })
    )

    const client = createApiClient(fetcher)

    await expect(client.createInstance({
      id: 'bot-1',
      account_id: 'acc-main',
      server_address: 'mc.example.com',
      enabled: true,
      auto_start: false,
    })).rejects.toMatchObject({
      message: 'Instance already exists',
      statusCode: 400,
      retryable: false,
    })
  })
})
```

- [ ] **Step 2: Run the API client test to verify it fails**

Run: `npm run test -- src/__tests__/api-client.test.ts`
Expected: FAIL with module not found errors for `@/api/client`

- [ ] **Step 3: Add the type definitions and client implementation**

Create the typed models and minimal client.

```ts
export type ClusterStatus = 'running' | 'partial' | 'stopped'

export type AccountApiItem = {
  id: string
  label?: string
  note?: string
  enabled?: boolean
  auth_status?: string
  has_token?: boolean
  player_id?: string
}

export type InstanceApiItem = {
  id: string
  account_id: string
  player_id?: string
  server_address: string
  status?: string
  online_duration?: string
  last_seen?: string
  has_token?: boolean
  health?: number
  food?: number
  position?: { x: number; y: number; z: number }
}

export type StatusResponse = {
  cluster_status: ClusterStatus
  total_instances: number
  running_instances: number
  uptime: number
}

export type ResourcesResponse = {
  cpu_percent: number
  memory: {
    total_bytes: number
    used_bytes: number
    available_bytes: number
    used_percent: number
  }
  collected_at: string
}

export type AccountsResponse = { accounts: AccountApiItem[] }
export type InstancesResponse = { instances: InstanceApiItem[] }

export type LogsResponse = {
  logs: Array<{
    id: string
    timestamp: string
    action: string
    target_instance_id: string
    target_account_id: string
    success: boolean
    client_ip?: string
    user_agent?: string
  }>
}

export type MicrosoftInitResponse = {
  success: boolean
  user_code: string
  verification_uri: string
  verification_uri_complete: string
  expires_in: number
  interval: number
  account_id: string
}

export type MicrosoftPollResponse = {
  success: boolean
  status: string
  message: string
  account_id: string
}
```

```ts
export type NavMode = 'overview' | 'accounts' | 'instances' | 'logs'

export type OperationLogItem = {
  id: string
  timestamp: string
  action: string
  targetInstanceId: string
  targetAccountId: string
  success: boolean
  clientIp?: string
  userAgent?: string
}

export type CreateAccountInput = {
  id: string
  label: string
  note: string
}

export type CreateInstanceInput = {
  id: string
  account_id: string
  server_address: string
  enabled: boolean
  auto_start: boolean
}

export type MicrosoftSession = {
  accountId: string
  userCode: string
  verificationUri: string
  verificationUriComplete: string
  expiresAt: number
  pollIntervalSeconds: number
  status: string
  message: string
}
```

```ts
export class ApiError extends Error {
  statusCode: number
  retryable: boolean

  constructor(message: string, statusCode: number, retryable: boolean) {
    super(message)
    this.name = 'ApiError'
    this.statusCode = statusCode
    this.retryable = retryable
  }
}

export function toApiError(statusCode: number, payload?: Record<string, unknown>) {
  const message =
    typeof payload?.message === 'string'
      ? payload.message
      : typeof payload?.error === 'string'
        ? payload.error
        : statusCode >= 500
          ? '服务暂时不可用'
          : '请求失败'

  return new ApiError(message, statusCode, statusCode >= 500 || statusCode === 429)
}
```

```ts
import { toApiError } from '@/lib/errors'

export async function httpJson<T>(
  input: RequestInfo | URL,
  init?: RequestInit,
  fetcher: typeof fetch = fetch,
): Promise<T> {
  const response = await fetcher(input, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers ?? {}),
    },
    ...init,
  })

  let payload: Record<string, unknown> | undefined

  try {
    payload = (await response.json()) as Record<string, unknown>
  } catch {
    payload = undefined
  }

  if (!response.ok) {
    throw toApiError(response.status, payload)
  }

  return payload as T
}
```

```ts
export function formatPercent(value: number) {
  return `${value.toFixed(1)}%`
}

export function formatBytes(value: number) {
  if (value < 1024) return `${value} B`
  const units = ['KB', 'MB', 'GB', 'TB']
  let current = value / 1024
  let idx = 0
  while (current >= 1024 && idx < units.length - 1) {
    current /= 1024
    idx += 1
  }
  return `${current.toFixed(1)} ${units[idx]}`
}
```

```ts
import { httpJson } from '@/lib/http'
import type {
  AccountsResponse,
  InstancesResponse,
  LogsResponse,
  MicrosoftInitResponse,
  MicrosoftPollResponse,
  ResourcesResponse,
  StatusResponse,
} from '@/types/api'
import type { CreateAccountInput, CreateInstanceInput, OperationLogItem } from '@/types/domain'

export function createApiClient(fetcher: typeof fetch = fetch) {
  return {
    getStatus() {
      return httpJson<StatusResponse>('/api/status', undefined, fetcher)
    },
    getResources() {
      return httpJson<ResourcesResponse>('/api/resources', undefined, fetcher)
    },
    listAccounts() {
      return httpJson<AccountsResponse>('/api/accounts', undefined, fetcher)
    },
    createAccount(payload: CreateAccountInput) {
      return httpJson<{ success: boolean; message: string }>('/api/accounts', {
        method: 'POST',
        body: JSON.stringify(payload),
      }, fetcher)
    },
    deleteAccount(accountId: string) {
      return httpJson<{ success: boolean; message: string }>(`/api/accounts/${accountId}`, {
        method: 'DELETE',
      }, fetcher)
    },
    initMicrosoftLogin(accountId: string) {
      return httpJson<MicrosoftInitResponse>('/api/auth/microsoft/init', {
        method: 'POST',
        body: JSON.stringify({ account_id: accountId }),
      }, fetcher)
    },
    pollMicrosoftLogin(accountId: string) {
      return httpJson<MicrosoftPollResponse>('/api/auth/microsoft/poll', {
        method: 'POST',
        body: JSON.stringify({ account_id: accountId }),
      }, fetcher)
    },
    listInstances() {
      return httpJson<InstancesResponse>('/api/instances', undefined, fetcher)
    },
    createInstance(payload: CreateInstanceInput) {
      return httpJson<{ success: boolean; message: string }>('/api/instances', {
        method: 'POST',
        body: JSON.stringify(payload),
      }, fetcher)
    },
    startInstance(instanceId: string) {
      return httpJson<{ success: boolean; message: string }>(`/api/instances/${instanceId}/start`, {
        method: 'POST',
      }, fetcher)
    },
    stopInstance(instanceId: string) {
      return httpJson<{ success: boolean; message: string }>(`/api/instances/${instanceId}/stop`, {
        method: 'POST',
      }, fetcher)
    },
    restartInstance(instanceId: string) {
      return httpJson<{ success: boolean; message: string }>(`/api/instances/${instanceId}/restart`, {
        method: 'POST',
      }, fetcher)
    },
    deleteInstance(instanceId: string) {
      return httpJson<{ success: boolean; message: string }>(`/api/instances/${instanceId}`, {
        method: 'DELETE',
      }, fetcher)
    },
    async listOperationLogs() {
      const payload = await httpJson<LogsResponse>('/api/logs/operations', undefined, fetcher)
      return payload.logs.map<OperationLogItem>((item) => ({
        id: item.id,
        timestamp: item.timestamp,
        action: item.action,
        targetInstanceId: item.target_instance_id,
        targetAccountId: item.target_account_id,
        success: item.success,
        clientIp: item.client_ip,
        userAgent: item.user_agent,
      }))
    },
  }
}

export const apiClient = createApiClient()
```

- [ ] **Step 4: Run the API client test again**

Run: `npm run test -- src/__tests__/api-client.test.ts`
Expected: PASS

- [ ] **Step 5: Commit the typed client layer**

```bash
git add frontend/src/types frontend/src/lib frontend/src/api frontend/src/__tests__/api-client.test.ts
git commit -m "feat: add typed frontend api client"
```

### Task 3: Add Pinia stores and derived stats

**Files:**
- Create: `frontend/src/stores/status.ts`
- Create: `frontend/src/stores/accounts.ts`
- Create: `frontend/src/stores/instances.ts`
- Create: `frontend/src/stores/logs.ts`
- Create: `frontend/src/stores/ui.ts`
- Create: `frontend/src/__tests__/accounts-store.test.ts`
- Create: `frontend/src/__tests__/instances-store.test.ts`

- [ ] **Step 1: Write failing store tests for account eligibility and instance counts**

Create two tests.

```ts
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useAccountsStore } from '@/stores/accounts'

describe('accounts store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('treats only enabled logged-in tokenized accounts as instance-eligible', () => {
    const store = useAccountsStore()
    store.items = [
      { id: 'a', enabled: true, auth_status: 'logged_in', has_token: true },
      { id: 'b', enabled: true, auth_status: 'pending', has_token: true },
      { id: 'c', enabled: false, auth_status: 'logged_in', has_token: true },
    ]

    expect(store.loggedInCount).toBe(2)
    expect(store.canCreateInstance('a')).toBe(true)
    expect(store.canCreateInstance('b')).toBe(false)
    expect(store.canCreateInstance('c')).toBe(false)
  })
})
```

```ts
import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useInstancesStore } from '@/stores/instances'

describe('instances store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('derives running and non-running counts from status', () => {
    const store = useInstancesStore()
    store.items = [
      { id: 'i1', account_id: 'a', server_address: 'one', status: 'running' },
      { id: 'i2', account_id: 'a', server_address: 'two', status: 'stopped' },
    ]

    expect(store.totalCount).toBe(2)
    expect(store.runningCount).toBe(1)
    expect(store.nonRunningCount).toBe(1)
  })
})
```

- [ ] **Step 2: Run the store tests to verify they fail**

Run: `npm run test -- src/__tests__/accounts-store.test.ts src/__tests__/instances-store.test.ts`
Expected: FAIL with missing store modules

- [ ] **Step 3: Implement the stores with explicit actions and computed getters**

Create these minimal stores.

```ts
import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { ResourcesResponse, StatusResponse } from '@/types/api'

export const useStatusStore = defineStore('status', () => {
  const status = shallowRef<StatusResponse | null>(null)
  const resources = shallowRef<ResourcesResponse | null>(null)
  const loading = shallowRef(false)
  const error = shallowRef('')
  const lastRefreshedAt = shallowRef(0)

  const clusterStatus = computed(() => status.value?.cluster_status ?? 'stopped')
  const cpuPercent = computed(() => resources.value?.cpu_percent ?? 0)
  const memorySummary = computed(() => resources.value?.memory ?? null)

  async function refreshStatus() {
    loading.value = true
    error.value = ''
    try {
      status.value = await apiClient.getStatus()
      lastRefreshedAt.value = Date.now()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '状态刷新失败'
    } finally {
      loading.value = false
    }
  }

  async function refreshResources() {
    try {
      resources.value = await apiClient.getResources()
      lastRefreshedAt.value = Date.now()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '资源刷新失败'
    }
  }

  return { status, resources, loading, error, lastRefreshedAt, clusterStatus, cpuPercent, memorySummary, refreshStatus, refreshResources }
})
```

```ts
import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { AccountApiItem } from '@/types/api'
import type { CreateAccountInput, MicrosoftSession } from '@/types/domain'

export const useAccountsStore = defineStore('accounts', () => {
  const items = shallowRef<AccountApiItem[]>([])
  const loading = shallowRef(false)
  const error = shallowRef('')
  const actionLoadingIds = shallowRef<string[]>([])
  const loginSession = shallowRef<MicrosoftSession | null>(null)

  const enabledCount = computed(() => items.value.filter((item) => item.enabled === true).length)
  const loggedInCount = computed(() => items.value.filter((item) => item.auth_status === 'logged_in').length)

  function canCreateInstance(accountId: string) {
    const account = items.value.find((item) => item.id === accountId)
    if (!account) return false
    return account.enabled === true && account.auth_status === 'logged_in' && account.has_token === true
  }

  async function refreshAccounts() {
    loading.value = true
    error.value = ''
    try {
      const payload = await apiClient.listAccounts()
      items.value = payload.accounts
    } catch (err) {
      error.value = err instanceof Error ? err.message : '账号刷新失败'
    } finally {
      loading.value = false
    }
  }

  async function createAccount(payload: CreateAccountInput) {
    await apiClient.createAccount(payload)
    await refreshAccounts()
  }

  async function deleteAccount(accountId: string) {
    await apiClient.deleteAccount(accountId)
    await refreshAccounts()
  }

  async function startMicrosoftLogin(accountId: string) {
    const response = await apiClient.initMicrosoftLogin(accountId)
    loginSession.value = {
      accountId: response.account_id,
      userCode: response.user_code,
      verificationUri: response.verification_uri,
      verificationUriComplete: response.verification_uri_complete,
      expiresAt: Date.now() + response.expires_in * 1000,
      pollIntervalSeconds: response.interval,
      status: 'pending',
      message: 'Waiting for user authorization...',
    }
  }

  async function pollMicrosoftLogin() {
    if (!loginSession.value) return
    const response = await apiClient.pollMicrosoftLogin(loginSession.value.accountId)
    loginSession.value = {
      ...loginSession.value,
      status: response.status,
      message: response.message,
    }
    if (response.status === 'succeeded') {
      await refreshAccounts()
    }
  }

  function clearLoginSession() {
    loginSession.value = null
  }

  return {
    items,
    loading,
    error,
    actionLoadingIds,
    loginSession,
    enabledCount,
    loggedInCount,
    canCreateInstance,
    refreshAccounts,
    createAccount,
    deleteAccount,
    startMicrosoftLogin,
    pollMicrosoftLogin,
    clearLoginSession,
  }
})
```

```ts
import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { InstanceApiItem } from '@/types/api'
import type { CreateInstanceInput } from '@/types/domain'

export const useInstancesStore = defineStore('instances', () => {
  const items = shallowRef<InstanceApiItem[]>([])
  const loading = shallowRef(false)
  const error = shallowRef('')
  const actionLoadingIds = shallowRef<string[]>([])

  const totalCount = computed(() => items.value.length)
  const runningCount = computed(() => items.value.filter((item) => item.status === 'running').length)
  const nonRunningCount = computed(() => totalCount.value - runningCount.value)

  async function refreshInstances() {
    loading.value = true
    error.value = ''
    try {
      const payload = await apiClient.listInstances()
      items.value = payload.instances
    } catch (err) {
      error.value = err instanceof Error ? err.message : '实例刷新失败'
    } finally {
      loading.value = false
    }
  }

  async function createInstance(payload: CreateInstanceInput) {
    await apiClient.createInstance(payload)
    await refreshInstances()
  }

  async function startInstance(instanceId: string) {
    await apiClient.startInstance(instanceId)
    await refreshInstances()
  }

  async function stopInstance(instanceId: string) {
    await apiClient.stopInstance(instanceId)
    await refreshInstances()
  }

  async function restartInstance(instanceId: string) {
    await apiClient.restartInstance(instanceId)
    await refreshInstances()
  }

  async function deleteInstance(instanceId: string) {
    await apiClient.deleteInstance(instanceId)
    await refreshInstances()
  }

  return { items, loading, error, actionLoadingIds, totalCount, runningCount, nonRunningCount, refreshInstances, createInstance, startInstance, stopInstance, restartInstance, deleteInstance }
})
```

```ts
import { computed, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import { apiClient } from '@/api/client'
import type { OperationLogItem } from '@/types/domain'

export const useLogsStore = defineStore('logs', () => {
  const items = shallowRef<OperationLogItem[]>([])
  const loading = shallowRef(false)
  const error = shallowRef('')

  const recentSummary = computed(() => [...items.value].sort((a, b) => b.timestamp.localeCompare(a.timestamp)).slice(0, 5))

  async function refreshLogs() {
    loading.value = true
    error.value = ''
    try {
      items.value = await apiClient.listOperationLogs()
    } catch (err) {
      error.value = err instanceof Error ? err.message : '日志刷新失败'
    } finally {
      loading.value = false
    }
  }

  return { items, loading, error, recentSummary, refreshLogs }
})
```

```ts
import { shallowRef } from 'vue'
import { defineStore } from 'pinia'
import type { NavMode } from '@/types/domain'

export const useUiStore = defineStore('ui', () => {
  const mode = shallowRef<NavMode>('overview')
  const selectedAccountId = shallowRef('')
  const selectedInstanceId = shallowRef('')
  const toastMessage = shallowRef('')

  function selectAccount(accountId: string) {
    selectedAccountId.value = accountId
    selectedInstanceId.value = ''
  }

  function selectInstance(instanceId: string) {
    selectedInstanceId.value = instanceId
    selectedAccountId.value = ''
  }

  function clearSelection() {
    selectedAccountId.value = ''
    selectedInstanceId.value = ''
  }

  return { mode, selectedAccountId, selectedInstanceId, toastMessage, selectAccount, selectInstance, clearSelection }
})
```

- [ ] **Step 4: Run the store tests again**

Run: `npm run test -- src/__tests__/accounts-store.test.ts src/__tests__/instances-store.test.ts`
Expected: PASS

- [ ] **Step 5: Commit the store layer**

```bash
git add frontend/src/stores frontend/src/__tests__/accounts-store.test.ts frontend/src/__tests__/instances-store.test.ts
git commit -m "feat: add frontend state stores"
```

### Task 4: Add refresh coordination and lifecycle rules

**Files:**
- Create: `frontend/src/composables/useRefreshCoordinator.ts`
- Create: `frontend/src/__tests__/refresh-coordinator.test.ts`
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Write a failing test for mode-based log refresh coordination**

```ts
import { describe, expect, it, vi } from 'vitest'
import { useRefreshCoordinator } from '@/composables/useRefreshCoordinator'

describe('refresh coordinator', () => {
  it('refreshes logs only in overview and logs mode', async () => {
    const calls: string[] = []
    const coordinator = useRefreshCoordinator({
      refreshStatus: async () => calls.push('status'),
      refreshResources: async () => calls.push('resources'),
      refreshAccounts: async () => calls.push('accounts'),
      refreshInstances: async () => calls.push('instances'),
      refreshLogs: async () => calls.push('logs'),
      setTimeout: vi.fn((fn: () => void) => {
        fn()
        return 1 as unknown as ReturnType<typeof globalThis.setTimeout>
      }),
      clearTimeout: vi.fn(),
    })

    await coordinator.triggerManualRefresh('all', 'instances')

    expect(calls).not.toContain('logs')
  })
})
```

- [ ] **Step 2: Run the coordinator test to verify it fails**

Run: `npm run test -- src/__tests__/refresh-coordinator.test.ts`
Expected: FAIL with missing composable module

- [ ] **Step 3: Implement the coordinator and wire it into `App.vue`**

```ts
type RefreshScope = 'all' | 'status' | 'accounts' | 'instances' | 'logs'
type RefreshMode = 'overview' | 'accounts' | 'instances' | 'logs'

type RefreshCoordinatorDeps = {
  refreshStatus: () => Promise<void>
  refreshResources: () => Promise<void>
  refreshAccounts: () => Promise<void>
  refreshInstances: () => Promise<void>
  refreshLogs: () => Promise<void>
  setTimeout?: typeof globalThis.setTimeout
  clearTimeout?: typeof globalThis.clearTimeout
}

export function useRefreshCoordinator(deps: RefreshCoordinatorDeps) {
  const timers = new Set<ReturnType<typeof globalThis.setTimeout>>()
  const setTimer = deps.setTimeout ?? globalThis.setTimeout
  const clearTimer = deps.clearTimeout ?? globalThis.clearTimeout

  async function runBaseRefresh() {
    await Promise.all([
      deps.refreshStatus(),
      deps.refreshResources(),
      deps.refreshAccounts(),
      deps.refreshInstances(),
    ])
  }

  async function triggerManualRefresh(scope: RefreshScope, mode: RefreshMode) {
    if (scope === 'all') {
      await runBaseRefresh()
      if (mode === 'overview' || mode === 'logs') {
        await deps.refreshLogs()
      }
      return
    }

    if (scope === 'status') {
      await Promise.all([deps.refreshStatus(), deps.refreshResources()])
    }
    if (scope === 'accounts') await deps.refreshAccounts()
    if (scope === 'instances') await deps.refreshInstances()
    if (scope === 'logs') await deps.refreshLogs()
  }

  function schedule(mode: RefreshMode, visible = true) {
    const delay = visible ? 10_000 : 30_000
    const timer = setTimer(async () => {
      await runBaseRefresh()
      if (mode === 'overview' || mode === 'logs') {
        await deps.refreshLogs()
      }
    }, delay)
    timers.add(timer)
  }

  function stop() {
    for (const timer of timers) clearTimer(timer)
    timers.clear()
  }

  return {
    start(mode: RefreshMode) {
      stop()
      schedule(mode, true)
    },
    stop,
    setMode(mode: RefreshMode) {
      stop()
      schedule(mode, true)
    },
    setPageVisible(visible: boolean, mode: RefreshMode) {
      stop()
      schedule(mode, visible)
    },
    triggerManualRefresh,
  }
}
```

Replace `frontend/src/App.vue` with:

```vue
<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'
import AppShell from '@/components/layout/AppShell.vue'
import { useRefreshCoordinator } from '@/composables/useRefreshCoordinator'
import { useAccountsStore } from '@/stores/accounts'
import { useInstancesStore } from '@/stores/instances'
import { useLogsStore } from '@/stores/logs'
import { useStatusStore } from '@/stores/status'
import { useUiStore } from '@/stores/ui'

const statusStore = useStatusStore()
const accountsStore = useAccountsStore()
const instancesStore = useInstancesStore()
const logsStore = useLogsStore()
const uiStore = useUiStore()

const coordinator = useRefreshCoordinator({
  refreshStatus: statusStore.refreshStatus,
  refreshResources: statusStore.refreshResources,
  refreshAccounts: accountsStore.refreshAccounts,
  refreshInstances: instancesStore.refreshInstances,
  refreshLogs: logsStore.refreshLogs,
})

function handleVisibilityChange() {
  coordinator.setPageVisible(document.visibilityState === 'visible', uiStore.mode)
}

watch(() => uiStore.mode, (mode) => {
  coordinator.setMode(mode)
})

onMounted(async () => {
  await coordinator.triggerManualRefresh('all', uiStore.mode)
  coordinator.start(uiStore.mode)
  document.addEventListener('visibilitychange', handleVisibilityChange)
})

onUnmounted(() => {
  coordinator.stop()
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})
</script>

<template>
  <AppShell />
</template>
```

- [ ] **Step 4: Run the coordinator test again**

Run: `npm run test -- src/__tests__/refresh-coordinator.test.ts`
Expected: PASS

- [ ] **Step 5: Commit refresh coordination**

```bash
git add frontend/src/composables frontend/src/__tests__/refresh-coordinator.test.ts frontend/src/App.vue
git commit -m "feat: add frontend refresh coordinator"
```

### Task 5: Build the responsive shell and shared UI primitives

**Files:**
- Create: `frontend/src/components/layout/AppShell.vue`
- Create: `frontend/src/components/layout/AppHeader.vue`
- Create: `frontend/src/components/layout/DesktopSidebar.vue`
- Create: `frontend/src/components/layout/MobileTabbar.vue`
- Create: `frontend/src/components/shared/StatCard.vue`
- Create: `frontend/src/components/shared/ResourceCard.vue`
- Create: `frontend/src/components/shared/StatusBadge.vue`
- Create: `frontend/src/components/shared/ToastHost.vue`
- Create: `frontend/src/components/shared/EmptyState.vue`
- Create: `frontend/src/components/shared/LoadingBlock.vue`

- [ ] **Step 1: Create the layout shell and base components**

Implement the shell with desktop three-column structure and mobile tabs.

```vue
<script setup lang="ts">
const props = defineProps<{ title: string; value: string; hint: string }>()
</script>

<template>
  <article class="rounded-2xl border border-slate-800 bg-slate-900/80 p-5 shadow-lg shadow-slate-950/30">
    <p class="text-xs uppercase tracking-[0.25em] text-slate-400">{{ props.title }}</p>
    <p class="mt-3 text-3xl font-semibold text-white">{{ props.value }}</p>
    <p class="mt-2 text-sm text-slate-400">{{ props.hint }}</p>
  </article>
</template>
```

```vue
<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ status: string }>()

const tone = computed(() => {
  if (props.status === 'running' || props.status === 'logged_in') return 'bg-emerald-500/15 text-emerald-300 ring-emerald-400/30'
  if (props.status === 'partial' || props.status === 'pending') return 'bg-amber-500/15 text-amber-300 ring-amber-400/30'
  return 'bg-rose-500/15 text-rose-300 ring-rose-400/30'
})
</script>

<template>
  <span :class="['inline-flex rounded-full px-3 py-1 text-xs font-medium ring-1', tone]">
    {{ props.status }}
  </span>
</template>
```

```vue
<script setup lang="ts">
import { formatBytes, formatPercent } from '@/lib/format'

const props = defineProps<{
  cpuPercent: number
  usedBytes: number
  totalBytes: number
  memoryPercent: number
}>()
</script>

<template>
  <article class="rounded-2xl border border-cyan-900/50 bg-cyan-950/20 p-5">
    <div class="flex items-center justify-between">
      <h3 class="text-lg font-medium text-white">系统资源</h3>
      <span class="text-xs uppercase tracking-[0.25em] text-cyan-300/80">Live</span>
    </div>
    <div class="mt-5 grid gap-4 md:grid-cols-2">
      <div>
        <p class="text-sm text-slate-400">CPU</p>
        <p class="mt-2 text-2xl font-semibold text-white">{{ formatPercent(props.cpuPercent) }}</p>
      </div>
      <div>
        <p class="text-sm text-slate-400">Memory</p>
        <p class="mt-2 text-2xl font-semibold text-white">{{ formatPercent(props.memoryPercent) }}</p>
        <p class="mt-1 text-sm text-slate-400">{{ formatBytes(props.usedBytes) }} / {{ formatBytes(props.totalBytes) }}</p>
      </div>
    </div>
  </article>
</template>
```

```vue
<script setup lang="ts">
defineProps<{ title: string; body: string; actionLabel?: string }>()
const emit = defineEmits<{ action: [] }>()
</script>

<template>
  <div class="rounded-2xl border border-dashed border-slate-700 bg-slate-900/50 p-6 text-left">
    <h3 class="text-lg font-medium text-white">{{ title }}</h3>
    <p class="mt-2 text-sm leading-7 text-slate-400">{{ body }}</p>
    <button
      v-if="actionLabel"
      class="mt-4 rounded-xl bg-cyan-400 px-4 py-2 text-sm font-medium text-slate-950"
      @click="emit('action')"
    >
      {{ actionLabel }}
    </button>
  </div>
</template>
```

`AppShell.vue`, `AppHeader.vue`, `DesktopSidebar.vue`, `MobileTabbar.vue`, `ToastHost.vue`, and `LoadingBlock.vue` should compose these primitives and reserve slots for `OverviewPanel`, `ContextPanel`, `AccountsPanel`, `InstancesPanel`, and `LogsPanel`.

- [ ] **Step 2: Build the shared shell**

Run: `npm run build`
Expected: PASS with the new component tree compiling

- [ ] **Step 3: Commit the shell**

```bash
git add frontend/src/components/layout frontend/src/components/shared
git commit -m "feat: add responsive frontend shell"
```

### Task 6: Implement dashboard, context, and logs presentation

**Files:**
- Create: `frontend/src/components/panels/OverviewPanel.vue`
- Create: `frontend/src/components/panels/ContextPanel.vue`
- Create: `frontend/src/components/panels/LogsPanel.vue`
- Create: `frontend/src/components/logs/LogList.vue`
- Modify: `frontend/src/components/layout/AppShell.vue`

- [ ] **Step 1: Implement the overview and log list panels**

Use the store-derived state and keep global overview separate from selected context.

```vue
<script setup lang="ts">
import { storeToRefs } from 'pinia'
import ResourceCard from '@/components/shared/ResourceCard.vue'
import StatCard from '@/components/shared/StatCard.vue'
import { useAccountsStore } from '@/stores/accounts'
import { useInstancesStore } from '@/stores/instances'
import { useLogsStore } from '@/stores/logs'
import { useStatusStore } from '@/stores/status'

const statusStore = useStatusStore()
const accountsStore = useAccountsStore()
const instancesStore = useInstancesStore()
const logsStore = useLogsStore()

const { cpuPercent, memorySummary, clusterStatus } = storeToRefs(statusStore)
const { enabledCount, loggedInCount } = storeToRefs(accountsStore)
const { totalCount, runningCount, nonRunningCount } = storeToRefs(instancesStore)
const { recentSummary } = storeToRefs(logsStore)
</script>

<template>
  <section class="space-y-6">
    <div class="grid gap-4 xl:grid-cols-3">
      <StatCard title="账号总数" :value="String(accountsStore.items.length)" hint="管理中的全部账号" />
      <StatCard title="已启用账号" :value="String(enabledCount)" hint="enabled = true" />
      <StatCard title="已登录账号" :value="String(loggedInCount)" hint="auth_status = logged_in" />
      <StatCard title="实例总数" :value="String(totalCount)" hint="全部实例" />
      <StatCard title="运行中实例" :value="String(runningCount)" hint="status = running" />
      <StatCard title="非运行实例" :value="String(nonRunningCount)" hint="所有非 running 状态" />
    </div>
    <div class="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <ResourceCard
        :cpu-percent="cpuPercent"
        :used-bytes="memorySummary?.used_bytes ?? 0"
        :total-bytes="memorySummary?.total_bytes ?? 0"
        :memory-percent="memorySummary?.used_percent ?? 0"
      />
      <article class="rounded-2xl border border-slate-800 bg-slate-900/80 p-5">
        <p class="text-xs uppercase tracking-[0.25em] text-slate-400">集群状态</p>
        <p class="mt-3 text-3xl font-semibold text-white">{{ clusterStatus }}</p>
      </article>
    </div>
    <LogList :items="recentSummary" title="最近操作" />
  </section>
</template>
```

- [ ] **Step 2: Wire the shell so overview only shows when no selection exists**

Update `AppShell.vue` to render `OverviewPanel` when both `selectedAccountId` and `selectedInstanceId` are empty, and otherwise render `ContextPanel` in the center column.

- [ ] **Step 3: Build and verify the dashboard compiles**

Run: `npm run build`
Expected: PASS

- [ ] **Step 4: Commit the overview layer**

```bash
git add frontend/src/components/panels frontend/src/components/logs frontend/src/components/layout/AppShell.vue
git commit -m "feat: add overview and context panels"
```

### Task 7: Implement account management and Microsoft device-code login

**Files:**
- Create: `frontend/src/components/accounts/AccountList.vue`
- Create: `frontend/src/components/accounts/AccountCard.vue`
- Create: `frontend/src/components/accounts/AccountCreateDialog.vue`
- Create: `frontend/src/components/accounts/MicrosoftLoginDialog.vue`
- Create: `frontend/src/components/panels/AccountsPanel.vue`
- Create: `frontend/src/__tests__/microsoft-login-dialog.test.ts`

- [ ] **Step 1: Write a failing test for single active device-code sessions**

```ts
import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { useAccountsStore } from '@/stores/accounts'

describe('microsoft login session', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('replaces the previous active session when a new login starts', () => {
    const store = useAccountsStore()
    store.loginSession = {
      accountId: 'acc-a',
      userCode: 'AAAA-BBBB',
      verificationUri: 'https://example.com',
      verificationUriComplete: 'https://example.com/full',
      expiresAt: 1,
      pollIntervalSeconds: 5,
      status: 'pending',
      message: 'waiting',
    }

    store.clearLoginSession()

    expect(store.loginSession).toBeNull()
  })
})
```

- [ ] **Step 2: Run the login dialog test**

Run: `npm run test -- src/__tests__/microsoft-login-dialog.test.ts`
Expected: FAIL until the dialog exists

- [ ] **Step 3: Implement the account panel and dialogs**

Use props-down/events-up components. `MicrosoftLoginDialog.vue` must show one active code only, poll using the store session, and close by clearing the session.

- [ ] **Step 4: Run the relevant tests**

Run: `npm run test -- src/__tests__/accounts-store.test.ts src/__tests__/microsoft-login-dialog.test.ts`
Expected: PASS

- [ ] **Step 5: Commit account workflows**

```bash
git add frontend/src/components/accounts frontend/src/components/panels/AccountsPanel.vue frontend/src/__tests__/microsoft-login-dialog.test.ts
git commit -m "feat: add account management workflows"
```

### Task 8: Implement instance management with required account binding

**Files:**
- Create: `frontend/src/components/instances/InstanceList.vue`
- Create: `frontend/src/components/instances/InstanceCard.vue`
- Create: `frontend/src/components/instances/InstanceCreateDialog.vue`
- Create: `frontend/src/components/panels/InstancesPanel.vue`
- Create: `frontend/src/__tests__/instance-create-dialog.test.ts`

- [ ] **Step 1: Write a failing test for create-instance guard rails**

```ts
import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import InstanceCreateDialog from '@/components/instances/InstanceCreateDialog.vue'

describe('instance create dialog', () => {
  it('disables submit when the selected account is not eligible', async () => {
    const wrapper = mount(InstanceCreateDialog, {
      props: {
        open: true,
        accounts: [
          { id: 'acc-main', label: 'Main', enabled: false, auth_status: 'logged_in', has_token: true },
        ],
      },
    })

    await wrapper.get('select').setValue('acc-main')

    expect(wrapper.get('button[type="submit"]').attributes('disabled')).toBeDefined()
  })
})
```

- [ ] **Step 2: Run the create-instance dialog test to verify failure**

Run: `npm run test -- src/__tests__/instance-create-dialog.test.ts`
Expected: FAIL with missing component module

- [ ] **Step 3: Implement the instance panel and dialog**

`InstanceCreateDialog.vue` must:

```vue
<script setup lang="ts">
import { computed, reactive } from 'vue'
import type { AccountApiItem } from '@/types/api'
import type { CreateInstanceInput } from '@/types/domain'

const props = defineProps<{
  open: boolean
  accounts: AccountApiItem[]
  presetAccountId?: string
}>()

const emit = defineEmits<{
  close: []
  submit: [payload: CreateInstanceInput]
}>()

const form = reactive<CreateInstanceInput>({
  id: '',
  account_id: props.presetAccountId ?? '',
  server_address: '',
  enabled: true,
  auto_start: false,
})

const selectedAccount = computed(() => props.accounts.find((item) => item.id === form.account_id))
const accountEligible = computed(() => {
  const account = selectedAccount.value
  return !!account && account.enabled === true && account.auth_status === 'logged_in' && account.has_token === true
})
const submitDisabled = computed(() => !form.id || !form.server_address || !form.account_id || !accountEligible.value || (!form.enabled && form.auto_start))

function onSubmit() {
  if (submitDisabled.value) return
  emit('submit', { ...form, auto_start: form.enabled ? form.auto_start : false })
}
</script>
```

- [ ] **Step 4: Run the instance dialog and store tests**

Run: `npm run test -- src/__tests__/instance-create-dialog.test.ts src/__tests__/instances-store.test.ts`
Expected: PASS

- [ ] **Step 5: Commit instance workflows**

```bash
git add frontend/src/components/instances frontend/src/components/panels/InstancesPanel.vue frontend/src/__tests__/instance-create-dialog.test.ts
git commit -m "feat: add instance management workflows"
```

### Task 9: Polish the app shell wiring and mobile behavior

**Files:**
- Modify: `frontend/src/components/layout/AppShell.vue`
- Modify: `frontend/src/components/layout/AppHeader.vue`
- Modify: `frontend/src/components/layout/DesktopSidebar.vue`
- Modify: `frontend/src/components/layout/MobileTabbar.vue`
- Modify: `frontend/src/components/panels/AccountsPanel.vue`
- Modify: `frontend/src/components/panels/InstancesPanel.vue`
- Modify: `frontend/src/components/panels/LogsPanel.vue`
- Modify: `frontend/src/style.css`

- [ ] **Step 1: Wire all panels into the shell**

Desktop behavior:

- left column: `AccountsPanel`
- center column: `OverviewPanel` or `ContextPanel`
- right column: `InstancesPanel`
- logs summary only in `overview`

Mobile behavior:

- bottom tabs switch `uiStore.mode`
- render one primary panel at a time
- dialogs fill most of the viewport width

- [ ] **Step 2: Verify the app works in dev and build mode**

Run: `npm run build`
Expected: PASS

- [ ] **Step 3: Commit the shell integration**

```bash
git add frontend/src/components frontend/src/style.css
git commit -m "feat: compose responsive gmcc console"
```

### Task 10: Verify packager compatibility and docs

**Files:**
- Modify: `tools/packager/main_test.go`
- Modify: `README.md`

- [ ] **Step 1: Add or update packager coverage for frontend output compatibility**

Extend the packager tests so the generated Vite-style output remains accepted. Add assertions for:

```go
func TestCollectFrontendAssetsIncludesViteStyleOutput(t *testing.T) {
	// ensure index.html and hashed assets are copied
}
```

- [ ] **Step 2: Update the README with frontend workflow commands**

Add a short section showing:

```bash
cd frontend
npm install
npm run dev
npm run build
cd ..
go run ./tools/packager
```

Also note that API requests use same-origin `/api` and that device-code login is single-session in the first frontend release.

- [ ] **Step 3: Run all required verification commands**

Run these commands in order:

```bash
cd frontend && npm run test
cd frontend && npm run build
go test ./tools/packager
go test ./...
go build -o gmcc.exe ./cmd/gmcc
```

Expected: all commands PASS

- [ ] **Step 4: Commit the frontend release integration**

```bash
git add README.md tools/packager/main_test.go frontend
git commit -m "feat: add embedded vue management console"
```

## Self-Review

- Spec coverage: this plan covers SPA scaffolding, typed API client, stores, refresh coordination, dashboard, account management, Microsoft device-code login, instance creation with account binding, logs, responsive shell, packager compatibility, docs, and verification.
- Placeholder scan: no `TODO`, `TBD`, or “implement later” markers remain; each task includes file paths, commands, and concrete code.
- Type consistency: `CreateInstanceInput`, `MicrosoftSession`, `OperationLogItem`, store action names, and panel names are consistent across tasks.
