<script setup lang="ts">
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import { useLogsStore } from '@/stores/logs'

const logsStore = useLogsStore()
const filterStart = ref('')
const filterEnd = ref('')

const {
  currentState,
  currentErrorMessage,
  visibleItems,
  hasActiveFilters,
} = storeToRefs(logsStore)

const isSubmitting = computed(() => currentState.value === 'loading')

async function applyFilters() {
  await logsStore.applyDateFilter({
    start: filterStart.value,
    end: filterEnd.value,
  })
}

function clearFilters() {
  filterStart.value = ''
  filterEnd.value = ''
  logsStore.clearDateFilter()
}

async function retry() {
  if (hasActiveFilters.value) {
    await logsStore.loadFilteredLogs(true)
    return
  }
  await logsStore.loadLogs(true)
}
</script>

<template>
  <BaseCard eyebrow="操作档案" title="日志中心" accent="gold">
    <div class="logs-view__toolbar">
      <label class="logs-view__field">
        <span>开始时间</span>
        <input v-model="filterStart" class="logs-view__input" type="datetime-local" />
      </label>
      <label class="logs-view__field">
        <span>结束时间</span>
        <input v-model="filterEnd" class="logs-view__input" type="datetime-local" />
      </label>
      <div class="logs-view__actions">
        <button type="button" class="logs-view__ghost" :disabled="isSubmitting" @click="clearFilters">清空筛选</button>
        <button type="button" class="logs-view__primary" :disabled="isSubmitting" @click="applyFilters">应用筛选</button>
      </div>
    </div>

    <p class="logs-view__hint">
      {{ hasActiveFilters ? '当前展示已筛选的操作日志。' : '当前展示最近的操作日志。' }}
    </p>

    <InlineError v-if="currentState === 'error'" :message="currentErrorMessage" button-label="重试" @retry="retry" />
    <EmptyState
      v-else-if="currentState === 'success' && !visibleItems.length"
      title="暂无日志记录"
      description="当前筛选条件下没有可展示的操作日志。"
    />
    <div v-else class="logs-view__list">
      <article v-for="log in visibleItems" :key="log.id" class="logs-view__item" :data-tone="log.tone">
        <div class="logs-view__row">
          <strong>{{ log.actionLabel }}</strong>
          <span>{{ log.timestampLabel }}</span>
        </div>
        <p class="logs-view__target">{{ log.targetLabel }}</p>
        <p class="logs-view__details">{{ log.details }}</p>
      </article>
    </div>
  </BaseCard>
</template>

<style scoped>
.logs-view__toolbar {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr)) auto;
  gap: 0.9rem;
  align-items: end;
}

.logs-view__field {
  display: grid;
  gap: 0.35rem;
}

.logs-view__field span,
.logs-view__hint,
.logs-view__row span,
.logs-view__target,
.logs-view__details {
  color: var(--text-soft);
}

.logs-view__input {
  width: 100%;
  padding: 0.8rem 0.9rem;
  border-radius: 1rem;
  border: 1px solid rgba(104, 128, 145, 0.18);
  background: rgba(255, 255, 255, 0.68);
}

.logs-view__actions {
  display: flex;
  gap: 0.7rem;
  flex-wrap: wrap;
}

.logs-view__ghost,
.logs-view__primary {
  border: 0;
  border-radius: 999px;
  padding: 0.8rem 1rem;
  cursor: pointer;
}

.logs-view__ghost {
  background: rgba(255, 255, 255, 0.68);
}

.logs-view__primary {
  background: linear-gradient(135deg, #bb7a22, #db9d3f);
  color: #fff;
}

.logs-view__hint {
  margin: 0.9rem 0 0;
}

.logs-view__list {
  display: grid;
  gap: 0.8rem;
  margin-top: 1rem;
}

.logs-view__item {
  padding: 0.95rem 1rem;
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.58);
}

.logs-view__item[data-tone='error'] {
  background: rgba(255, 240, 244, 0.84);
}

.logs-view__row {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.logs-view__target,
.logs-view__details {
  margin: 0.45rem 0 0;
}

@media (max-width: 900px) {
  .logs-view__toolbar {
    grid-template-columns: 1fr;
  }
}
</style>
