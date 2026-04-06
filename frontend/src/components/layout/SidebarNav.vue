<script setup lang="ts">
import type { ActiveView } from '@/types/view'

defineProps<{
  activeView: ActiveView
}>()

const emit = defineEmits<{
  select: [view: ActiveView]
}>()

const items: Array<{ key: ActiveView; label: string; icon: string }> = [
  { key: 'home', label: '总览', icon: '01' },
  { key: 'accounts', label: '账号', icon: '02' },
  { key: 'instances', label: '实例', icon: '03' },
  { key: 'logs', label: '日志', icon: '04' },
]
</script>

<template>
  <nav class="sidebar-nav">
    <div class="sidebar-nav__brand">
      <p class="sidebar-nav__eyebrow">gmcc</p>
      <h1 class="sidebar-nav__title">gmcc 控制台</h1>
      <p class="sidebar-nav__subtitle">统一查看账号、实例和运行状态。</p>
    </div>

    <button
      v-for="item in items"
      :key="item.key"
      type="button"
      class="sidebar-nav__item"
      :data-active="item.key === activeView"
      @click="emit('select', item.key)"
    >
      <span class="sidebar-nav__icon">{{ item.icon }}</span>
      <span>{{ item.label }}</span>
    </button>
  </nav>
</template>

<style scoped>
.sidebar-nav {
  display: grid;
  gap: 0.75rem;
}

.sidebar-nav__brand {
  padding: 0.8rem 0.2rem 1.2rem;
  margin-bottom: 0.35rem;
  border-bottom: 1px dashed rgba(122, 141, 156, 0.18);
}

.sidebar-nav__eyebrow {
  margin: 0 0 0.35rem;
  color: var(--text-soft);
  text-transform: uppercase;
  letter-spacing: 0.24em;
  font-size: 0.7rem;
}

.sidebar-nav__title {
  margin: 0;
  font-size: 1.55rem;
  color: var(--text-main);
  letter-spacing: 0.01em;
}

.sidebar-nav__subtitle {
  margin: 0.45rem 0 0;
  color: var(--text-soft);
  line-height: 1.55;
  font-size: 0.92rem;
}

.sidebar-nav__item {
  display: flex;
  align-items: center;
  gap: 0.8rem;
  width: 100%;
  border: 1px solid rgba(94, 123, 140, 0.14);
  border-radius: 1rem;
  padding: 0.78rem 0.82rem;
  background: rgba(255, 255, 255, 0.52);
  color: var(--text-main);
  cursor: pointer;
  text-align: left;
}

.sidebar-nav__item[data-active='true'] {
  background: linear-gradient(135deg, rgba(238, 247, 242, 0.95), rgba(244, 236, 229, 0.92));
  transform: translateX(4px);
  border-color: rgba(242, 142, 22, 0.18);
  box-shadow: 0 10px 24px rgba(190, 140, 72, 0.1);
}

.sidebar-nav__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.85rem;
  height: 1.85rem;
  border-radius: 0.8rem;
  background: rgba(255, 255, 255, 0.84);
  font-size: 0.68rem;
  font-weight: 700;
}
</style>
