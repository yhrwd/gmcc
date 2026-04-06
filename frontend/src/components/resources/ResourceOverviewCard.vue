<script setup lang="ts">
import BaseCard from '@/components/shared/BaseCard.vue'
import type { ViewResourceSnapshot } from '@/types/view'

defineProps<{
  snapshot: ViewResourceSnapshot | null
}>()
</script>

<template>
  <BaseCard eyebrow="资源概况" title="资源快照" accent="mint">
    <div v-if="snapshot" class="resource-card">
      <div class="resource-card__item">
        <span>CPU</span>
        <strong>{{ snapshot.cpuPercent ?? '--' }}%</strong>
      </div>
      <div class="resource-card__item">
        <span>内存</span>
        <strong>{{ snapshot.memoryPercent ?? '--' }}%</strong>
      </div>
      <div class="resource-card__item">
        <span>已用 / 总量</span>
        <strong>{{ snapshot.usedMemoryLabel }} / {{ snapshot.totalMemoryLabel }}</strong>
      </div>
      <div class="resource-card__item">
        <span>采样时间</span>
        <strong>{{ snapshot.collectedAtLabel }}</strong>
      </div>
    </div>
  </BaseCard>
</template>

<style scoped>
.resource-card {
  display: grid;
  gap: 0.8rem;
}

.resource-card__item {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.85rem 1rem;
  border-radius: 1.2rem;
  background: rgba(255, 255, 255, 0.58);
}

.resource-card__item span {
  color: var(--text-soft);
}
</style>
