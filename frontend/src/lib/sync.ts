export type SyncGroup = 'overview' | 'accounts' | 'instances' | 'logs'

type SyncReason = string
type Runner = () => Promise<void>
type TimerHandle = ReturnType<typeof setTimeout>

type GroupInterval = {
    foreground: number
    background: number
}

type CoordinatorDeps = {
    now: () => number
    isDocumentHidden: () => boolean
    setTimer: (fn: () => void, ms: number) => TimerHandle
    clearTimer: (id: TimerHandle) => void
}

type GroupState = {
    running: boolean
    pending: boolean
    failureCount: number
    timerId?: TimerHandle
    lastStartedAt: number
    lastFinishedAt: number
    lastReason: string
    lastError: string
    pendingReason: string
}

export type SyncDebugState = {
    started: boolean
    hidden: boolean
    activeGroups: SyncGroup[]
    groups: Record<SyncGroup, {
        running: boolean
        pending: boolean
        failureCount: number
        lastStartedAt: number
        lastFinishedAt: number
        lastReason: string
        lastError: string
        nextDelayMs: number | null
    }>
}

const GROUPS: SyncGroup[] = ['overview', 'accounts', 'instances', 'logs']

const GROUP_INTERVALS: Record<SyncGroup, GroupInterval> = {
    overview: { foreground: 8000, background: 20000 },
    accounts: { foreground: 8000, background: 20000 },
    instances: { foreground: 8000, background: 20000 },
    logs: { foreground: 20000, background: 60000 },
}

function createInitialState(): GroupState {
    return {
        running: false,
        pending: false,
        failureCount: 0,
        timerId: undefined,
        lastStartedAt: 0,
        lastFinishedAt: 0,
        lastReason: '',
        lastError: '',
        pendingReason: '',
    }
}

function getNextDelay(group: SyncGroup, hidden: boolean, failureCount: number): number {
    const base = hidden ? GROUP_INTERVALS[group].background : GROUP_INTERVALS[group].foreground
    if (failureCount >= 3) {
        return Math.max(base, 30000)
    }
    if (failureCount >= 2) {
        return Math.max(base, 15000)
    }
    return base
}

export class ConfigError extends Error {
    constructor(message: string) {
        super(message)
        this.name = 'ConfigError'
    }
}

class SyncCoordinator {
    private readonly deps: CoordinatorDeps
    private readonly runners: Partial<Record<SyncGroup, Runner>>
    private readonly states: Record<SyncGroup, GroupState>
    private readonly activeGroups: Set<SyncGroup>
    private started: boolean
    private wasHidden: boolean
    private visibilityHandler?: () => void

    constructor(deps: CoordinatorDeps) {
        this.deps = deps
        this.runners = {}
        this.states = {
            overview: createInitialState(),
            accounts: createInitialState(),
            instances: createInitialState(),
            logs: createInitialState(),
        }
        this.activeGroups = new Set<SyncGroup>()
        this.started = false
        this.wasHidden = false
    }

    register(group: SyncGroup, runner: Runner): void {
        this.runners[group] = runner
    }

    start(groups: SyncGroup[]): void {
        if (this.started) {
            return
        }

        this.assertRegistered(groups)
        this.activeGroups.clear()
        for (const group of groups) {
            this.activeGroups.add(group)
        }

        this.started = true
        this.wasHidden = this.deps.isDocumentHidden()
        if (typeof document !== 'undefined') {
            this.visibilityHandler = () => {
                const hidden = this.deps.isDocumentHidden()
                this.notifyVisibilityChange(hidden, this.wasHidden)
                this.wasHidden = hidden
            }
            document.addEventListener('visibilitychange', this.visibilityHandler)
        }

        for (const group of this.activeGroups) {
            this.requestNow([group], 'startup')
        }
    }

    stop(): void {
        if (!this.started) {
            return
        }

        this.started = false
        if (this.visibilityHandler && typeof document !== 'undefined') {
            document.removeEventListener('visibilitychange', this.visibilityHandler)
            this.visibilityHandler = undefined
        }
        for (const group of GROUPS) {
            const state = this.states[group]
            state.pending = false
            this.clearGroupTimer(group)
        }
    }

    request(groups: SyncGroup[], _reason: SyncReason): void {
        const reason = _reason || 'requested'
        this.assertRegistered(groups)
        if (!this.started) {
            return
        }

        for (const group of groups) {
            if (!this.activeGroups.has(group)) {
                continue
            }
            const state = this.states[group]
            if (state.running) {
                state.pending = true
                state.pendingReason = reason
                continue
            }
            state.pendingReason = reason
            this.scheduleGroup(group)
        }
    }

    requestNow(groups: SyncGroup[], _reason: SyncReason): void {
        const reason = _reason || 'requested-now'
        this.assertRegistered(groups)
        if (!this.started) {
            return
        }

        for (const group of groups) {
            if (!this.activeGroups.has(group)) {
                continue
            }
            const state = this.states[group]
            if (state.running) {
                state.pending = true
                state.pendingReason = reason
                continue
            }
            state.pendingReason = reason
            void this.runGroup(group)
        }
    }

    notifyVisibilityChange(isHidden: boolean, wasHidden: boolean): void {
        if (!this.started || !wasHidden || isHidden) {
            return
        }

        const visibleRefreshGroups: SyncGroup[] = ['overview', 'accounts', 'instances']
        this.requestNow(visibleRefreshGroups.filter((group) => this.activeGroups.has(group)), 'visibility-visible')
    }

    resetForTest(): void {
        this.stop()
        this.activeGroups.clear()
        this.wasHidden = false
        for (const group of GROUPS) {
            this.clearGroupTimer(group)
            this.states[group] = createInitialState()
            delete this.runners[group]
        }
    }

    getDebugSnapshot(): SyncDebugState {
        return {
            started: this.started,
            hidden: this.deps.isDocumentHidden(),
            activeGroups: [...this.activeGroups],
            groups: {
                overview: this.describeGroup('overview'),
                accounts: this.describeGroup('accounts'),
                instances: this.describeGroup('instances'),
                logs: this.describeGroup('logs'),
            },
        }
    }

    private assertRegistered(groups: SyncGroup[]): void {
        for (const group of groups) {
            if (!this.runners[group]) {
                throw new ConfigError(`sync coordinator misconfigured: missing runner for ${group}`)
            }
        }
    }

    private scheduleGroup(group: SyncGroup): void {
        if (!this.started || !this.activeGroups.has(group)) {
            return
        }

        const state = this.states[group]
        if (state.running) {
            state.pending = true
            return
        }

        this.clearGroupTimer(group)
        const delay = getNextDelay(group, this.deps.isDocumentHidden(), state.failureCount)
        state.timerId = this.deps.setTimer(() => {
            void this.runGroup(group)
        }, delay)
    }

    private clearGroupTimer(group: SyncGroup): void {
        const state = this.states[group]
        if (state.timerId === undefined) {
            return
        }
        this.deps.clearTimer(state.timerId)
        state.timerId = undefined
    }

    private async runGroup(group: SyncGroup): Promise<void> {
        const state = this.states[group]
        if (!this.started || !this.activeGroups.has(group) || state.running) {
            return
        }

        const runner = this.runners[group]
        if (!runner) {
            throw new ConfigError(`sync coordinator misconfigured: missing runner for ${group}`)
        }

        state.running = true
        state.pending = false
        state.lastStartedAt = this.deps.now()
        state.lastReason = state.pendingReason || 'scheduled'
        state.pendingReason = ''
        state.lastError = ''
        this.clearGroupTimer(group)

        try {
            await runner()
            state.failureCount = 0
            state.lastFinishedAt = this.deps.now()
        } catch (error) {
            state.failureCount += 1
            state.lastFinishedAt = this.deps.now()
            state.lastError = error instanceof Error ? error.message : 'unknown sync error'
        } finally {
            state.running = false
            if (!this.started || !this.activeGroups.has(group)) {
                return
            }
            if (state.pending) {
                queueMicrotask(() => {
                    void this.runGroup(group)
                })
                return
            }
            this.scheduleGroup(group)
        }
    }

    private describeGroup(group: SyncGroup) {
        const state = this.states[group]
        return {
            running: state.running,
            pending: state.pending,
            failureCount: state.failureCount,
            lastStartedAt: state.lastStartedAt,
            lastFinishedAt: state.lastFinishedAt,
            lastReason: state.lastReason,
            lastError: state.lastError,
            nextDelayMs: this.started && this.activeGroups.has(group) && !state.running ? getNextDelay(group, this.deps.isDocumentHidden(), state.failureCount) : null,
        }
    }
}

export function createSyncCoordinator(deps: Partial<CoordinatorDeps> = {}) {
    return new SyncCoordinator({
        now: () => Date.now(),
        isDocumentHidden: () => typeof document !== 'undefined' && document.hidden,
        setTimer: (fn, ms) => globalThis.setTimeout(fn, ms),
        clearTimer: (id) => globalThis.clearTimeout(id),
        ...deps,
    })
}

export const syncCoordinator = createSyncCoordinator()
