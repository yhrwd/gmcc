# Frontend Sync Coordinator Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace fragile overlapping frontend polling with a coordinated sync layer that serializes same-group refreshes, triggers immediate refresh after successful account and instance actions, and stays stable on slow networks.

**Architecture:** Add a singleton sync coordinator in `frontend/src/lib/sync.ts` that owns runner registration, same-group serialization, pending reruns, timed scheduling, visibility-aware intervals, and retry backoff. `App.vue` becomes the lifecycle hook that registers store-backed runners and starts/stops the coordinator, while account and instance stores notify the coordinator instead of directly chaining ad-hoc refresh calls.

**Tech Stack:** Vue 3, Pinia, TypeScript, Vitest, Vite

---

## File Map

- Create: `frontend/src/lib/sync.ts`
  - Own the singleton coordinator, group types, registration, scheduling, and test reset support.
- Create: `frontend/src/__tests__/sync.test.ts`
  - Cover serialization, pending reruns, visibility refresh, stop behavior, and invalid call boundaries.
- Modify: `frontend/src/App.vue`
  - Register four runners from live Pinia stores and connect app mount/unmount to `syncCoordinator.start()` / `syncCoordinator.stop()`.
- Modify: `frontend/src/stores/home.ts`
  - Make `loadStatus`, `loadResources`, and `loadHome` reject on refresh failure while still preserving store error state.
- Modify: `frontend/src/stores/accounts.ts`
  - Preserve current optimistic UI and toasts, but replace direct refresh chaining with coordinator-triggered `accounts` / `overview` / `instances` refresh.
- Modify: `frontend/src/stores/instances.ts`
  - Preserve current optimistic UI and toasts, but replace direct refresh chaining with coordinator-triggered `instances` / `overview` refresh.
- Modify: `frontend/src/stores/logs.ts`
  - Make log refresh reject on failure so the coordinator can count retry backoff correctly.
- Modify: `frontend/src/__tests__/accounts-store.test.ts`
  - Verify successful account actions request the expected sync groups.
- Modify: `frontend/src/__tests__/instances-store.test.ts`
  - Verify successful instance actions request the expected sync groups.
- Modify: `frontend/src/__tests__/home-store.test.ts`
  - Verify `loadHome()` still preserves partial state locally while rejecting to the coordinator when either child request fails.

### Task 1: Build and test the sync coordinator core

**Files:**
- Create: `frontend/src/lib/sync.ts`
- Test: `frontend/src/__tests__/sync.test.ts`

- [ ] **Step 1: Write the failing coordinator tests**

```ts
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createSyncCoordinator, type SyncGroup } from '@/lib/sync'

describe('sync coordinator', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.restoreAllMocks()
    Object.defineProperty(document, 'hidden', { value: false, configurable: true })
  })

  it('reruns once after an in-flight request receives requestNow', async () => {
    let release!: () => void
    const runner = vi.fn(
      () =>
        new Promise<void>((resolve) => {
          release = resolve
        }),
    )
    const coordinator = createSyncCoordinator({
      now: () => 100,
      isDocumentHidden: () => false,
      setTimer: (fn, ms) => window.setTimeout(fn, ms),
      clearTimer: (id) => window.clearTimeout(id),
    })

    coordinator.register('accounts', runner)
    coordinator.start(['accounts'])
    await Promise.resolve()
    coordinator.requestNow(['accounts'], 'manual-refresh')
    expect(runner).toHaveBeenCalledTimes(1)

    release()
    await Promise.resolve()
    await Promise.resolve()

    expect(runner).toHaveBeenCalledTimes(2)
  })

  it('throws ConfigError for unregistered groups', () => {
    const coordinator = createSyncCoordinator()
    expect(() => coordinator.requestNow(['overview'], 'missing-runner')).toThrow('sync coordinator misconfigured')
  })

  it('refreshes overview/accounts/instances when page becomes visible again', async () => {
    const calls: SyncGroup[] = []
    const coordinator = createSyncCoordinator({
      isDocumentHidden: () => false,
      now: () => 1,
    })

    for (const group of ['overview', 'accounts', 'instances', 'logs'] as const) {
      coordinator.register(group, async () => {
        calls.push(group)
      })
    }

    coordinator.start(['overview', 'accounts', 'instances', 'logs'])
    calls.length = 0
    coordinator.notifyVisibilityChange(false, true)
    await Promise.resolve()

    expect(calls.sort()).toEqual(['accounts', 'instances', 'overview'])
  })
})
```

- [ ] **Step 2: Run the coordinator test file and confirm it fails**

Run: `npm test -- src/__tests__/sync.test.ts`
Expected: FAIL with missing `@/lib/sync` exports such as `createSyncCoordinator`

- [ ] **Step 3: Write the minimal coordinator implementation**

```ts
export type SyncGroup = 'overview' | 'accounts' | 'instances' | 'logs'
type SyncReason = string
type Runner = () => Promise<void>

export class ConfigError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'ConfigError'
  }
}

export function createSyncCoordinator(deps?: Partial<CoordinatorDeps>) {
  const coordinator = new SyncCoordinator({
    now: () => Date.now(),
    isDocumentHidden: () => typeof document !== 'undefined' && document.hidden,
    setTimer: (fn, ms) => window.setTimeout(fn, ms),
    clearTimer: (id) => window.clearTimeout(id),
    ...deps,
  })
  return coordinator
}

class SyncCoordinator {
  register(group: SyncGroup, runner: Runner) { /* keep map, block register after start */ }
  start(groups: SyncGroup[]) { /* validate groups, mark started, fire requestNow on each */ }
  stop() { /* stop timers, block reschedule */ }
  request(groups: SyncGroup[], reason: SyncReason) { /* natural scheduling only */ }
  requestNow(groups: SyncGroup[], reason: SyncReason) { /* immediate run or pending */ }
  notifyVisibilityChange(isHidden: boolean, wasHidden: boolean) { /* rerun main groups on visible */ }
  resetForTest() { /* clear timers, runtime state, registrations */ }
}

export const syncCoordinator = createSyncCoordinator()
```

- [ ] **Step 4: Fill in scheduling rules in `frontend/src/lib/sync.ts`**

```ts
const GROUP_INTERVALS: Record<SyncGroup, { foreground: number; background: number }> = {
  overview: { foreground: 8000, background: 20000 },
  accounts: { foreground: 8000, background: 20000 },
  instances: { foreground: 8000, background: 20000 },
  logs: { foreground: 20000, background: 60000 },
}

function nextDelay(group: SyncGroup, hidden: boolean, failureCount: number) {
  const base = hidden ? GROUP_INTERVALS[group].background : GROUP_INTERVALS[group].foreground
  if (failureCount >= 3) return Math.max(base, 30000)
  if (failureCount >= 2) return Math.max(base, 15000)
  return base
}

private async runGroup(group: SyncGroup) {
  const state = this.states[group]
  if (!this.started || state.running) return
  state.running = true
  state.pending = false
  state.lastStartedAt = this.deps.now()
  this.clearGroupTimer(group)

  try {
    await this.runners[group]!()
    state.failureCount = 0
  } catch (error) {
    state.failureCount += 1
  } finally {
    state.running = false
    if (!this.started) return
    if (state.pending) {
      queueMicrotask(() => void this.runGroup(group))
      return
    }
    this.scheduleGroup(group)
  }
}
```

- [ ] **Step 5: Run the coordinator tests and make them pass**

Run: `npm test -- src/__tests__/sync.test.ts`
Expected: PASS

- [ ] **Step 6: Commit the coordinator core**

```bash
git add frontend/src/lib/sync.ts frontend/src/__tests__/sync.test.ts
git commit -m "feat: add frontend sync coordinator"
```

### Task 2: Make store refresh methods reject correctly and document partial overview failure

**Files:**
- Modify: `frontend/src/stores/home.ts`
- Modify: `frontend/src/stores/logs.ts`
- Test: `frontend/src/__tests__/home-store.test.ts`

- [ ] **Step 1: Write failing tests for overview and logs failure behavior**

```ts
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { apiClient } from '@/api/client'
import { useHomeStore } from '@/stores/home'

describe('home store loadHome', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('rejects when either overview subrequest fails while preserving partial local state', async () => {
    vi.spyOn(apiClient, 'getStatus').mockResolvedValue({ cluster_status: 'running', total_instances: 1, running_instances: 1 } as never)
    vi.spyOn(apiClient, 'getResources').mockRejectedValue(new Error('resource down'))
    const store = useHomeStore()

    await expect(store.loadHome(true, true)).rejects.toThrow('resource down')
    expect(store.statusState).toBe('success')
    expect(store.resourcesState).toBe('error')
  })
})
```

- [ ] **Step 2: Run the home-store test and confirm it fails**

Run: `npm test -- src/__tests__/home-store.test.ts`
Expected: FAIL because `loadHome()` currently resolves even when one child request fails

- [ ] **Step 3: Update `frontend/src/stores/home.ts` so coordinator-visible failures reject**

```ts
async loadStatus(force = false, silent = false) {
  // keep current state handling
  try {
    const status = await apiClient.getStatus()
    this.statusSummary = mapStatusSummary(status)
    this.statusState = 'success'
  } catch (error) {
    this.statusState = 'error'
    this.statusError = error instanceof Error ? error.message : '状态获取失败'
    throw error instanceof Error ? error : new Error('状态获取失败')
  }
}

async loadHome(force = false, silent = false) {
  const [statusResult, resourcesResult] = await Promise.allSettled([
    this.loadStatus(force, silent),
    this.loadResources(force, silent),
  ])
  if (statusResult.status === 'rejected') throw statusResult.reason
  if (resourcesResult.status === 'rejected') throw resourcesResult.reason
}
```

- [ ] **Step 4: Update `frontend/src/stores/logs.ts` to reject on refresh failure**

```ts
async loadLogs(force = false, silent = false) {
  if (!force && this.state === 'success') return
  if (!silent || this.state === 'idle') this.state = 'loading'
  try {
    const result = await apiClient.getOperationLogs()
    this.items = (result.logs ?? []).map(mapLogItem)
    this.state = 'success'
  } catch (error) {
    this.state = 'error'
    this.errorMessage = error instanceof Error ? error.message : '日志读取失败'
    throw error instanceof Error ? error : new Error('日志读取失败')
  }
}
```

- [ ] **Step 5: Run the store failure tests and make them pass**

Run: `npm test -- src/__tests__/home-store.test.ts`
Expected: PASS

- [ ] **Step 6: Commit the store failure contract changes**

```bash
git add frontend/src/stores/home.ts frontend/src/stores/logs.ts frontend/src/__tests__/home-store.test.ts
git commit -m "fix: expose frontend refresh failures to coordinator"
```

### Task 3: Wire the app lifecycle into the coordinator

**Files:**
- Modify: `frontend/src/App.vue`
- Test: `frontend/src/__tests__/sync.test.ts`

- [ ] **Step 1: Extend the sync tests to cover app-facing lifecycle assumptions**

```ts
it('starts with all four groups and stops without rescheduling pending work', async () => {
  const runner = vi.fn(async () => {})
  const coordinator = createSyncCoordinator()

  coordinator.register('overview', runner)
  coordinator.register('accounts', runner)
  coordinator.register('instances', runner)
  coordinator.register('logs', runner)

  coordinator.start(['overview', 'accounts', 'instances', 'logs'])
  await Promise.resolve()
  expect(runner).toHaveBeenCalledTimes(4)

  coordinator.stop()
  coordinator.requestNow(['overview'], 'after-stop')
  await Promise.resolve()
  expect(runner).toHaveBeenCalledTimes(4)
})
```

- [ ] **Step 2: Run the sync tests and confirm any new lifecycle expectation fails before wiring `App.vue`**

Run: `npm test -- src/__tests__/sync.test.ts`
Expected: FAIL if `start()` / `stop()` semantics are incomplete

- [ ] **Step 3: Replace `App.vue` interval polling with coordinator registration**

```ts
onMounted(() => {
  syncCoordinator.register('overview', () => homeStore.loadHome(true, true))
  syncCoordinator.register('accounts', () => accountsStore.loadAccounts(true, true))
  syncCoordinator.register('instances', () => instancesStore.loadInstances(true, true))
  syncCoordinator.register('logs', () => logsStore.refreshLogs())
  syncCoordinator.start(['overview', 'accounts', 'instances', 'logs'])

  const handleVisibility = () => {
    syncCoordinator.notifyVisibilityChange(document.hidden, false)
  }
  document.addEventListener('visibilitychange', handleVisibility)
})

onBeforeUnmount(() => {
  syncCoordinator.stop()
})
```

- [ ] **Step 4: Keep visibility tracking inside the coordinator instead of `App.vue`**

```ts
start(groups: SyncGroup[]) {
  if (this.started) return
  this.started = true
  this.wasHidden = this.deps.isDocumentHidden()
  this.visibilityHandler = () => {
    const hidden = this.deps.isDocumentHidden()
    this.notifyVisibilityChange(hidden, this.wasHidden)
    this.wasHidden = hidden
  }
  document.addEventListener('visibilitychange', this.visibilityHandler)
  for (const group of groups) this.requestNow([group], 'initial-load')
}
```

- [ ] **Step 5: Run the sync tests again and confirm lifecycle behavior passes**

Run: `npm test -- src/__tests__/sync.test.ts`
Expected: PASS

- [ ] **Step 6: Commit the app lifecycle wiring**

```bash
git add frontend/src/App.vue frontend/src/lib/sync.ts frontend/src/__tests__/sync.test.ts
git commit -m "refactor: drive app refresh through sync coordinator"
```

### Task 4: Move account and instance success refreshes onto the coordinator

**Files:**
- Modify: `frontend/src/stores/accounts.ts`
- Modify: `frontend/src/stores/instances.ts`
- Modify: `frontend/src/__tests__/accounts-store.test.ts`
- Modify: `frontend/src/__tests__/instances-store.test.ts`

- [ ] **Step 1: Write failing tests for account action refresh groups**

```ts
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { apiClient } from '@/api/client'
import { syncCoordinator } from '@/lib/sync'
import { useAccountsStore } from '@/stores/accounts'

describe('accounts store sync integration', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('requests accounts refresh after create success', async () => {
    vi.spyOn(apiClient, 'createAccount').mockResolvedValue({ success: true } as never)
    const syncSpy = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => undefined)
    const store = useAccountsStore()

    await store.createAccount({ id: 'acc-main', label: 'Main', note: '' })

    expect(syncSpy).toHaveBeenCalledWith(['accounts'], 'account-created')
  })
})
```

- [ ] **Step 2: Write failing tests for instance action refresh groups**

```ts
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import { apiClient } from '@/api/client'
import { syncCoordinator } from '@/lib/sync'
import { useInstancesStore } from '@/stores/instances'

describe('instances store sync integration', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.restoreAllMocks()
  })

  it('requests instances and overview refresh after start success', async () => {
    vi.spyOn(apiClient, 'startInstance').mockResolvedValue({ success: true } as never)
    const syncSpy = vi.spyOn(syncCoordinator, 'requestNow').mockImplementation(() => undefined)
    const store = useInstancesStore()
    store.items = [{ id: 'bot-1', accountId: 'acc-1', serverAddress: 'mc.test', statusTone: 'pending', statusLabel: '待出勤', onlineDurationLabel: '0s', health: null, food: null, positionLabel: '' }]

    await store.runAction('bot-1', 'start')

    expect(syncSpy).toHaveBeenCalledWith(['instances', 'overview'], 'instance-started')
  })
})
```

- [ ] **Step 3: Run the account and instance store tests and confirm they fail**

Run: `npm test -- src/__tests__/accounts-store.test.ts src/__tests__/instances-store.test.ts`
Expected: FAIL because stores still call direct `loadXxx(true, true)` refreshes

- [ ] **Step 4: Replace direct refresh chaining in `frontend/src/stores/accounts.ts`**

```ts
import { syncCoordinator } from '@/lib/sync'

await apiClient.createAccount({ /* existing payload */ })
this.state = 'success'
syncCoordinator.requestNow(['accounts'], 'account-created')

// delete success
syncCoordinator.requestNow(['accounts', 'instances', 'overview'], 'account-deleted')

// login success
syncCoordinator.requestNow(['accounts', 'overview'], 'account-login-succeeded')
```

- [ ] **Step 5: Replace direct refresh chaining in `frontend/src/stores/instances.ts`**

```ts
import { syncCoordinator } from '@/lib/sync'

await apiClient.createInstance(payload)
this.state = 'success'
syncCoordinator.requestNow(['instances', 'overview'], 'instance-created')

// runAction success
syncCoordinator.requestNow(['instances', 'overview'], `instance-${action}ed`)

// delete success
syncCoordinator.requestNow(['instances', 'overview'], 'instance-deleted')
```

- [ ] **Step 6: Run the targeted store tests and make them pass**

Run: `npm test -- src/__tests__/accounts-store.test.ts src/__tests__/instances-store.test.ts`
Expected: PASS

- [ ] **Step 7: Commit the store/coordinator integration**

```bash
git add frontend/src/stores/accounts.ts frontend/src/stores/instances.ts frontend/src/__tests__/accounts-store.test.ts frontend/src/__tests__/instances-store.test.ts
git commit -m "fix: trigger coordinated refresh after frontend actions"
```

### Task 5: Run full frontend verification

**Files:**
- Modify if needed: `frontend/src/lib/sync.ts`
- Modify if needed: `frontend/src/stores/home.ts`
- Modify if needed: `frontend/src/stores/accounts.ts`
- Modify if needed: `frontend/src/stores/instances.ts`
- Modify if needed: `frontend/src/stores/logs.ts`

- [ ] **Step 1: Run the full frontend test suite**

Run: `npm test`
Expected: PASS

- [ ] **Step 2: Run the frontend production build**

Run: `npm run build`
Expected: PASS with generated Vite bundle output

- [ ] **Step 3: If tests or build fail, apply the smallest fix and rerun the exact failing command**

```ts
// Example minimal fix shape if a visibility mock leaks between tests
afterEach(() => {
  vi.useRealTimers()
  syncCoordinator.resetForTest()
})
```

- [ ] **Step 4: Commit the final verification fixes**

```bash
git add frontend/src/lib/sync.ts frontend/src/stores/home.ts frontend/src/stores/accounts.ts frontend/src/stores/instances.ts frontend/src/stores/logs.ts frontend/src/__tests__/sync.test.ts frontend/src/__tests__/home-store.test.ts frontend/src/__tests__/accounts-store.test.ts frontend/src/__tests__/instances-store.test.ts
git commit -m "test: cover coordinated frontend sync behavior"
```

## Self-Review

- Spec coverage: coordinator singleton, registration, same-group serialization, pending reruns, visibility refresh, fixed intervals/backoff, store-triggered refresh groups, and overview failure semantics are each mapped to a concrete task.
- Placeholder scan: no `TODO` / `TBD` / “similar to above” placeholders remain; each task includes explicit files, commands, and concrete code shapes.
- Type consistency: `SyncGroup`, `syncCoordinator`, `requestNow`, `register`, and `resetForTest` names are used consistently across tasks.
