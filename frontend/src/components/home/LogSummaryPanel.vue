<script setup lang="ts">
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import type { LoadState, ViewLogItem } from '@/types/view'

defineProps<{
  state: LoadState
  errorMessage: string
  logs: ViewLogItem[]
}>()

const emit = defineEmits<{
  retry: []
}>()
</script>

<template>
  <BaseCard eyebrow="操作日志" title="操作记录" accent="gold">
    <InlineError v-if="state === 'error'" :message="errorMessage" button-label="重试" @retry="emit('retry')" />
    <EmptyState v-else-if="state === 'success' && !logs.length" title="暂无操作记录" description="当前还没有可显示的操作日志。" />
    <div v-else class="log-list log-list--scroll">
      <article v-for="log in logs" :key="log.id" class="log-list__item" :data-tone="log.tone">
        <div class="log-list__row">
          <strong>{{ log.actionLabel }}</strong>
          <span>{{ log.timestampLabel }}</span>
        </div>
        <p class="log-list__target">{{ log.targetLabel }}</p>
        <p class="log-list__details">{{ log.details }}</p>
      </article>
    </div>
  </BaseCard>
</template>

<style scoped>
.log-list {
  display: grid;
  gap: 0.8rem;
  min-height: 18rem;
  align-content: start;
}

.log-list--scroll {
  max-height: 23rem;
  overflow: auto;
  padding-right: 0.25rem;
}

.log-list__item {
  padding: 0.82rem 0.9rem;
  border-radius: 1.2rem;
  background: rgba(255, 255, 255, 0.58);
}

.log-list__item[data-tone='error'] {
  background: rgba(255, 240, 244, 0.84);
}

.log-list__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.log-list__row span,
.log-list__target,
.log-list__details {
  color: var(--text-soft);
}

.log-list__target,
.log-list__details {
  margin: 0.45rem 0 0;
}
</style>
