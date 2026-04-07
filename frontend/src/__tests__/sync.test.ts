import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ConfigError, createSyncCoordinator, type SyncGroup } from '@/lib/sync'

function flushMicrotasks() {
    return Promise.resolve()
}

describe('sync coordinator', () => {
    beforeEach(() => {
        vi.useFakeTimers()
        vi.restoreAllMocks()
        Object.defineProperty(document, 'hidden', {
            value: false,
            configurable: true,
        })
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
        await flushMicrotasks()

        coordinator.requestNow(['accounts'], 'manual-refresh')
        coordinator.requestNow(['accounts'], 'manual-refresh-again')
        expect(runner).toHaveBeenCalledTimes(1)

        release()
        await flushMicrotasks()
        await flushMicrotasks()

        expect(runner).toHaveBeenCalledTimes(2)
    })

    it('throws ConfigError for unregistered groups', () => {
        const coordinator = createSyncCoordinator()

        expect(() => coordinator.requestNow(['overview'], 'missing-runner')).toThrow(ConfigError)
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
        await flushMicrotasks()
        calls.length = 0

        coordinator.notifyVisibilityChange(false, true)
        await flushMicrotasks()

        expect(calls.sort()).toEqual(['accounts', 'instances', 'overview'])
    })

    it('stop prevents after-stop execution and start/stop are idempotent', async () => {
        const runner = vi.fn(async () => undefined)
        const coordinator = createSyncCoordinator({
            now: () => 1,
            isDocumentHidden: () => false,
            setTimer: (fn, ms) => window.setTimeout(fn, ms),
            clearTimer: (id) => window.clearTimeout(id),
        })

        coordinator.register('overview', runner)
        coordinator.start(['overview'])
        coordinator.start(['overview'])
        await flushMicrotasks()

        expect(runner).toHaveBeenCalledTimes(1)

        coordinator.stop()
        coordinator.stop()
        coordinator.requestNow(['overview'], 'after-stop-immediate')
        coordinator.request(['overview'], 'after-stop-scheduled')
        await vi.advanceTimersByTimeAsync(30000)

        expect(runner).toHaveBeenCalledTimes(1)
    })

    it('exposes debug snapshot with latest reason and error state', async () => {
        let fail = true
        const coordinator = createSyncCoordinator({
            now: (() => {
                let tick = 10
                return () => tick++
            })(),
            isDocumentHidden: () => false,
            setTimer: (fn, ms) => window.setTimeout(fn, ms),
            clearTimer: (id) => window.clearTimeout(id),
        })

        coordinator.register('overview', async () => {
            if (fail) {
                fail = false
                throw new Error('temporary failure')
            }
        })

        coordinator.start(['overview'])
        await flushMicrotasks()

        const failedSnapshot = coordinator.getDebugSnapshot()
        expect(failedSnapshot.groups.overview.lastReason).toBe('startup')
        expect(failedSnapshot.groups.overview.lastError).toBe('temporary failure')
        expect(failedSnapshot.groups.overview.failureCount).toBe(1)

        coordinator.requestNow(['overview'], 'manual-retry')
        await flushMicrotasks()

        const successSnapshot = coordinator.getDebugSnapshot()
        expect(successSnapshot.groups.overview.lastReason).toBe('manual-retry')
        expect(successSnapshot.groups.overview.lastError).toBe('')
        expect(successSnapshot.groups.overview.failureCount).toBe(0)
    })

    it('starts with all four groups and stops without rescheduling pending work', async () => {
        const runner = vi.fn(async () => undefined)
        const coordinator = createSyncCoordinator({
            now: () => 1,
            isDocumentHidden: () => false,
            setTimer: (fn, ms) => window.setTimeout(fn, ms),
            clearTimer: (id) => window.clearTimeout(id),
        })

        coordinator.register('overview', runner)
        coordinator.register('accounts', runner)
        coordinator.register('instances', runner)
        coordinator.register('logs', runner)

        coordinator.start(['overview', 'accounts', 'instances', 'logs'])
        await flushMicrotasks()

        expect(runner).toHaveBeenCalledTimes(4)

        coordinator.stop()
        coordinator.requestNow(['overview'], 'after-stop')
        await flushMicrotasks()

        expect(runner).toHaveBeenCalledTimes(4)
    })
})
