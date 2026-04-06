<script setup lang="ts">
import { computed } from 'vue'
import StatusBadge from '@/components/shared/StatusBadge.vue'
import type { ViewStatusSummary } from '@/types/view'

const props = defineProps<{
  status: ViewStatusSummary | null
}>()

const clusterTone = computed(() => {
  if (!props.status) return 'unknown'
  if (props.status.clusterStatus === 'running') return 'running'
  if (props.status.clusterStatus === 'stopped') return 'stopped'
  return 'pending'
})

const clusterLabel = computed(() => {
  if (!props.status) return '状态未连接'
  if (props.status.clusterStatus === 'running') return '运行中'
  if (props.status.clusterStatus === 'stopped') return '已停止'
  return '状态异常'
})
</script>

<template>
  <div class="top-status-bar">
    <div class="top-status-bar__meta top-status-bar__meta--left">
      <StatusBadge :label="clusterLabel" :tone="clusterTone" />
      <span class="top-status-bar__chip">状态总览</span>
    </div>
  </div>
</template>

<style scoped>
.top-status-bar {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-start;
  gap: 0.65rem 1rem;
  padding: 0.72rem 0.92rem;
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.66);
  border: 1px solid rgba(103, 132, 148, 0.14);
  align-items: center;
}

.top-status-bar__meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.6rem;
  color: var(--text-soft);
  font-size: 0.88rem;
}

.top-status-bar__chip {
  display: inline-flex;
  align-items: center;
  padding: 0.38rem 0.65rem;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.64);
  color: var(--text-soft);
}

</style>
