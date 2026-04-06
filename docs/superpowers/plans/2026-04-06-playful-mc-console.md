# Playful MC Console Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a playful, production-usable Vue 3 frontend for `gmcc` that covers dashboard, account management, Microsoft device-code login, instance management, resource monitoring, operation logs, and embedded frontend delivery.

**Architecture:** Keep the frontend in `frontend/` as a Vite + Vue + TypeScript SPA using a small global UI store, page-level query/store modules, and reusable display components. The app uses same-origin `/api` requests, shares one global Microsoft login task across views, and ships through the existing `frontend/dist -> tools/packager -> internal/webui/dist` embedding flow.

**Tech Stack:** Vue 3, TypeScript, Vite, Pinia, plain CSS, Vitest, Vue Test Utils, Go, Gin

---

## File Map

- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tsconfig.app.json`
- Create: `frontend/tsconfig.node.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/style.css`
- Create: `frontend/src/env.d.ts`
- Create: `frontend/src/types/api.ts`
- Create: `frontend/src/types/view.ts`
- Create: `frontend/src/lib/http.ts`
- Create: `frontend/src/lib/format.ts`
- Create: `frontend/src/lib/mappers.ts`
- Create: `frontend/src/lib/theme.ts`
- Create: `frontend/src/api/client.ts`
- Create: `frontend/src/stores/ui.ts`
- Create: `frontend/src/stores/home.ts`
- Create: `frontend/src/stores/accounts.ts`
- Create: `frontend/src/stores/instances.ts`
- Create: `frontend/src/stores/resources.ts`
- Create: `frontend/src/stores/logs.ts`
- Create: `frontend/src/components/layout/AppShell.vue`
- Create: `frontend/src/components/layout/SidebarNav.vue`
- Create: `frontend/src/components/layout/MobileNav.vue`
- Create: `frontend/src/components/layout/TopStatusBar.vue`
- Create: `frontend/src/components/shared/BaseCard.vue`
- Create: `frontend/src/components/shared/StatusBadge.vue`
- Create: `frontend/src/components/shared/EmptyState.vue`
- Create: `frontend/src/components/shared/InlineError.vue`
- Create: `frontend/src/components/shared/ToastHost.vue`
- Create: `frontend/src/components/home/HomeView.vue`
- Create: `frontend/src/components/home/HeroSummaryCard.vue`
- Create: `frontend/src/components/home/AccountSummaryPanel.vue`
- Create: `frontend/src/components/home/InstanceSummaryPanel.vue`
- Create: `frontend/src/components/home/LogSummaryPanel.vue`
- Create: `frontend/src/components/accounts/AccountsView.vue`
- Create: `frontend/src/components/accounts/AccountCard.vue`
- Create: `frontend/src/components/accounts/CreateAccountDialog.vue`
- Create: `frontend/src/components/accounts/MicrosoftLoginPanel.vue`
- Create: `frontend/src/components/instances/InstancesView.vue`
- Create: `frontend/src/components/instances/InstanceCard.vue`
- Create: `frontend/src/components/instances/CreateInstanceDialog.vue`
- Create: `frontend/src/components/resources/ResourcesView.vue`
- Create: `frontend/src/components/resources/ResourceOverviewCard.vue`
- Create: `frontend/src/components/logs/LogsView.vue`
- Create: `frontend/src/components/logs/LogTimeline.vue`
- Create: `frontend/src/__tests__/mappers.test.ts`
- Create: `frontend/src/__tests__/home-store.test.ts`
- Create: `frontend/src/__tests__/accounts-store.test.ts`
- Create: `frontend/src/__tests__/instances-store.test.ts`
- Create: `frontend/src/__tests__/ui-store.test.ts`
- Create: `frontend/src/__tests__/create-instance-dialog.test.ts`
- Create: `frontend/src/__tests__/microsoft-login-panel.test.ts`
- Modify: `README.md`
- Modify: `tools/packager/main_test.go`

## Implementation Notes

- Use Vue 3 Composition API with `<script setup lang="ts">` in all SFCs.
- Keep page switching inside the app shell; first version does not need Vue Router.
- Use one global login task in `ui` store; both Home and Accounts views consume it.
- Use same-origin relative `/api` requests only.
- Prefer CSS variables and handcrafted styles over a component library.
- Treat first version delete flows as out of scope.
- Build output must work with the current embed workflow and root-path deployment.

### Task 1: Scaffold the frontend project

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tsconfig.app.json`
- Create: `frontend/tsconfig.node.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/src/style.css`
- Create: `frontend/src/env.d.ts`

- [ ] **Step 1: Create the frontend toolchain files**

```json
{
  "name": "gmcc-playful-console",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "test": "vitest run"
  },
  "dependencies": {
    "pinia": "^3.0.3",
    "vue": "^3.5.13"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.2.3",
    "@vue/test-utils": "^2.4.6",
    "jsdom": "^26.1.0",
    "typescript": "^5.8.3",
    "vite": "^6.2.4",
    "vitest": "^3.1.1",
    "vue-tsc": "^2.2.8"
  }
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

- [ ] **Step 2: Add the minimal app bootstrap**

```ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import './style.css'

createApp(App).use(createPinia()).mount('#app')
```

```vue
<template>
  <div id="app-root">gmcc playful console bootstrap</div>
</template>
```

- [ ] **Step 3: Install dependencies**

Run: `npm install`
Expected: install completes and generates `package-lock.json`

- [ ] **Step 4: Verify build scaffolding**

Run: `npm run build`
Expected: Vite build succeeds and writes `frontend/dist`

- [ ] **Step 5: Commit**

```bash
git add frontend/package.json frontend/package-lock.json frontend/tsconfig.json frontend/tsconfig.app.json frontend/tsconfig.node.json frontend/vite.config.ts frontend/index.html frontend/src/main.ts frontend/src/App.vue frontend/src/style.css frontend/src/env.d.ts
git commit -m "feat: scaffold frontend console app"
```

### Task 2: Define API types, HTTP client, and view-model mappers

**Files:**
- Create: `frontend/src/types/api.ts`
- Create: `frontend/src/types/view.ts`
- Create: `frontend/src/lib/http.ts`
- Create: `frontend/src/lib/format.ts`
- Create: `frontend/src/lib/mappers.ts`
- Create: `frontend/src/api/client.ts`
- Create: `frontend/src/__tests__/mappers.test.ts`

- [ ] **Step 1: Write the failing mapper tests**

```ts
import { describe, expect, it } from 'vitest'
import { mapAccountStatus, mapInstanceStatus, mapComfortLevel } from '@/lib/mappers'

describe('mappers', () => {
  it('maps logged in enabled account to ready state', () => {
    expect(mapAccountStatus({ enabled: true, auth_status: 'logged_in', has_token: true })).toBe('ready')
  })

  it('maps unknown instance status to cautious state', () => {
    expect(mapInstanceStatus('mystery')).toBe('unknown')
  })

  it('maps healthy resource snapshot to comfort label', () => {
    expect(mapComfortLevel({ clusterStatus: 'running', cpuPercent: 25, memoryPercent: 32 })).toBe('comfort')
  })
})
```

- [ ] **Step 2: Run tests to verify failure**

Run: `npm run test -- mappers.test.ts`
Expected: FAIL because mapper functions do not exist yet

- [ ] **Step 3: Implement the API contracts and mapper helpers**

```ts
export type ApiAccount = {
  id: string
  label?: string
  note?: string
  enabled: boolean
  auth_status?: string
  has_token?: boolean
}

export function mapAccountStatus(account: Pick<ApiAccount, 'enabled' | 'auth_status' | 'has_token'>) {
  if (!account.enabled) return 'disabled'
  if (account.auth_status === 'logged_in' && account.has_token) return 'ready'
  if (account.has_token === false) return 'pending'
  return 'unknown'
}
```

```ts
export function mapInstanceStatus(status?: string) {
  if (status === 'running') return 'running'
  if (status === 'stopped') return 'stopped'
  return 'unknown'
}

export function mapComfortLevel(input: { clusterStatus?: string; cpuPercent?: number; memoryPercent?: number }) {
  const peak = Math.max(input.cpuPercent ?? 0, input.memoryPercent ?? 0)
  if (input.clusterStatus === 'stopped') return 'quiet'
  if (peak >= 85) return 'tense'
  if (peak >= 60) return 'busy'
  return 'comfort'
}
```

- [ ] **Step 4: Verify tests pass**

Run: `npm run test -- mappers.test.ts`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add frontend/src/types/api.ts frontend/src/types/view.ts frontend/src/lib/http.ts frontend/src/lib/format.ts frontend/src/lib/mappers.ts frontend/src/api/client.ts frontend/src/__tests__/mappers.test.ts
git commit -m "feat: add frontend api contracts and mappers"
```

### Task 3: Build the global UI store and app shell

**Files:**
- Create: `frontend/src/stores/ui.ts`
- Create: `frontend/src/lib/theme.ts`
- Create: `frontend/src/components/layout/AppShell.vue`
- Create: `frontend/src/components/layout/SidebarNav.vue`
- Create: `frontend/src/components/layout/MobileNav.vue`
- Create: `frontend/src/components/layout/TopStatusBar.vue`
- Create: `frontend/src/components/shared/BaseCard.vue`
- Create: `frontend/src/components/shared/StatusBadge.vue`
- Create: `frontend/src/components/shared/ToastHost.vue`
- Modify: `frontend/src/App.vue`
- Modify: `frontend/src/style.css`
- Create: `frontend/src/__tests__/ui-store.test.ts`

- [ ] **Step 1: Write the failing UI store test**

```ts
import { setActivePinia, createPinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import { useUiStore } from '@/stores/ui'

describe('ui store', () => {
  it('shares one login task across views', () => {
    setActivePinia(createPinia())
    const ui = useUiStore()
    ui.openLoginPanel('acc-main')
    expect(ui.loginTask.accountId).toBe('acc-main')
    expect(ui.loginPanelOpen).toBe(true)
  })
})
```

- [ ] **Step 2: Run the test to verify failure**

Run: `npm run test -- ui-store.test.ts`
Expected: FAIL because the store does not exist yet

- [ ] **Step 3: Implement the shell and shared UI state**

```ts
export const useUiStore = defineStore('ui', {
  state: () => ({
    activeView: 'home' as 'home' | 'accounts' | 'instances' | 'resources' | 'logs',
    loginPanelOpen: false,
    loginTask: {
      accountId: '',
      taskStatus: 'idle' as 'idle' | 'initializing' | 'polling' | 'succeeded' | 'failed' | 'expired' | 'replaced',
      userCode: '',
      verificationUriComplete: '',
      intervalSeconds: 5,
      errorMessage: '',
    },
    toasts: [] as Array<{ id: string; tone: 'success' | 'error'; message: string }>,
  }),
  actions: {
    setActiveView(view) { this.activeView = view },
    openLoginPanel(accountId: string) { this.loginPanelOpen = true; this.loginTask.accountId = accountId },
    closeLoginPanel() { this.loginPanelOpen = false },
  },
})
```

- [ ] **Step 4: Verify tests and build**

Run: `npm run test -- ui-store.test.ts && npm run build`
Expected: PASS, then Vite build succeeds

- [ ] **Step 5: Commit**

```bash
git add frontend/src/stores/ui.ts frontend/src/lib/theme.ts frontend/src/components/layout/AppShell.vue frontend/src/components/layout/SidebarNav.vue frontend/src/components/layout/MobileNav.vue frontend/src/components/layout/TopStatusBar.vue frontend/src/components/shared/BaseCard.vue frontend/src/components/shared/StatusBadge.vue frontend/src/components/shared/ToastHost.vue frontend/src/App.vue frontend/src/style.css frontend/src/__tests__/ui-store.test.ts
git commit -m "feat: add playful app shell and global ui state"
```

### Task 4: Implement the home dashboard query store and summary panels

**Files:**
- Create: `frontend/src/stores/home.ts`
- Create: `frontend/src/components/home/HomeView.vue`
- Create: `frontend/src/components/home/HeroSummaryCard.vue`
- Create: `frontend/src/components/home/AccountSummaryPanel.vue`
- Create: `frontend/src/components/home/InstanceSummaryPanel.vue`
- Create: `frontend/src/components/home/LogSummaryPanel.vue`
- Create: `frontend/src/components/shared/InlineError.vue`
- Create: `frontend/src/components/shared/EmptyState.vue`
- Create: `frontend/src/__tests__/home-store.test.ts`
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Write the failing home store test**

```ts
import { describe, expect, it, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useHomeStore } from '@/stores/home'

describe('home store', () => {
  it('keeps partial success visible when one module fails', async () => {
    setActivePinia(createPinia())
    const store = useHomeStore()
    store.$patch({ statusState: 'success', resourcesState: 'error' })
    expect(store.statusState).toBe('success')
    expect(store.resourcesState).toBe('error')
  })
})
```

- [ ] **Step 2: Run the test to verify failure**

Run: `npm run test -- home-store.test.ts`
Expected: FAIL because the store does not exist yet

- [ ] **Step 3: Implement the dashboard store and panels**

```ts
export const useHomeStore = defineStore('home', {
  state: () => ({
    statusState: 'idle' as LoadState,
    resourcesState: 'idle' as LoadState,
    accountsState: 'idle' as LoadState,
    instancesState: 'idle' as LoadState,
    logsState: 'idle' as LoadState,
    hero: null as HomeHeroView | null,
  }),
  actions: {
    async loadHome() {
      await Promise.allSettled([
        this.loadStatus(),
        this.loadResources(),
        this.loadAccounts(),
        this.loadInstances(),
        this.loadLogs(),
      ])
    },
  },
})
```

- [ ] **Step 4: Verify tests and render build**

Run: `npm run test -- home-store.test.ts && npm run build`
Expected: PASS, then Vite build succeeds

- [ ] **Step 5: Commit**

```bash
git add frontend/src/stores/home.ts frontend/src/components/home/HomeView.vue frontend/src/components/home/HeroSummaryCard.vue frontend/src/components/home/AccountSummaryPanel.vue frontend/src/components/home/InstanceSummaryPanel.vue frontend/src/components/home/LogSummaryPanel.vue frontend/src/components/shared/InlineError.vue frontend/src/components/shared/EmptyState.vue frontend/src/__tests__/home-store.test.ts frontend/src/App.vue
git commit -m "feat: add playful home dashboard"
```

### Task 5: Implement accounts view, create-account flow, and shared Microsoft login panel

**Files:**
- Create: `frontend/src/stores/accounts.ts`
- Create: `frontend/src/components/accounts/AccountsView.vue`
- Create: `frontend/src/components/accounts/AccountCard.vue`
- Create: `frontend/src/components/accounts/CreateAccountDialog.vue`
- Create: `frontend/src/components/accounts/MicrosoftLoginPanel.vue`
- Create: `frontend/src/__tests__/accounts-store.test.ts`
- Create: `frontend/src/__tests__/microsoft-login-panel.test.ts`
- Modify: `frontend/src/stores/ui.ts`

- [ ] **Step 1: Write the failing account and login tests**

```ts
it('blocks empty account id before submit', async () => {
  const payload = { id: '   ', label: '', note: '' }
  expect(validateCreateAccount(payload)).toEqual({ valid: false, message: '账号 ID 不能为空' })
})

it('marks previous login task as replaced when a new login starts', async () => {
  const ui = useUiStore()
  ui.loginTask.taskStatus = 'polling'
  ui.openLoginPanel('acc-old')
  ui.replaceLoginTask('acc-new')
  expect(ui.loginTask.accountId).toBe('acc-new')
  expect(ui.loginTask.taskStatus).toBe('initializing')
})
```

- [ ] **Step 2: Run tests to verify failure**

Run: `npm run test -- accounts-store.test.ts microsoft-login-panel.test.ts`
Expected: FAIL because validation and login panel logic are not implemented yet

- [ ] **Step 3: Implement account creation and login polling flow**

```ts
function validateCreateAccount(payload: { id: string }) {
  if (!payload.id.trim()) {
    return { valid: false, message: '账号 ID 不能为空' }
  }
  return { valid: true, message: '' }
}
```

```ts
async function startMicrosoftLogin(accountId: string) {
  const ui = useUiStore()
  ui.loginTask = { ...ui.loginTask, accountId, taskStatus: 'initializing', errorMessage: '' }
  const session = await api.initMicrosoftLogin(accountId)
  ui.loginTask = {
    ...ui.loginTask,
    taskStatus: 'polling',
    userCode: session.user_code,
    verificationUriComplete: session.verification_uri_complete,
    intervalSeconds: session.interval,
  }
}
```

- [ ] **Step 4: Verify tests and build**

Run: `npm run test -- accounts-store.test.ts microsoft-login-panel.test.ts && npm run build`
Expected: PASS, then build succeeds

- [ ] **Step 5: Commit**

```bash
git add frontend/src/stores/accounts.ts frontend/src/components/accounts/AccountsView.vue frontend/src/components/accounts/AccountCard.vue frontend/src/components/accounts/CreateAccountDialog.vue frontend/src/components/accounts/MicrosoftLoginPanel.vue frontend/src/__tests__/accounts-store.test.ts frontend/src/__tests__/microsoft-login-panel.test.ts frontend/src/stores/ui.ts
git commit -m "feat: add account management and microsoft login panel"
```

### Task 6: Implement instances view, create-instance validation, and lifecycle actions

**Files:**
- Create: `frontend/src/stores/instances.ts`
- Create: `frontend/src/components/instances/InstancesView.vue`
- Create: `frontend/src/components/instances/InstanceCard.vue`
- Create: `frontend/src/components/instances/CreateInstanceDialog.vue`
- Create: `frontend/src/__tests__/instances-store.test.ts`
- Create: `frontend/src/__tests__/create-instance-dialog.test.ts`

- [ ] **Step 1: Write the failing instance tests**

```ts
it('rejects auto start when instance is disabled', () => {
  expect(validateCreateInstance({ id: 'bot-1', accountId: 'acc-1', serverAddress: 'mc.test', enabled: false, autoStart: true })).toEqual({
    valid: false,
    message: '禁用实例不能启用自动启动',
  })
})

it('uses cautious action set for unknown status', () => {
  expect(getInstanceActions('unknown')).toEqual(['restart'])
})
```

- [ ] **Step 2: Run tests to verify failure**

Run: `npm run test -- instances-store.test.ts create-instance-dialog.test.ts`
Expected: FAIL because validation and action mapping do not exist yet

- [ ] **Step 3: Implement instance flows**

```ts
function validateCreateInstance(payload: CreateInstanceForm) {
  if (!payload.id.trim()) return { valid: false, message: '实例 ID 不能为空' }
  if (!payload.accountId) return { valid: false, message: '请选择账号' }
  if (!payload.serverAddress.trim()) return { valid: false, message: '服务器地址不能为空' }
  if (!payload.enabled && payload.autoStart) return { valid: false, message: '禁用实例不能启用自动启动' }
  return { valid: true, message: '' }
}
```

```ts
function getInstanceActions(status: 'running' | 'stopped' | 'unknown') {
  if (status === 'running') return ['stop', 'restart']
  if (status === 'stopped') return ['start']
  return ['restart']
}
```

- [ ] **Step 4: Verify tests and build**

Run: `npm run test -- instances-store.test.ts create-instance-dialog.test.ts && npm run build`
Expected: PASS, then build succeeds

- [ ] **Step 5: Commit**

```bash
git add frontend/src/stores/instances.ts frontend/src/components/instances/InstancesView.vue frontend/src/components/instances/InstanceCard.vue frontend/src/components/instances/CreateInstanceDialog.vue frontend/src/__tests__/instances-store.test.ts frontend/src/__tests__/create-instance-dialog.test.ts
git commit -m "feat: add instance management workflows"
```

### Task 7: Implement resources view and logs view

**Files:**
- Create: `frontend/src/stores/resources.ts`
- Create: `frontend/src/stores/logs.ts`
- Create: `frontend/src/components/resources/ResourcesView.vue`
- Create: `frontend/src/components/resources/ResourceOverviewCard.vue`
- Create: `frontend/src/components/logs/LogsView.vue`
- Create: `frontend/src/components/logs/LogTimeline.vue`
- Modify: `frontend/src/components/home/LogSummaryPanel.vue`

- [ ] **Step 1: Add the view modules**

```ts
export const useResourcesStore = defineStore('resources', {
  state: () => ({ state: 'idle' as LoadState, snapshot: null as ResourceView | null }),
  actions: {
    async loadResources() { /* call /api/resources and map to ResourceView */ },
  },
})
```

```ts
export const useLogsStore = defineStore('logs', {
  state: () => ({ state: 'idle' as LoadState, items: [] as LogTimelineItem[] }),
  actions: {
    async loadLogs() { /* call /api/logs/operations and map timeline items */ },
  },
})
```

- [ ] **Step 2: Verify the views compile**

Run: `npm run build`
Expected: build succeeds with resources and logs views wired into the shell

- [ ] **Step 3: Add smoke assertions if needed**

```ts
expect(mapComfortLevel({ clusterStatus: 'partial', cpuPercent: 62, memoryPercent: 70 })).toBe('busy')
```

- [ ] **Step 4: Run the full frontend tests**

Run: `npm run test`
Expected: all Vitest tests pass

- [ ] **Step 5: Commit**

```bash
git add frontend/src/stores/resources.ts frontend/src/stores/logs.ts frontend/src/components/resources/ResourcesView.vue frontend/src/components/resources/ResourceOverviewCard.vue frontend/src/components/logs/LogsView.vue frontend/src/components/logs/LogTimeline.vue frontend/src/components/home/LogSummaryPanel.vue frontend/src/__tests__
git commit -m "feat: add resource and log views"
```

### Task 8: Validate embed flow, update docs, and finish QA

**Files:**
- Modify: `README.md`
- Modify: `tools/packager/main_test.go`

- [ ] **Step 1: Update README frontend instructions**

```md
### Frontend development

```bash
cd frontend
npm install
npm run dev
```

### Frontend verification

```bash
cd frontend
npm run test
npm run build

cd ..
go run ./tools/packager
```
```

- [ ] **Step 2: Extend packager coverage if needed**

```go
func TestPackagerCopiesViteStyleAssets(t *testing.T) {
    // verify index.html and assets/* are preserved for embed flow
}
```

- [ ] **Step 3: Run full verification**

Run: `npm run test && npm run build && go test ./tools/packager && go test ./... && go run ./tools/packager`
Expected: frontend tests pass, build succeeds, Go tests pass, packager produces embedded binary output

- [ ] **Step 4: Manually verify the integrated app**

Run: `build/gmcc.exe` (or the platform-equivalent binary)
Expected: `/api/status` still responds and the root page serves the frontend; refreshing a non-root frontend view does not white-screen

- [ ] **Step 5: Commit**

```bash
git add README.md tools/packager/main_test.go frontend/dist internal/webui/dist
git commit -m "docs: document playful frontend workflow"
```

## Self-Review

- Spec coverage: this plan covers shell, dashboard, accounts, Microsoft login, instances, resources, logs, responsive shell, and embed verification.
- Placeholder scan: all tasks name exact files, commands, and minimal code direction; implementation details still need real code during execution, but no task is left as TBD.
- Type consistency: the plan keeps one global login task, uses `running | stopped | unknown` for instance status, and keeps delete flows out of scope for first version.
