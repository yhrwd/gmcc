import { describe, expect, it } from 'vitest'
import { validateCreateAccount } from '@/stores/accounts'
import { createPinia, setActivePinia } from 'pinia'
import { useAccountsStore } from '@/stores/accounts'

describe('accounts validation', () => {
  it('blocks empty account id before submit', () => {
    expect(validateCreateAccount({ id: '   ', label: '', note: '' })).toEqual({ valid: false, message: '账号 ID 不能为空' })
  })

  it('returns cached data without forcing reload', async () => {
    setActivePinia(createPinia())
      const store = useAccountsStore()
      store.state = 'success'
      store.items = [{ id: 'a', label: 'A', note: '', enabled: true, authStatus: 'logged_in', authStatusTone: 'ready', authStatusLabel: '在线冒险者', hasToken: true }]
      await store.loadAccounts()
      expect(store.items).toHaveLength(1)
  })
})
