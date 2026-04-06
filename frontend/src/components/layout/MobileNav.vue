<script setup lang="ts">
import type { ActiveView } from '@/types/view'

defineProps<{
  activeView: ActiveView
}>()

const emit = defineEmits<{
  select: [view: ActiveView]
}>()

const items: Array<{ key: ActiveView; label: string }> = [
  { key: 'home', label: '总览' },
  { key: 'accounts', label: '账号' },
  { key: 'instances', label: '实例' },
  { key: 'logs', label: '日志' },
]
</script>

<template>
  <nav class="mobile-nav">
    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="mobile-nav__item"
      :data-active="item.key === activeView"
      @click="emit('select', item.key)"
    >
      {{ item.label }}
    </button>
  </nav>
</template>

<style scoped>
.mobile-nav {
  position: fixed;
  left: 0.75rem;
  right: 0.75rem;
  bottom: max(0.75rem, env(safe-area-inset-bottom));
  z-index: 30;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.45rem;
  padding: 0.55rem;
  border-radius: 1.2rem;
  background: linear-gradient(180deg, rgba(255, 252, 246, 0.94), rgba(245, 240, 233, 0.88));
  backdrop-filter: blur(18px);
  border: 1px solid rgba(120, 139, 160, 0.14);
  box-shadow: 0 14px 36px rgba(77, 90, 111, 0.14);
}

.mobile-nav__item {
  border: none;
  border-radius: 999px 999px 999px 0.9rem;
  padding: 0.68rem 0.2rem;
  background: transparent;
  color: var(--text-soft);
  font-weight: 700;
}

.mobile-nav__item[data-active='true'] {
  background: linear-gradient(135deg, rgba(238, 247, 242, 0.96), rgba(244, 236, 229, 0.94));
  color: var(--text-main);
  box-shadow: 0 8px 20px rgba(157, 136, 108, 0.14);
  border: 1px solid rgba(242, 142, 22, 0.14);
}
</style>
