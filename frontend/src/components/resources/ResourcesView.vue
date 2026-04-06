<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import ResourceOverviewCard from '@/components/resources/ResourceOverviewCard.vue'
import { useHomeStore } from '@/stores/home'

const homeStore = useHomeStore()
const { resourcesState: state, resourcesError: errorMessage, resourceSnapshot: snapshot } = storeToRefs(homeStore)

const comfortCopy = computed(() => {
  if (!snapshot.value) return '等待新的资源快照...'
  if (snapshot.value.comfortLevel === 'comfort') return '当前资源状态稳定。'
  if (snapshot.value.comfortLevel === 'busy') return '当前资源负载有所上升，请注意观察。'
  if (snapshot.value.comfortLevel === 'tense') return '当前资源较为紧张，建议控制实例数量。'
  return '当前没有可用资源数据。'
})

onMounted(() => {
  if (state.value === 'idle') {
    void homeStore.loadResources()
  }
})
</script>

<template>
  <section class="resources-view">
    <ResourceOverviewCard :snapshot="snapshot" />
    <BaseCard eyebrow="资源说明" title="资源状态" accent="gold">
      <InlineError v-if="state === 'error'" :message="errorMessage" button-label="重试" @retry="homeStore.loadResources(true)" />
      <EmptyState v-else-if="state === 'success' && !snapshot" title="暂无资源数据" description="当前没有采集到可显示的资源信息。" />
      <p v-else class="resources-view__copy">{{ comfortCopy }}</p>
    </BaseCard>
  </section>
</template>

<style scoped>
.resources-view {
  display: grid;
  gap: 1rem;
}

.resources-view__copy {
  margin: 0;
  color: var(--text-soft);
  line-height: 1.7;
}
</style>
