import { describe, expect, it } from 'vitest'
import { validateCreateInstance } from '@/stores/instances'

describe('instance validation', () => {
  it('rejects auto start when instance is disabled', () => {
    expect(validateCreateInstance({ id: 'bot-1', accountId: 'acc-1', serverAddress: 'mc.test', enabled: false, autoStart: true })).toEqual({
      valid: false,
      message: '禁用实例不能启用自动启动',
    })
  })
})
