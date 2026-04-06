<script setup lang="ts">
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import { computed } from 'vue'
import type { LoadState, ViewInstance } from '@/types/view'

const props = defineProps<{
  state: LoadState
  errorMessage: string
  instances: ViewInstance[]
}>()

const emit = defineEmits<{
  retry: []
}>()

const stats = computed(() => {
  const running = props.instances.filter((instance) => instance.statusTone === 'running').length
  const pending = props.instances.filter((instance) => instance.statusTone === 'pending' || instance.statusTone === 'starting' || instance.statusTone === 'reconnecting').length
  const stopped = props.instances.filter((instance) => instance.statusTone === 'stopped').length
  const error = props.instances.filter((instance) => instance.statusTone === 'error' || instance.statusTone === 'unknown').length
  const availability = props.instances.length > 0 ? `${Math.round((running / props.instances.length) * 100)}%` : '0%'

  return [
    { label: '总实例数', value: props.instances.length, tone: 'neutral' },
    { label: '运行中', value: running, tone: 'ready' },
    { label: '处理中', value: pending, tone: 'pending' },
    { label: '可用率', value: availability, tone: 'muted' },
  ]
})
</script>

<template>
  <BaseCard eyebrow="实例" title="实例状态" accent="sky">
    <InlineError v-if="state === 'error'" :message="errorMessage" button-label="重试" @retry="emit('retry')" />
    <EmptyState
      v-else-if="state === 'success' && !instances.length"
      title="暂无实例"
      description="请先完成账号登录，再创建实例。"
    />
    <div v-else class="instance-stats">
      <article v-for="item in stats" :key="item.label" class="instance-stats__item" :data-tone="item.tone">
        <span class="instance-stats__label">{{ item.label }}</span>
        <strong class="instance-stats__value">{{ item.value }}</strong>
      </article>
    </div>
  </BaseCard>
</template>

<style scoped>
.instance-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.85rem;
}

.instance-stats__item {
  display: grid;
  gap: 0.55rem;
  min-height: 8rem;
  align-content: space-between;
  padding: 1rem;
  border-radius: 1.25rem;
  background: rgba(255, 255, 255, 0.58);
  border: 1px solid rgba(110, 129, 147, 0.12);
}

.instance-stats__item[data-tone='ready'] {
  box-shadow: inset 0 0 0 1px rgba(111, 202, 164, 0.16);
}

.instance-stats__item[data-tone='pending'] {
  box-shadow: inset 0 0 0 1px rgba(232, 179, 103, 0.16);
}

.instance-stats__item[data-tone='muted'] {
  box-shadow: inset 0 0 0 1px rgba(154, 156, 182, 0.14);
}

.instance-stats__label {
  color: var(--text-soft);
  font-size: 0.88rem;
}

.instance-stats__value {
  font-size: clamp(1.45rem, 4vw, 2.15rem);
  color: var(--text-main);
  line-height: 1;
}

@media (max-width: 640px) {
  .instance-stats {
    grid-template-columns: 1fr;
  }
}
</style>
