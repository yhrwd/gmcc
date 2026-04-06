<script setup lang="ts">
import StatusBadge from '@/components/shared/StatusBadge.vue'
import { getInstanceActions } from '@/stores/instances'
import type { ViewInstance } from '@/types/view'

const props = defineProps<{
  instance: ViewInstance
  busyAction?: 'start' | 'stop' | 'restart' | 'delete' | ''
}>()

const emit = defineEmits<{
  action: [payload: { id: string; action: 'start' | 'stop' | 'restart' }]
  delete: [id: string]
}>()

const actions = getInstanceActions(props.instance.statusTone)

function getActionLabel(action: string) {
  if (props.busyAction === action) {
    if (action === 'start') return '启动中...'
    if (action === 'stop') return '停止中...'
    return '重启中...'
  }
  if (action === 'start') return '启动'
  if (action === 'stop') return '停止'
  return '重启'
}

function isBusy(action: string) {
  return props.busyAction === action || props.busyAction === 'delete'
}
</script>

<template>
  <article class="instance-card" :data-tone="instance.statusTone">
    <div class="instance-card__header">
      <div>
        <h3 class="instance-card__title">{{ instance.id }}</h3>
        <p class="instance-card__address">{{ instance.serverAddress }}</p>
      </div>
      <StatusBadge :label="instance.statusLabel" :tone="instance.statusTone" />
    </div>

    <div class="instance-card__meta">
      <span>账号 {{ instance.accountId }}</span>
      <span>在线 {{ instance.onlineDurationLabel }}</span>
      <span v-if="instance.positionLabel">坐标 {{ instance.positionLabel }}</span>
    </div>

    <div class="instance-card__bars">
      <span v-if="instance.health !== null">生命 {{ instance.health }}</span>
      <span v-if="instance.food !== null">饱食 {{ instance.food }}</span>
    </div>

    <div class="instance-card__actions">
      <button
        v-for="action in actions"
        :key="action"
        type="button"
        class="instance-card__button"
        :disabled="isBusy(action)"
        @click="emit('action', { id: instance.id, action: action as 'start' | 'stop' | 'restart' })"
      >
        {{ getActionLabel(action) }}
      </button>
      <button type="button" class="instance-card__button instance-card__button--ghost" :disabled="Boolean(props.busyAction)" @click="emit('delete', instance.id)">
        {{ props.busyAction === 'delete' ? '删除中...' : '删除' }}
      </button>
    </div>
  </article>
</template>

<style scoped>
.instance-card {
  display: grid;
  gap: 0.8rem;
  padding: 1rem;
  border-radius: 1.4rem;
  background: rgba(255, 255, 255, 0.58);
  border: 1px solid rgba(109, 132, 151, 0.14);
}

.instance-card[data-tone='running'] {
  box-shadow: inset 0 0 0 1px rgba(100, 205, 171, 0.16), 0 16px 34px rgba(75, 177, 152, 0.12);
}

.instance-card__header {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.instance-card__title,
.instance-card__address {
  margin: 0;
}

.instance-card__address,
.instance-card__meta,
.instance-card__bars {
  color: var(--text-soft);
}

.instance-card__meta,
.instance-card__bars,
.instance-card__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
}

.instance-card__button {
  border: none;
  border-radius: 999px 999px 999px 0.95rem;
  padding: 0.65rem 0.95rem;
  background: linear-gradient(135deg, rgba(219, 241, 255, 0.96), rgba(210, 226, 255, 0.92));
  color: #35556f;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 10px 20px rgba(120, 167, 219, 0.16);
}

.instance-card__button:disabled {
  cursor: wait;
  opacity: 0.7;
}

.instance-card__button--ghost {
  background: rgba(236, 241, 248, 0.95);
  color: #576980;
  box-shadow: none;
}
</style>
