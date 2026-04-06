<script setup lang="ts">
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import { computed } from 'vue'
import type { LoadState, ViewAccount } from '@/types/view'

const props = defineProps<{
  state: LoadState
  errorMessage: string
  accounts: ViewAccount[]
}>()

const emit = defineEmits<{
  retry: []
}>()

const stats = computed(() => {
  const ready = props.accounts.filter((account) => account.authStatusTone === 'ready').length
  const pending = props.accounts.filter((account) => account.authStatusTone === 'pending').length
  const disabled = props.accounts.filter((account) => account.authStatusTone === 'disabled').length
  const invalid = props.accounts.filter((account) => account.authStatusTone === 'invalid' || account.authStatusTone === 'unknown').length
  const availability = props.accounts.length > 0 ? `${Math.round((ready / props.accounts.length) * 100)}%` : '0%'

  return [
    { label: '总账号数', value: props.accounts.length, tone: 'neutral' },
    { label: '已就绪', value: ready, tone: 'ready' },
    { label: '待处理', value: pending, tone: 'pending' },
    { label: '可用率', value: availability, tone: 'muted' },
  ]
})
</script>

<template>
  <BaseCard eyebrow="账号" title="账号状态" accent="peach">
    <InlineError v-if="state === 'error'" :message="errorMessage" button-label="重试" @retry="emit('retry')" />
    <EmptyState
      v-else-if="state === 'success' && !accounts.length"
      title="暂无账号"
      description="请先创建账号，再进行 Microsoft 登录。"
    />
    <div v-else class="summary-stats">
      <article v-for="item in stats" :key="item.label" class="summary-stats__item" :data-tone="item.tone">
        <span class="summary-stats__label">{{ item.label }}</span>
        <strong class="summary-stats__value">{{ item.value }}</strong>
      </article>
    </div>
  </BaseCard>
</template>

<style scoped>
.summary-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.85rem;
}

.summary-stats__item {
  display: grid;
  gap: 0.55rem;
  min-height: 8rem;
  align-content: space-between;
  padding: 1rem;
  border-radius: 1.25rem;
  background: rgba(255, 255, 255, 0.58);
  border: 1px solid rgba(110, 129, 147, 0.12);
}

.summary-stats__item[data-tone='ready'] {
  box-shadow: inset 0 0 0 1px rgba(111, 202, 164, 0.16);
}

.summary-stats__item[data-tone='pending'] {
  box-shadow: inset 0 0 0 1px rgba(232, 179, 103, 0.16);
}

.summary-stats__item[data-tone='muted'] {
  box-shadow: inset 0 0 0 1px rgba(154, 156, 182, 0.14);
}

.summary-stats__label {
  color: var(--text-soft);
  font-size: 0.88rem;
}

.summary-stats__value {
  font-size: clamp(1.45rem, 4vw, 2.15rem);
  color: var(--text-main);
  line-height: 1;
}

@media (max-width: 640px) {
  .summary-stats {
    grid-template-columns: 1fr;
  }
}
</style>
