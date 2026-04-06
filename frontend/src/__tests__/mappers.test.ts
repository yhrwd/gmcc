import { describe, expect, it } from 'vitest'
import { buildHeroView, mapAccountStatus, mapComfortLevel, mapInstanceStatus } from '@/lib/mappers'

describe('mappers', () => {
  it('maps logged in enabled account to ready state', () => {
    expect(mapAccountStatus({ enabled: true, auth_status: 'logged_in', has_token: true })).toBe('ready')
  })

  it('maps not logged in account to pending state', () => {
    expect(mapAccountStatus({ enabled: true, auth_status: 'not_logged_in', has_token: false })).toBe('pending')
  })

  it('maps invalid account to invalid state', () => {
    expect(mapAccountStatus({ enabled: true, auth_status: 'auth_invalid', has_token: true })).toBe('invalid')
  })

  it('maps unknown instance status to unknown', () => {
    expect(mapInstanceStatus('mystery')).toBe('unknown')
  })

  it('maps pending instance status correctly', () => {
    expect(mapInstanceStatus('pending')).toBe('pending')
  })

  it('maps healthy resource snapshot to comfort label', () => {
    expect(mapComfortLevel({ clusterStatus: 'running', cpuPercent: 25, memoryPercent: 32 })).toBe('comfort')
  })

  it('builds a readable hero summary', () => {
    const hero = buildHeroView({
      status: { clusterStatus: 'running', totalInstances: 4, runningInstances: 2, uptimeLabel: '2h' },
      resources: { cpuPercent: 20, memoryPercent: 40, usedMemoryLabel: '2 GB', totalMemoryLabel: '8 GB', collectedAtLabel: '04-06 12:00', comfortLevel: 'comfort' },
      accounts: [],
    })

    expect(hero.activeInstanceRatioLabel).toBe('2/4 运行中')
  })
})
