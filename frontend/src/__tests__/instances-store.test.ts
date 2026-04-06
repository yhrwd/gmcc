import { describe, expect, it } from 'vitest'
import { getInstanceActions } from '@/stores/instances'

describe('instances action mapping', () => {
  it('uses start action for pending status', () => {
    expect(getInstanceActions('pending')).toEqual(['start'])
  })

  it('uses start and restart for error status', () => {
    expect(getInstanceActions('error')).toEqual(['start', 'restart'])
  })

  it('uses cautious action set for unknown status', () => {
    expect(getInstanceActions('unknown')).toEqual(['restart'])
  })
})
