import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import { useHomeStore } from '@/stores/home'

describe('home store', () => {
  it('keeps partial success visible when one module fails', () => {
    setActivePinia(createPinia())
    const store = useHomeStore()
    store.$patch({ statusState: 'success', resourcesState: 'error' })
    expect(store.statusState).toBe('success')
    expect(store.resourcesState).toBe('error')
  })
})
