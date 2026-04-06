<script setup lang="ts">
import { shallowRef, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import BaseCard from '@/components/shared/BaseCard.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import InlineError from '@/components/shared/InlineError.vue'
import CreateInstanceDialog from '@/components/instances/CreateInstanceDialog.vue'
import InstanceCard from '@/components/instances/InstanceCard.vue'
import { useAccountsStore } from '@/stores/accounts'
import { useInstancesStore } from '@/stores/instances'

const instancesStore = useInstancesStore()
const accountsStore = useAccountsStore()
const { state, errorMessage, items, submitting, actionBusyMap } = storeToRefs(instancesStore)
const { items: accounts } = storeToRefs(accountsStore)
const createOpen = shallowRef(false)

onMounted(() => {
  if (state.value === 'idle') {
    void instancesStore.loadInstances()
  }
  if (accounts.value.length === 0) {
    void accountsStore.loadAccounts()
  }
})

async function handleSubmit(payload: { id: string; accountId: string; serverAddress: string; enabled: boolean; autoStart: boolean }) {
  const result = await instancesStore.createInstance(payload)
  if (result.valid) {
    createOpen.value = false
  }
}
</script>

<template>
  <section class="instances-view">
    <BaseCard eyebrow="实例管理" title="实例管理" accent="sky">
      <div class="instances-view__toolbar">
        <div>
          <h2 class="instances-view__headline">管理实例</h2>
          <p class="instances-view__copy">创建实例并执行启动、停止、重启等操作。</p>
        </div>
        <button type="button" class="instances-view__button" @click="createOpen = true">新增实例</button>
      </div>

      <InlineError v-if="state === 'error'" :message="errorMessage" button-label="重试" @retry="instancesStore.loadInstances" />
      <EmptyState
        v-else-if="state === 'success' && !items.length"
        title="暂无实例"
        description="请先准备已登录账号，再创建实例。"
        button-label="创建实例"
        @action="createOpen = true"
      />
      <div v-else class="instances-view__grid">
        <InstanceCard
          v-for="instance in items"
          :key="instance.id"
          :instance="instance"
          :busy-action="actionBusyMap[instance.id] || ''"
          @action="instancesStore.runAction($event.id, $event.action)"
          @delete="instancesStore.deleteInstance"
        />
      </div>
    </BaseCard>

    <CreateInstanceDialog :open="createOpen" :submitting="submitting" :accounts="accounts" @close="createOpen = false" @submit="handleSubmit" />
  </section>
</template>

<style scoped>
.instances-view {
  display: grid;
  gap: 1rem;
}

.instances-view__toolbar {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}

.instances-view__headline,
.instances-view__copy {
  margin: 0;
}

.instances-view__copy {
  margin-top: 0.45rem;
  color: var(--text-soft);
}

.instances-view__button {
  align-self: start;
  border: none;
  border-radius: 999px;
  padding: 0.8rem 1rem;
  background: linear-gradient(135deg, #b8ebff, #bdd1ff);
  color: #31516a;
  font-weight: 700;
  cursor: pointer;
}

.instances-view__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.85rem;
}

@media (max-width: 900px) {
  .instances-view__toolbar {
    flex-direction: column;
  }

  .instances-view__grid {
    grid-template-columns: 1fr;
  }
}
</style>
